package repo

import (
	"datauptwo/app/model"
	"datauptwo/global"
)

type OperationLogRepo struct{}

func NewOperationLogRepo() *OperationLogRepo {
	return &OperationLogRepo{}
}

func (r *OperationLogRepo) Create(log *model.OperationLog) error {
	return global.DB.Create(log).Error
}

func (r *OperationLogRepo) List(page, pageSize int, userID uint, username, menuName, action, resourceType, keyword, startDate, endDate string) ([]model.OperationLog, int64, error) {
	var logs []model.OperationLog
	var total int64

	query := global.DB.Model(&model.OperationLog{})

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if menuName != "" {
		query = query.Where("menu_name LIKE ?", "%"+menuName+"%")
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}
	if keyword != "" {
		query = query.Where("resource_name LIKE ? OR detail LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate+" 23:59:59")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error
	return logs, total, err
}

const exportLimit = 10000

func (r *OperationLogRepo) ListForExport(userID uint, username, menuName, action, resourceType, keyword, startDate, endDate string) ([]model.OperationLog, int64, bool, error) {
	var logs []model.OperationLog

	query := global.DB.Model(&model.OperationLog{})

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if menuName != "" {
		query = query.Where("menu_name LIKE ?", "%"+menuName+"%")
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}
	if keyword != "" {
		query = query.Where("resource_name LIKE ? OR detail LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate+" 23:59:59")
	}

	// 先查实际总数
	var total int64
	query.Count(&total)

	truncated := total > exportLimit
	if truncated {
		err := query.Order("created_at DESC").Limit(exportLimit).Find(&logs).Error
		return logs, total, truncated, err
	}

	err := query.Order("created_at DESC").Find(&logs).Error
	return logs, total, false, err
}

type LoginLogRepo struct{}

func NewLoginLogRepo() *LoginLogRepo {
	return &LoginLogRepo{}
}

func (r *LoginLogRepo) Create(log *model.LoginLog) error {
	return global.DB.Create(log).Error
}

func (r *LoginLogRepo) List(page, pageSize int, userID uint, username, status, keyword, startDate, endDate string) ([]model.LoginLog, int64, error) {
	var logs []model.LoginLog
	var total int64

	query := global.DB.Model(&model.LoginLog{})

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("fail_reason LIKE ? OR ip_address LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate+" 23:59:59")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error
	return logs, total, err
}

func (r *LoginLogRepo) ListForExport(userID uint, username, status, keyword, startDate, endDate string) ([]model.LoginLog, int64, bool, error) {
	var logs []model.LoginLog

	query := global.DB.Model(&model.LoginLog{})

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("fail_reason LIKE ? OR ip_address LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate+" 23:59:59")
	}

	var total int64
	query.Count(&total)

	truncated := total > exportLimit
	if truncated {
		err := query.Order("created_at DESC").Limit(exportLimit).Find(&logs).Error
		return logs, total, truncated, err
	}

	err := query.Order("created_at DESC").Find(&logs).Error
	return logs, total, false, err
}

func (r *LoginLogRepo) GetRecentByUser(userID uint, limit int) ([]model.LoginLog, error) {
	var logs []model.LoginLog
	err := global.DB.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Find(&logs).Error
	return logs, err
}
