package repo

import (
	"gorm.io/gorm"
)

// DBOption 数据库查询选项
type DBOption func(*gorm.DB) *gorm.DB

// WithByID 根据ID查询
func WithByID(id uint) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

// WithByUUID 根据UUID查询
func WithByUUID(uuid string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", uuid)
	}
}

// WithByName 根据名称查询
func WithByName(name string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", name)
	}
}

// WithByCode 根据编码查询
func WithByCode(code string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("code = ?", code)
	}
}

// WithStatus 根据状态查询
func WithStatus(status string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", status)
	}
}

// WithKeyword 根据关键词模糊查询
func WithKeyword(keyword string, fields ...string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		if keyword == "" {
			return db
		}
		query := ""
		for i, field := range fields {
			if i > 0 {
				query += " OR "
			}
			query += field + " LIKE ?"
		}
		return db.Where("("+query+")", "%"+keyword+"%", "%"+keyword+"%")
	}
}

// WithDateRange 日期范围查询
func WithDateRange(field string, startDate, endDate string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		if startDate != "" {
			db = db.Where(field+" >= ?", startDate)
		}
		if endDate != "" {
			db = db.Where(field+" <= ?", endDate+" 23:59:59")
		}
		return db
	}
}

// WithOrderDesc 降序排序
func WithOrderDesc(orderBy string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(orderBy + " DESC")
	}
}

// WithOrderAsc 升序排序
func WithOrderAsc(orderBy string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(orderBy + " ASC")
	}
}

// WithPagination 分页
func WithPagination(page, pageSize int) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// WithNotDeleted 未删除的记录
func WithNotDeleted() DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("is_deleted = ?", false)
	}
}

// WithRoleID 根据角色ID查询
func WithRoleID(roleID uint) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("role_id = ?", roleID)
	}
}

// WithFormID 根据表单ID查询
func WithFormID(formID uint) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("form_id = ?", formID)
	}
}

// WithFormCode 根据表单编码查询
func WithFormCode(formCode string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("form_code = ?", formCode)
	}
}

// WithRecordID 根据记录ID查询
func WithRecordID(recordID uint) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("record_id = ?", recordID)
	}
}

// Count 计数
func Count(db *gorm.DB) (int64, error) {
	var total int64
	err := db.Count(&total).Error
	return total, err
}

// Paginate 分页查询
func Paginate(page, pageSize int) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
