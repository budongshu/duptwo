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
	settings, err := s.settingsRepo.Get()
	if err != nil {
		return nil, err
	}
	if settings != nil {
		return settings, nil
	}
	// 如果记录不存在，返回默认设置
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
		PasswordMinLength:       6,
		PasswordRequireUppercase: false,
		PasswordRequireLowercase: false,
		PasswordRequireDigit:    false,
		PasswordRequireSpecial:  false,
		SessionTimeoutHours:     24,
	}, nil
}

// UpdateSettings 更新安全设置
func (s *SecurityService) UpdateSettings(req dto.SecuritySettingsUpdateReq) error {
	settings, err := s.settingsRepo.Get()
	if err != nil {
		return err
	}
	if settings == nil {
		// 记录不存在，创建新的
		settings = &model.SecuritySettings{}
	}

	settings.CaptchaEnabled = req.CaptchaEnabled
	settings.CaptchaMinLen = req.CaptchaMinLen
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

// RecordLoginFailure 记录登录失败（增加失败计数，可能触发锁定）
func (s *SecurityService) RecordLoginFailure(username, ip string) (locked bool, lockType string, message string) {
	settings, err := s.settingsRepo.Get()
	if err != nil || settings == nil {
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

// ValidatePasswordPolicy 验证密码策略
func (s *SecurityService) ValidatePasswordPolicy(password string) (valid bool, reason string) {
	settings, err := s.settingsRepo.Get()
	if err != nil || settings == nil {
		return true, "" // 出错时或无记录时跳过
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
