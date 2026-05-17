package repo

import (
	"datauptwo/app/model"
	"datauptwo/global"
	"time"

	"gorm.io/gorm"
)

type SecuritySettingsRepo struct{}

func NewSecuritySettingsRepo() *SecuritySettingsRepo {
	return &SecuritySettingsRepo{}
}

func (r *SecuritySettingsRepo) Get() (*model.SecuritySettings, error) {
	var settings model.SecuritySettings
	err := global.DB.First(&settings).Error
	if err != nil {
		// 如果记录不存在，返回nil而不是错误，让调用方处理
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &settings, nil
}

func (r *SecuritySettingsRepo) Update(settings *model.SecuritySettings) error {
	if settings.ID == 0 {
		// 记录不存在，创建新记录
		return global.DB.Create(settings).Error
	}
	return global.DB.Save(settings).Error
}

func (r *SecuritySettingsRepo) Create(settings *model.SecuritySettings) error {
	return global.DB.Create(settings).Error
}

// SaveADConfig 保存 AD 域配置到 security_settings 表（仅更新 AD 相关字段，保留其他安全设置）
func (r *SecuritySettingsRepo) SaveADConfig(settings *model.SecuritySettings) error {
	// 获取现有记录，保留其他字段
	var existing model.SecuritySettings
	err := global.DB.First(&existing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return global.DB.Create(settings).Error
		}
		return err
	}
	// 仅更新 AD 相关字段
	existing.ADEnabled = settings.ADEnabled
	existing.ADServer = settings.ADServer
	existing.ADPort = settings.ADPort
	existing.ADUseSSL = settings.ADUseSSL
	existing.ADBaseDN = settings.ADBaseDN
	existing.ADBindDN = settings.ADBindDN
	existing.ADBindPassword = settings.ADBindPassword
	existing.ADUserFilter = settings.ADUserFilter
	existing.ADAutoRegister = settings.ADAutoRegister
	existing.ADDefaultRoleID = settings.ADDefaultRoleID
	existing.ADLastSyncAt = settings.ADLastSyncAt
	existing.ADLastSyncCount = settings.ADLastSyncCount
	return global.DB.Save(&existing).Error
}

type LoginLockoutRepo struct{}

func NewLoginLockoutRepo() *LoginLockoutRepo {
	return &LoginLockoutRepo{}
}

func (r *LoginLockoutRepo) GetByTarget(target, lockType string) (*model.LoginLockout, error) {
	var record model.LoginLockout
	err := global.DB.Where("`target` = ? AND `type` = ?", target, lockType).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

func (r *LoginLockoutRepo) CreateOrUpdate(record *model.LoginLockout) error {
	var existing model.LoginLockout
	err := global.DB.Where("`target` = ? AND `type` = ?", record.Target, record.Type).First(&existing).Error
	if err != nil {
		// 不存在，创建
		return global.DB.Create(record).Error
	}
	// 存在，只更新需要变更的字段，保留锁定相关状态
	existing.FailCount = record.FailCount
	existing.Locked = record.Locked
	existing.LockedAt = record.LockedAt
	existing.Reason = record.Reason
	return global.DB.Save(&existing).Error
}

func (r *LoginLockoutRepo) ResetByTarget(target, lockType string) error {
	return global.DB.Model(&model.LoginLockout{}).
		Where("`target` = ? AND `type` = ?", target, lockType).
		Updates(map[string]interface{}{
			"fail_count":   0,
			"locked":       false,
			"locked_at":    nil,
			"unlocked_at":  time.Now(),
		}).Error
}

func (r *LoginLockoutRepo) Unlock(target, lockType string) error {
	return global.DB.Model(&model.LoginLockout{}).
		Where("`target` = ? AND `type` = ?", target, lockType).
		Updates(map[string]interface{}{
			"locked":       false,
			"locked_at":   nil,
			"unlocked_at": time.Now(),
		}).Error
}

func (r *LoginLockoutRepo) CountLockedUsers() (int64, error) {
	var count int64
	err := global.DB.Model(&model.LoginLockout{}).
		Where("`type` = ? AND locked = ?", "user", true).Count(&count).Error
	return count, err
}

func (r *LoginLockoutRepo) CountLockedIPs() (int64, error) {
	var count int64
	err := global.DB.Model(&model.LoginLockout{}).
		Where("`type` = ? AND locked = ?", "ip", true).Count(&count).Error
	return count, err
}

func (r *LoginLockoutRepo) ListLockedUsers(limit int) ([]model.LoginLockout, error) {
	var records []model.LoginLockout
	err := global.DB.Where("`type` = ? AND locked = ?", "user", true).
		Order("locked_at DESC").Limit(limit).Find(&records).Error
	return records, err
}

func (r *LoginLockoutRepo) ListLockedIPs(limit int) ([]model.LoginLockout, error) {
	var records []model.LoginLockout
	err := global.DB.Where("`type` = ? AND locked = ?", "ip", true).
		Order("locked_at DESC").Limit(limit).Find(&records).Error
	return records, err
}

// DeleteExpiredLockouts 删除已过期的锁定记录
func (r *LoginLockoutRepo) DeleteExpiredLockouts(lockMinutes int) error {
	return global.DB.Where("locked = ? AND locked_at < ?",
		true,
		time.Now().Add(-time.Duration(lockMinutes)*time.Minute),
	).Delete(&model.LoginLockout{}).Error
}
