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
	// 存在，更新
	record.ID = existing.ID
	record.CreatedAt = existing.CreatedAt
	return global.DB.Save(record).Error
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
