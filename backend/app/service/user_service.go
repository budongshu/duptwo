package service

import (
	"datauptwo/app/dto"
	"datauptwo/app/model"
	"datauptwo/app/repo"
	"errors"

	"github.com/pquerna/otp/totp"
)

type UserService struct {
	userRepo  *repo.UserRepo
	roleRepo  *repo.RoleRepo
	groupRepo *repo.UserGroupRepo
}

func NewUserService() *UserService {
	return &UserService{
		userRepo:  repo.NewUserRepo(),
		roleRepo:  repo.NewRoleRepo(),
		groupRepo: repo.NewUserGroupRepo(),
	}
}

// Create 创建用户
func (s *UserService) Create(req dto.UserCreateReq) (*dto.UserResp, error) {
	// 检查用户名是否存在
	count, err := s.userRepo.CountByUsername(req.Username)
	if err != nil {
		return nil, errors.New("检查用户名失败")
	}
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	// 校验密码策略
	if valid, reason := ValidatePassword(req.Password); !valid {
		return nil, errors.New(reason)
	}

	// 哈希密码
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	user := &model.User{
		Username:   req.Username,
		Password:   hashedPassword,
		Nickname:   req.Nickname,
		Email:      req.Email,
		Phone:      req.Phone,
		Status:     "active",
		RoleID:     req.RoleID,
		GroupID:    req.GroupID,
		MFAEnabled: req.MFAEnabled,
	}

	// 如果启用MFA，自动生成密钥
	if req.MFAEnabled {
		totpSecret, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "DataRegistry",
			AccountName: req.Username,
		})
		if err != nil {
			return nil, errors.New("生成MFA密钥失败")
		}
		user.MFASecret = totpSecret.Secret()
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return s.toUserResp(user, false), nil
}

// Update 更新用户
func (s *UserService) Update(req dto.UserUpdateReq) (*dto.UserResp, error) {
	user, err := s.userRepo.GetByID(req.ID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	user.Nickname = req.Nickname
	user.Email = req.Email
	user.Phone = req.Phone
	user.RoleID = req.RoleID
	user.GroupID = req.GroupID

	if req.Status != "" {
		user.Status = req.Status
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return s.toUserResp(user, false), nil
}

// Delete 删除用户
func (s *UserService) Delete(id uint) error {
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}
	return s.userRepo.Delete(id)
}

// BatchDelete 批量删除用户
func (s *UserService) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return s.userRepo.BatchDelete(ids)
}

// BatchUpdateRole 批量更新用户角色
func (s *UserService) BatchUpdateRole(ids []uint, roleId uint) error {
	if len(ids) == 0 {
		return nil
	}
	// 验证角色是否存在
	role, err := s.roleRepo.GetByID(roleId)
	if err != nil {
		return errors.New("角色不存在")
	}
	if role == nil {
		return errors.New("角色不存在")
	}
	return s.userRepo.BatchUpdateRole(ids, roleId)
}

// GetByID 获取用户详情
func (s *UserService) GetByID(id uint) (*dto.UserResp, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return s.toUserResp(user, false), nil
}

// List 分页列表
func (s *UserService) List(req dto.UserListReq) (*dto.PageResult, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}

	users, total, err := s.userRepo.List(req)
	if err != nil {
		return nil, err
	}

	// 批量获取锁定状态
	lockedUsers, _ := repo.NewLoginLockoutRepo().ListLockedUsers(1000)
	lockedMap := make(map[string]bool)
	for _, lu := range lockedUsers {
		lockedMap[lu.Target] = true
	}

	items := make([]dto.UserResp, len(users))
	for i, user := range users {
		items[i] = *s.toUserResp(&user, lockedMap[user.Username])
	}

	return &dto.PageResult{Total: total, Items: items}, nil
}

// GetAll 获取所有用户
func (s *UserService) GetAll() ([]dto.UserResp, error) {
	req := dto.UserListReq{Page: 1, PageSize: 1000}
	users, _, err := s.userRepo.List(req)
	if err != nil {
		return nil, err
	}

	lockedUsers, _ := repo.NewLoginLockoutRepo().ListLockedUsers(1000)
	lockedMap := make(map[string]bool)
	for _, lu := range lockedUsers {
		lockedMap[lu.Target] = true
	}

	items := make([]dto.UserResp, len(users))
	for i, user := range users {
		items[i] = *s.toUserResp(&user, lockedMap[user.Username])
	}

	return items, nil
}

// toUserResp 转换为用户响应
func (s *UserService) toUserResp(user *model.User, locked bool) *dto.UserResp {
	resp := &dto.UserResp{
		ID:          user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Email:       user.Email,
		Phone:       user.Phone,
		Avatar:      user.Avatar,
		Status:      user.Status,
		StatusText:  "正常",
		Locked:      locked,
		RoleID:      user.RoleID,
		GroupID:     user.GroupID,
		Department:  user.Department,
		Title:       user.Title,
		Company:     user.Company,
		Source:      user.Source,
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

// ResetPassword 重置用户密码
func (s *UserService) ResetPassword(userID uint, newPassword string) error {
	if userID == 0 {
		return errors.New("用户ID无效")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 校验密码策略
	if valid, reason := ValidatePassword(newPassword); !valid {
		return errors.New(reason)
	}

	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	user.Password = hashedPassword
	return s.userRepo.Update(user)
}

// ResetMFA 重置用户MFA
func (s *UserService) ResetMFA(userID uint) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	user.MFAEnabled = false
	user.MFASecret = ""
	return s.userRepo.Update(user)
}

// GenerateMFASecret 生成MFA密钥（仅生成并保存，不启用）
func (s *UserService) GenerateMFASecret(userID uint) (secret, qrCode string, err error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		err = errors.New("用户不存在")
		return
	}

	// 生成新的secret
	totpSecret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "DataRegistry",
		AccountName: user.Username,
	})
	if err != nil {
		err = errors.New("生成MFA密钥失败")
		return
	}

	// 先将secret保存到数据库（不启用，等用户验证后再启用）
	user.MFASecret = totpSecret.Secret()
	if err = s.userRepo.Update(user); err != nil {
		err = errors.New("保存MFA密钥失败")
		return
	}

	secret = totpSecret.Secret()
	qrCode = totpSecret.URL()

	return
}

// AdminEnableMFA 管理员为用户启用MFA（验证验证码后启用）
func (s *UserService) AdminEnableMFA(userID uint, code string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 检查是否已生成secret
	if user.MFASecret == "" {
		return errors.New("请先生成MFA密钥")
	}

	// 验证TOTP验证码
	if !totp.Validate(code, user.MFASecret) {
		return errors.New("验证码错误")
	}

	user.MFAEnabled = true
	return s.userRepo.Update(user)
}

// ListForExport 获取所有用户（用于导出）
func (s *UserService) ListForExport(req dto.UserListReq) ([]model.User, error) {
	users, _, err := s.userRepo.ListForExport(req)
	return users, err
}

// GetRoleAndGroupMaps 获取角色和用户组的映射
func (s *UserService) GetRoleAndGroupMaps() (roleMap map[uint]string, groupMap map[uint]string) {
	roleMap = make(map[uint]string)
	groupMap = make(map[uint]string)

	roles, _ := s.roleRepo.GetAll()
	for _, r := range roles {
		roleMap[r.ID] = r.Name
	}

	groups, _ := s.groupRepo.GetAll()
	for _, g := range groups {
		groupMap[g.ID] = g.Name
	}

	return
}
