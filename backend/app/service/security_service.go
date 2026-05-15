package service

import (
	"datauptwo/app/dto"
	"datauptwo/app/model"
	"datauptwo/app/repo"
	"datauptwo/global"
	"net"
	"strings"
	"time"
)

type SecurityService struct {
	settingsRepo   *repo.SecuritySettingsRepo
	lockoutRepo    *repo.LoginLockoutRepo
	userRepo       *repo.UserRepo
}

func NewSecurityService() *SecurityService {
	return &SecurityService{
		settingsRepo:   repo.NewSecuritySettingsRepo(),
		lockoutRepo:    repo.NewLoginLockoutRepo(),
		userRepo:       repo.NewUserRepo(),
	}
}

// GetSettings 获取安全设置
func (s *SecurityService) GetSettings() (*model.SecuritySettings, error) {
	return s.getSettingsWithDefault(), nil
}

// UpdateSettings 更新安全设置
func (s *SecurityService) UpdateSettings(req dto.SecuritySettingsUpdateReq) error {
	settings, err := s.settingsRepo.Get()
	if err != nil {
		return err
	}
	isNew := false
	if settings == nil {
		// 记录不存在，创建新的
		isNew = true
		settings = &model.SecuritySettings{}
		// 设置默认值
		settings.CaptchaEnabled = true
		settings.CaptchaMinLen = 3
		settings.RegistrationEnabled = true
		settings.InactiveAutoDisable = false
		settings.InactiveDaysThreshold = 90
		settings.UserLoginMaxAttempts = 5
		settings.UserLoginLockMinutes = 30
		settings.IPLoginMaxAttempts = 20
		settings.IPLoginLockMinutes = 60
		settings.PasswordExpiryDays = 0
		settings.PasswordMinLength = 8
		settings.PasswordRequireUppercase = false
		settings.PasswordRequireLowercase = false
		settings.PasswordRequireDigit = false
		settings.PasswordRequireSpecial = false
		settings.SessionTimeoutHours = 24
	}

	settings.CaptchaEnabled = req.CaptchaEnabled
	settings.CaptchaMinLen = req.CaptchaMinLen
	settings.RegistrationEnabled = req.RegistrationEnabled
	settings.InactiveAutoDisable = req.InactiveAutoDisable
	settings.InactiveDaysThreshold = req.InactiveDaysThreshold
	settings.UserLoginMaxAttempts = req.UserLoginMaxAttempts
	settings.UserLoginLockMinutes = req.UserLoginLockMinutes
	settings.IPLoginMaxAttempts = req.IPLoginMaxAttempts
	settings.IPLoginLockMinutes = req.IPLoginLockMinutes
	settings.IPWhitelist = req.IPWhitelist
	settings.IPBlacklist = req.IPBlacklist
	settings.PasswordExpiryDays = req.PasswordExpiryDays
	settings.PasswordMinLength = req.PasswordMinLength
	settings.PasswordRequireUppercase = req.PasswordRequireUppercase
	settings.PasswordRequireLowercase = req.PasswordRequireLowercase
	settings.PasswordRequireDigit = req.PasswordRequireDigit
	settings.PasswordRequireSpecial = req.PasswordRequireSpecial
	settings.SessionTimeoutHours = req.SessionTimeoutHours

	if isNew {
		return s.settingsRepo.Create(settings)
	}
	return s.settingsRepo.Update(settings)
}

// CheckIPAccess 检查IP访问权限（黑名单/白名单）
func (s *SecurityService) CheckIPAccess(ip string) (allowed bool, reason string) {
	settings, err := s.settingsRepo.Get()
	if err != nil || settings == nil {
		return true, "" // 出错时或无记录时放行
	}

	// 检查黑名单
	if settings.IPBlacklist != "" {
		blacklist := strings.Split(settings.IPBlacklist, ",")
		for _, item := range blacklist {
			item = strings.TrimSpace(item)
			if item == "" {
				continue
			}
			if s.ipMatches(ip, item) {
				return false, "IP在黑名单中，禁止登录"
			}
		}
	}

	// 检查白名单（白名单优先，白名单非空则只允许白名单IP）
	if settings.IPWhitelist != "" {
		whitelist := strings.Split(settings.IPWhitelist, ",")
		whitelisted := false
		for _, item := range whitelist {
			item = strings.TrimSpace(item)
			if item == "" {
				continue
			}
			if s.ipMatches(ip, item) {
				whitelisted = true
				break
			}
		}
		if !whitelisted {
			return false, "IP不在白名单中，禁止登录"
		}
	}

	return true, ""
}

// CheckIPLockout 检查IP是否被锁定
func (s *SecurityService) CheckIPLockout(ip string) (locked bool, remainingMinutes int, err error) {
	settings, err := s.settingsRepo.Get()
	if err != nil || settings == nil {
		return false, 0, nil
	}

	record, err := s.lockoutRepo.GetByTarget(ip, "ip")
	if err != nil {
		return false, 0, nil // 无记录，未锁定
	}

	if record == nil || !record.Locked {
		return false, 0, nil
	}

	// 计算剩余锁定时间
	lockedAt := record.LockedAt
	if lockedAt == nil {
		return false, 0, nil
	}
	elapsed := time.Since(*lockedAt)
	remaining := time.Duration(settings.IPLoginLockMinutes)*time.Minute - elapsed
	if remaining <= 0 {
		// 已过期，自动解锁
		s.lockoutRepo.Unlock(ip, "ip")
		return false, 0, nil
	}
	return true, int(remaining.Minutes()), nil
}

// CheckUserLockout 检查用户是否被锁定
func (s *SecurityService) CheckUserLockout(username string) (locked bool, remainingMinutes int, message string) {
	settings, err := s.settingsRepo.Get()
	if err != nil || settings == nil {
		return false, 0, ""
	}

	record, err := s.lockoutRepo.GetByTarget(username, "user")
	if err != nil || record == nil {
		return false, 0, ""
	}

	if !record.Locked {
		return false, 0, ""
	}

	lockedAt := record.LockedAt
	if lockedAt == nil {
		return false, 0, ""
	}
	elapsed := time.Since(*lockedAt)
	remaining := time.Duration(settings.UserLoginLockMinutes)*time.Minute - elapsed
	if remaining <= 0 {
		s.lockoutRepo.Unlock(username, "user")
		return false, 0, ""
	}
	return true, int(remaining.Minutes()), "账号已被锁定，请在" + itoa(int(remaining.Minutes())) + "分钟后重试"
}

// GetRemainingAttempts 获取用户剩余登录尝试次数（自动清除过期失败记录）
func (s *SecurityService) GetRemainingAttempts(username string) (remaining int, maxAttempts int) {
	settings := s.getSettingsWithDefault()
	maxAttempts = settings.UserLoginMaxAttempts

	record, err := s.lockoutRepo.GetByTarget(username, "user")
	if err != nil || record == nil {
		return maxAttempts, maxAttempts
	}

	// 检查失败记录是否已过期（超过锁定时长则自动清除）
	if record.LockedAt != nil {
		elapsed := time.Since(*record.LockedAt)
		if time.Duration(settings.UserLoginLockMinutes)*time.Minute <= elapsed {
			// 已过锁定时间，自动清除
			s.lockoutRepo.ResetByTarget(username, "user")
			return maxAttempts, maxAttempts
		}
	}

	// 检查失败次数是否超过锁定阈值（超过则按锁定处理）
	if record.LockedAt != nil && record.FailCount >= maxAttempts {
		return 0, maxAttempts
	}

	// 普通失败记录也检查是否过期（超过锁定时长但未达到锁定阈值）
	if record.LockedAt == nil && record.FailCount > 0 {
		// 用 UpdatedAt 检查旧记录是否过期
		if !record.UpdatedAt.IsZero() {
			elapsed := time.Since(record.UpdatedAt)
			if time.Duration(settings.UserLoginLockMinutes)*time.Minute <= elapsed {
				// 已过有效期，清除旧记录
				s.lockoutRepo.ResetByTarget(username, "user")
				return maxAttempts, maxAttempts
			}
		}
	}

	remaining = maxAttempts - record.FailCount
	if remaining < 0 {
		remaining = 0
	}
	return remaining, maxAttempts
}

// RecordLoginFailure 记录登录失败（增加失败计数，可能触发锁定）
func (s *SecurityService) RecordLoginFailure(username, ip string) (locked bool, lockType string, message string) {
	settings := s.getSettingsWithDefault()
	if settings == nil {
		return false, "", ""
	}

	// 用户级别锁定
	s.recordFailure(username, "user", settings.UserLoginMaxAttempts)
	userRecord, _ := s.lockoutRepo.GetByTarget(username, "user")
	if userRecord != nil && userRecord.Locked {
		return true, "user", "账号已被锁定，请稍后再试或联系管理员"
	}

	// IP级别锁定
	s.recordFailure(ip, "ip", settings.IPLoginMaxAttempts)
	ipRecord, _ := s.lockoutRepo.GetByTarget(ip, "ip")
	if ipRecord != nil && ipRecord.Locked {
		return true, "ip", "IP已被锁定，请稍后再试"
	}

	return false, "", ""
}

func (s *SecurityService) recordFailure(target, lockType string, maxAttempts int) {
	settings := s.getSettingsWithDefault()
	lockMinutes := settings.UserLoginLockMinutes
	if lockType == "ip" {
		lockMinutes = settings.IPLoginLockMinutes
	}
	lockDuration := time.Duration(lockMinutes) * time.Minute

	record, err := s.lockoutRepo.GetByTarget(target, lockType)
	if err != nil || record == nil {
		// 无记录，创建
		record = &model.LoginLockout{
			Target:    target,
			Type:      lockType,
			FailCount: 1,
			Locked:    false,
		}
		s.lockoutRepo.CreateOrUpdate(record)
		return
	}

	// 检查是否已锁定且未解锁
	if record.Locked && record.LockedAt != nil {
		elapsed := time.Since(*record.LockedAt)
		if elapsed < lockDuration {
			return // 仍在锁定中，不处理失败计数
		}
		// 已过锁定时间，解锁并重置计数
		record.Locked = false
		record.LockedAt = nil
		record.FailCount = 0
		s.lockoutRepo.CreateOrUpdate(record)
		return
	}

	// 检查上次失败是否已超时（普通失败计数，超时后自动清除重新计数）
	if !record.Locked && record.FailCount > 0 {
		elapsed := time.Since(record.UpdatedAt)
		if elapsed >= lockDuration {
			record.FailCount = 0
		}
	}

	record.FailCount++
	if record.FailCount >= maxAttempts {
		record.Locked = true
		now := time.Now()
		record.LockedAt = &now
		record.Reason = "连续登录失败次数过多"
	}
	s.lockoutRepo.CreateOrUpdate(record)
}


// RecordLoginSuccess 登录成功时重置失败计数
func (s *SecurityService) RecordLoginSuccess(username, ip string) {
	s.lockoutRepo.ResetByTarget(username, "user")
	s.lockoutRepo.ResetByTarget(ip, "ip")
}

// UnlockUser 解锁用户
func (s *SecurityService) UnlockUser(username string) error {
	return s.lockoutRepo.Unlock(username, "user")
}

// UnlockIP 解锁IP
func (s *SecurityService) UnlockIP(ip string) error {
	return s.lockoutRepo.Unlock(ip, "ip")
}

// GetLockedUsers 获取被锁定的用户列表
func (s *SecurityService) GetLockedUsers() ([]model.LoginLockout, error) {
	return s.lockoutRepo.ListLockedUsers(100)
}

// GetLockedIPs 获取被锁定的IP列表
func (s *SecurityService) GetLockedIPs() ([]model.LoginLockout, error) {
	return s.lockoutRepo.ListLockedIPs(100)
}

// GetSecurityOverview 获取安全总览
func (s *SecurityService) GetSecurityOverview() (*dto.SecurityOverviewResp, error) {
	resp := &dto.SecurityOverviewResp{}

	lockedUsers, _ := s.lockoutRepo.CountLockedUsers()
	resp.TotalLockedUsers = int(lockedUsers)

	lockedIPs, _ := s.lockoutRepo.CountLockedIPs()
	resp.TotalLockedIPs = int(lockedIPs)

	settings, err := s.settingsRepo.Get()
	if err == nil && settings != nil && settings.IPWhitelist != "" {
		resp.WhitelistCount = len(strings.Split(settings.IPWhitelist, ","))
	}
	if err == nil && settings != nil && settings.IPBlacklist != "" {
		resp.BlacklistCount = len(strings.Split(settings.IPBlacklist, ","))
	}

	return resp, nil
}

// ShouldRequireCaptcha 检查登录时是否需要验证码
func (s *SecurityService) ShouldRequireCaptcha(username string) (required bool, err error) {
	settings, err := s.settingsRepo.Get()
	if err != nil || settings == nil {
		return false, nil
	}
	if !settings.CaptchaEnabled {
		return false, nil
	}

	// 检查该用户的失败次数是否达到启用验证码的阈值
	record, err := s.lockoutRepo.GetByTarget(username, "user")
	if err != nil || record == nil {
		return false, nil
	}
	return record.FailCount >= settings.CaptchaMinLen, nil
}

// IsRegistrationEnabled 检查是否允许自主注册
func (s *SecurityService) IsRegistrationEnabled() bool {
	settings := s.getSettingsWithDefault()
	if settings == nil {
		return true // 出错时默认允许
	}
	return settings.RegistrationEnabled
}

// ValidatePasswordPolicy 验证密码策略
func (s *SecurityService) ValidatePasswordPolicy(password string) (valid bool, reason string) {
	settings := s.getSettingsWithDefault()
	if settings == nil {
		return true, "" // 出错时跳过
	}

	if len(password) < settings.PasswordMinLength {
		return false, "密码长度不能少于" + itoa(settings.PasswordMinLength) + "位"
	}
	if settings.PasswordRequireUppercase && !containsUppercase(password) {
		return false, "密码必须包含大写字母"
	}
	if settings.PasswordRequireLowercase && !containsLowercase(password) {
		return false, "密码必须包含小写字母"
	}
	if settings.PasswordRequireDigit && !containsDigit(password) {
		return false, "密码必须包含数字"
	}
	if settings.PasswordRequireSpecial && !containsSpecial(password) {
		return false, "密码必须包含特殊字符"
	}
	return true, ""
}

// getSettingsWithDefault 获取设置，不存在时返回默认配置
func (s *SecurityService) getSettingsWithDefault() *model.SecuritySettings {
	settings, err := s.settingsRepo.Get()
	if err == nil && settings != nil {
		return settings
	}
	// 返回默认安全设置
	return &model.SecuritySettings{
		CaptchaEnabled:           true,
		CaptchaMinLen:            3,
		InactiveAutoDisable:       false,
		InactiveDaysThreshold:    90,
		UserLoginMaxAttempts:     5,
		UserLoginLockMinutes:    30,
		IPLoginMaxAttempts:       20,
		IPLoginLockMinutes:      60,
		IPWhitelist:             "",
		IPBlacklist:             "",
		PasswordExpiryDays:      0,
		PasswordMinLength:       8,
		PasswordRequireUppercase: false,
		PasswordRequireLowercase: false,
		PasswordRequireDigit:    false,
		PasswordRequireSpecial:  false,
		SessionTimeoutHours:     24,
		RegistrationEnabled:     true,
	}
}

// ValidatePassword 包级函数，校验密码复杂度（供其他 service 调用）
func ValidatePassword(password string) (bool, string) {
	return NewSecurityService().ValidatePasswordPolicy(password)
}

// CheckInactiveUsers 检查不活跃用户并自动禁用
func (s *SecurityService) CheckInactiveUsers() error {
	settings, err := s.settingsRepo.Get()
	if err != nil || settings == nil || !settings.InactiveAutoDisable {
		return nil
	}

	threshold := time.Now().AddDate(0, 0, -settings.InactiveDaysThreshold)
	return global.DB.Model(&model.User{}).
		Where("last_login_at < ? AND status = ? AND source = ?", threshold, "active", "LOCAL").
		Update("status", "inactive").Error
}

// ipMatches 判断IP是否匹配CIDR或精确IP
func (s *SecurityService) ipMatches(ip, cidr string) bool {
	// 精确匹配
	if ip == cidr {
		return true
	}
	// CIDR 匹配
	if strings.Contains(cidr, "/") {
		_, netIP, err := net.ParseCIDR(cidr)
		if err != nil {
			return false
		}
		return netIP.Contains(net.ParseIP(ip))
	}
	return false
}

func containsUppercase(s string) bool {
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			return true
		}
	}
	return false
}

func containsLowercase(s string) bool {
	for _, c := range s {
		if c >= 'a' && c <= 'z' {
			return true
		}
	}
	return false
}

func containsDigit(s string) bool {
	for _, c := range s {
		if c >= '0' && c <= '9' {
			return true
		}
	}
	return false
}

func containsSpecial(s string) bool {
	for _, c := range s {
		if (c >= '!' && c <= '/') || (c >= ':' && c <= '@') || (c >= '[' && c <= '`') || (c >= '{' && c <= '~') {
			return true
		}
	}
	return false
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	result := ""
	for i > 0 {
		result = string(rune('0'+i%10)) + result
		i /= 10
	}
	return result
}
