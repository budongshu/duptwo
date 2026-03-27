package service

import (
	"crypto/rand"
	"datauptwo/app/dto"
	"datauptwo/app/ldap"
	"datauptwo/app/model"
	"datauptwo/app/repo"
	"datauptwo/global"
	"datauptwo/middleware"
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo     *repo.UserRepo
	roleRepo     *repo.RoleRepo
	groupRepo    *repo.UserGroupRepo
	loginLogRepo *repo.LoginLogRepo
	adService    *ADService
	securityService *SecurityService
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo:     repo.NewUserRepo(),
		roleRepo:     repo.NewRoleRepo(),
		groupRepo:    repo.NewUserGroupRepo(),
		loginLogRepo: repo.NewLoginLogRepo(),
		adService:    NewADService(),
		securityService: NewSecurityService(),
	}
}

// Login 用户登录（统一入口，自动判断本地/AD用户）
func (s *AuthService) Login(req dto.LoginReq, ip, userAgent string) (*dto.LoginResp, error) {
	// ========== 安全检查 ==========
	security := s.securityService

	// 1. IP黑名单/白名单检查
	if allowed, reason := security.CheckIPAccess(ip); !allowed {
		s.loginLogRepo.Create(&model.LoginLog{
			Username:   req.Username,
			Status:     "failed",
			IPAddress:  ip,
			UserAgent:  userAgent,
			FailReason: reason,
		})
		return nil, errors.New(reason)
	}

	// 2. IP锁定检查
	if locked, remaining, _ := security.CheckIPLockout(ip); locked {
		return nil, errors.New("IP已被锁定，请在" + itoa(remaining) + "分钟后重试")
	}

	// 3. 用户登录失败记录（无论用户是否存在，都记录，用于触发验证码和锁定）
	security.RecordLoginFailure(req.Username, ip)

	// 4. 用户被锁定检查
	if locked, _, msg := security.CheckUserLockout(req.Username); locked {
		return nil, errors.New(msg)
	}

	// 5. 验证码检查（仅在用户名存在且失败次数达标时要求）
	if captchaRequired, _ := security.ShouldRequireCaptcha(req.Username); captchaRequired {
		if req.CaptchaID == "" || req.Captcha == "" {
			return nil, errors.New("请填写验证码")
		}
		if err := s.ValidateCaptcha(req.CaptchaID, req.Captcha); err != nil {
			return nil, err
		}
	}

	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		// 用户不存在
		// 如果AD已启用且允许自动注册，尝试AD认证后创建AD用户
		if s.adService.IsEnabled() && global.CONF.AD.AutoRegister {
			return s.ADAutoLogin(req.Username, req.Password, ip, userAgent)
		}
		s.loginLogRepo.Create(&model.LoginLog{
			Username:   req.Username,
			Status:     "failed",
			IPAddress:  ip,
			UserAgent:  userAgent,
			FailReason: "用户不存在",
		})
		return nil, errors.New("用户名或密码错误")
	}

	if user.Status != "active" {
		s.loginLogRepo.Create(&model.LoginLog{
			UserID:     user.ID,
			Username:   user.Username,
			Status:     "failed",
			IPAddress:  ip,
			UserAgent:  userAgent,
			FailReason: "账号已被禁用",
		})
		return nil, errors.New("账号已被禁用，请联系管理员")
	}

	// 根据用户来源验证密码
	if user.Source == "AD" {
		// AD用户，使用AD认证
		adInfo, err := s.adService.Authenticate(user.Username, req.Password)
		if err != nil {
			s.securityService.RecordLoginFailure(user.Username, ip)
			s.loginLogRepo.Create(&model.LoginLog{
				UserID:     user.ID,
				Username:   user.Username,
				Status:     "failed",
				IPAddress:  ip,
				UserAgent:  userAgent,
				FailReason: "AD认证失败",
				LoginMethod: "ad",
			})
			return nil, errors.New("用户名或密码错误")
		}
		// AD认证成功，更新用户信息（如果AD中有变化）
		s.syncADUserInfo(user, adInfo)
	} else {
		// 本地用户，使用本地密码验证
		if !CheckPassword(user.Password, req.Password) {
			s.securityService.RecordLoginFailure(user.Username, ip)
			s.loginLogRepo.Create(&model.LoginLog{
				UserID:     user.ID,
				Username:   user.Username,
				Status:     "failed",
				IPAddress:  ip,
				UserAgent:  userAgent,
				FailReason: "密码错误",
			})
			return nil, errors.New("用户名或密码错误")
		}
	}

	// 检查是否启用了MFA
	if user.MFAEnabled {
		return &dto.LoginResp{
			MFARequired: true,
			User:        s.toUserInfo(user),
			ExpireAt:    time.Now().Add(5 * time.Minute).Unix(),
		}, nil
	}

	// 生成 JWT token
	token, err := middleware.GenerateToken(user.ID, user.Username, user.RoleID == 0, s.getUserPermissions(user)...)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	// 更新登录信息
	s.userRepo.UpdateLastLogin(user.ID, ip)

	// 重置登录失败计数
	s.securityService.RecordLoginSuccess(user.Username, ip)

	// 记录成功的登录
	loginMethod := "password"
	if user.Source == "AD" {
		loginMethod = "ad"
	}
	s.loginLogRepo.Create(&model.LoginLog{
		UserID:      user.ID,
		Username:    user.Username,
		Status:      "success",
		IPAddress:   ip,
		UserAgent:   userAgent,
		LoginMethod: loginMethod,
	})

	return &dto.LoginResp{
		Token:    token,
		ExpireAt: time.Now().Add(24 * 7 * time.Hour).Unix(),
		User:     s.toUserInfo(user),
	}, nil
}

// ADLogin AD用户登录（显式调用，跳过本地用户查找）
func (s *AuthService) ADLogin(req dto.ADLoginReq, ip, userAgent string) (*dto.LoginResp, error) {
	if !s.adService.IsEnabled() {
		return nil, errors.New("AD认证未启用")
	}

	// 验证验证码（如果启用）
	if req.CaptchaID != "" && req.Captcha != "" {
		if err := s.ValidateCaptcha(req.CaptchaID, req.Captcha); err != nil {
			return nil, err
		}
	}

	// 检查本地是否已有该用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		// 本地不存在，尝试自动注册
		return s.ADAutoLogin(req.Username, req.Password, ip, userAgent)
	}

	// 本地已存在该用户
	if user.Source != "AD" {
		// 本地用户，不能用AD登录
		return nil, errors.New("该用户为本地用户，请使用密码登录")
	}

	if user.Status != "active" {
		s.loginLogRepo.Create(&model.LoginLog{
			UserID:     user.ID,
			Username:   user.Username,
			Status:     "failed",
			IPAddress:  ip,
			UserAgent:  userAgent,
			FailReason: "账号已被禁用",
			LoginMethod: "ad",
		})
		return nil, errors.New("账号已被禁用，请联系管理员")
	}

	// 使用AD认证
	adInfo, adErr := s.adService.Authenticate(user.Username, req.Password)
	if adErr != nil {
		s.loginLogRepo.Create(&model.LoginLog{
			UserID:     user.ID,
			Username:   user.Username,
			Status:     "failed",
			IPAddress:  ip,
			UserAgent:  userAgent,
			FailReason: "AD认证失败",
			LoginMethod: "ad",
		})
		return nil, errors.New("用户名或密码错误")
	}

	// 同步AD用户信息
	s.syncADUserInfo(user, adInfo)

	// 检查MFA
	if user.MFAEnabled {
		return &dto.LoginResp{
			MFARequired: true,
			User:        s.toUserInfo(user),
			ExpireAt:    time.Now().Add(5 * time.Minute).Unix(),
		}, nil
	}

	// 生成token
	token, err := middleware.GenerateToken(user.ID, user.Username, user.RoleID == 0, s.getUserPermissions(user)...)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	s.userRepo.UpdateLastLogin(user.ID, ip)
	s.loginLogRepo.Create(&model.LoginLog{
		UserID:      user.ID,
		Username:    user.Username,
		Status:      "success",
		IPAddress:   ip,
		UserAgent:   userAgent,
		LoginMethod: "ad",
	})

	return &dto.LoginResp{
		Token:    token,
		ExpireAt: time.Now().Add(24 * 7 * time.Hour).Unix(),
		User:     s.toUserInfo(user),
	}, nil
}

// ADAutoLogin AD认证并自动注册用户（当本地不存在该用户时）
func (s *AuthService) ADAutoLogin(username, password, ip, userAgent string) (*dto.LoginResp, error) {
	if !s.adService.IsEnabled() {
		return nil, errors.New("AD认证未启用")
	}
	if !global.CONF.AD.AutoRegister {
		return nil, errors.New("该用户不存在，且系统不允许自动注册AD用户")
	}

	// AD认证
	adInfo, err := s.adService.Authenticate(username, password)
	if err != nil {
		s.loginLogRepo.Create(&model.LoginLog{
			Username:   username,
			Status:     "failed",
			IPAddress:  ip,
			UserAgent:  userAgent,
			FailReason: "AD认证失败",
			LoginMethod: "ad",
		})
		return nil, errors.New("用户名或密码错误")
	}

	// 创建AD用户
	user := &model.User{
		Username: username,
		Status:   "active",
		Source:   "AD",
		ADDN:     adInfo.DN,
		Nickname: adInfo.Nickname,
		Email:    adInfo.Email,
		Password: "", // AD用户无本地密码
		RoleID:   global.CONF.AD.DefaultRoleID,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("自动注册AD用户失败")
	}

	// 生成token
	token, err := middleware.GenerateToken(user.ID, user.Username, user.RoleID == 0, s.getUserPermissions(user)...)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	s.userRepo.UpdateLastLogin(user.ID, ip)
	s.loginLogRepo.Create(&model.LoginLog{
		UserID:      user.ID,
		Username:    user.Username,
		Status:      "success",
		IPAddress:   ip,
		UserAgent:   userAgent,
		LoginMethod: "ad",
	})

	return &dto.LoginResp{
		Token:    token,
		ExpireAt: time.Now().Add(24 * 7 * time.Hour).Unix(),
		User:     s.toUserInfo(user),
	}, nil
}

// syncADUserInfo 同步AD用户信息到本地
func (s *AuthService) syncADUserInfo(user *model.User, adInfo *ldap.ADUserInfo) {
	updated := false
	if adInfo.Email != "" && user.Email != adInfo.Email {
		user.Email = adInfo.Email
		updated = true
	}
	if adInfo.Nickname != "" && user.Nickname != adInfo.Nickname {
		user.Nickname = adInfo.Nickname
		updated = true
	}
	if adInfo.DN != "" && user.ADDN != adInfo.DN {
		user.ADDN = adInfo.DN
		updated = true
	}
	if updated {
		s.userRepo.Update(user)
	}
}

// MFAVerify MFA验证
func (s *AuthService) MFAVerify(req dto.MFAVerifyReq, ip, userAgent string) (*dto.LoginResp, error) {
	user, err := s.userRepo.GetByID(req.UserID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if !user.MFAEnabled {
		return nil, errors.New("用户未启用MFA")
	}

	// 验证TOTP code
	if !totp.Validate(req.Code, user.MFASecret) {
		s.loginLogRepo.Create(&model.LoginLog{
			UserID:      user.ID,
			Username:    user.Username,
			Status:      "failed",
			IPAddress:   ip,
			UserAgent:   userAgent,
			FailReason:  "MFA验证码错误",
			MFAUsed:     true,
			LoginMethod: "mfa",
		})
		return nil, errors.New("验证码错误或已过期")
	}

	// 生成 JWT token
	token, err := middleware.GenerateToken(user.ID, user.Username, user.RoleID == 0, s.getUserPermissions(user)...)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	// 更新登录信息
	s.userRepo.UpdateLastLogin(user.ID, ip)

	// 记录成功的MFA登录
	s.loginLogRepo.Create(&model.LoginLog{
		UserID:      user.ID,
		Username:    user.Username,
		Status:      "success",
		IPAddress:   ip,
		UserAgent:   userAgent,
		MFAUsed:     true,
		LoginMethod: "mfa",
	})

	return &dto.LoginResp{
		Token:    token,
		ExpireAt: time.Now().Add(24 * 7 * time.Hour).Unix(),
		User:     s.toUserInfo(user),
	}, nil
}

// EnableMFA 启用MFA
func (s *AuthService) EnableMFA(userID uint, req dto.MFAEnableReq) (*dto.MFAStatusResp, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 生成新的secret
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "DataRegistry",
		AccountName: user.Username,
	})
	if err != nil {
		return nil, errors.New("生成MFA密钥失败")
	}

	// 验证初始code
	if !totp.Validate(req.Code, secret.Secret()) {
		return nil, errors.New("验证码错误，请确保时间同步")
	}

	// 保存secret
	user.MFASecret = secret.Secret()
	user.MFAEnabled = true
	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("保存MFA设置失败")
	}

	// 生成备份码
	backupCodes := s.generateBackupCodes()

	return &dto.MFAStatusResp{
		Enabled:     true,
		Secret:      secret.Secret(),
		BackupCodes: backupCodes,
	}, nil
}

// DisableMFA 禁用MFA
func (s *AuthService) DisableMFA(userID uint) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	user.MFAEnabled = false
	user.MFASecret = ""
	return s.userRepo.Update(user)
}

// GetMFAStatus 获取MFA状态
func (s *AuthService) GetMFAStatus(userID uint) (*dto.MFAStatusResp, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	return &dto.MFAStatusResp{
		Enabled: user.MFAEnabled,
	}, nil
}

// generateBackupCodes 生成备份码（使用加密安全的随机数）
func (s *AuthService) generateBackupCodes() []string {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	codes := make([]string, 10)
	for i := 0; i < 10; i++ {
		code := ""
		for j := 0; j < 8; j++ {
			if j > 0 && j%4 == 0 {
				code += "-"
			}
			idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			code += string(charset[idx.Int64()])
		}
		codes[i] = code
	}
	return codes
}

// Register 用户注册
func (s *AuthService) Register(req dto.RegisterReq) (*dto.LoginResp, error) {
	// 检查用户名是否存在
	count, err := s.userRepo.CountByUsername(req.Username)
	if err != nil {
		return nil, errors.New("检查用户名失败")
	}
	if count > 0 {
		// 返回通用错误，避免泄露用户名是否存在
		return nil, errors.New("该用户名不可用")
	}

	// 哈希密码
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	user := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Email:    req.Email,
		Status:   "active",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("创建用户失败")
	}

	// 生成 JWT token
	token, err := middleware.GenerateToken(user.ID, user.Username, user.RoleID == 0, s.getUserPermissions(user)...)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	return &dto.LoginResp{
		Token:    token,
		ExpireAt: time.Now().Add(24 * 7 * time.Hour).Unix(),
		User:     s.toUserInfo(user),
	}, nil
}

// GetCurrentUser 获取当前用户信息
func (s *AuthService) GetCurrentUser(userID uint) (*dto.UserResp, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return s.toUserResp(user), nil
}

// UpdateProfile 更新个人资料
func (s *AuthService) UpdateProfile(userID uint, req dto.ProfileUpdateReq) (*dto.UserResp, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("更新失败")
	}
	return s.toUserResp(user), nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(userID uint, req dto.ChangePasswordReq) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if !CheckPassword(user.Password, req.OldPassword) {
		return errors.New("原密码错误")
	}

	// 哈希新密码
	hashedPassword, err := HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	user.Password = hashedPassword
	return s.userRepo.Update(user)
}

// toUserInfo 转换为用户信息
func (s *AuthService) toUserInfo(user *model.User) dto.UserInfo {
	info := dto.UserInfo{
		ID:         user.ID,
		Username:   user.Username,
		Nickname:   user.Nickname,
		Email:      user.Email,
		Phone:      user.Phone,
		Avatar:     user.Avatar,
		Status:     user.Status,
		MFAEnabled: user.MFAEnabled,
		Source:     user.Source,
	}
	if info.Source == "" {
		info.Source = "LOCAL"
	}

	// 获取角色信息
	if user.RoleID > 0 {
		if role, err := s.roleRepo.GetByID(user.RoleID); err == nil {
			info.RoleID = role.ID
			info.RoleName = role.Name
			info.RoleCode = role.Code
		}
	}

	// 填充权限列表
	info.Permissions = s.getUserPermissions(user)

	return info
}

// toUserResp 转换为用户响应
func (s *AuthService) toUserResp(user *model.User) *dto.UserResp {
	resp := &dto.UserResp{
		ID:          user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Email:       user.Email,
		Phone:       user.Phone,
		Avatar:      user.Avatar,
		Status:      user.Status,
		StatusText:  "正常",
		RoleID:      user.RoleID,
		GroupID:     user.GroupID,
		MFAEnabled:  user.MFAEnabled,
		LastLoginAt: user.LastLoginAt,
		LastLoginIP: user.LastLoginIP,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	if user.Status == "inactive" {
		resp.StatusText = "已禁用"
	}

	// 获取角色名称
	if user.RoleID > 0 {
		if role, err := s.roleRepo.GetByID(user.RoleID); err == nil {
			resp.RoleName = role.Name
		}
	}

	// 获取用户组名称
	if user.GroupID > 0 {
		if group, err := s.groupRepo.GetByID(user.GroupID); err == nil {
			resp.GroupName = group.Name
		}
	}

	return resp
}

// HashPassword 哈希密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证密码
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// getUserPermissions 获取用户的权限列表
func (s *AuthService) getUserPermissions(user *model.User) []string {
	if user.RoleID == 0 {
		return nil
	}
	role, err := s.roleRepo.GetByID(user.RoleID)
	if err != nil || role.Permissions == "" {
		return nil
	}
	var perms []string
	if err := json.Unmarshal([]byte(role.Permissions), &perms); err != nil {
		return nil
	}
	return perms
}
