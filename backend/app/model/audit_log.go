package model

import (
	"time"
)

// OperationLog 操作日志
type OperationLog struct {
	BaseModel
	UserID       uint      `json:"userId" gorm:"index"`
	Username     string    `json:"username" gorm:"size:64;index"`
	MenuName     string    `json:"menuName" gorm:"size:64"`      // 功能菜单名称
	Action       string    `json:"action" gorm:"size:32"`        // action: view, create, update, delete, export
	ResourceType string    `json:"resourceType" gorm:"size:64"` // 资源类型：如 User, Role, UploadRecord
	ResourceID   uint      `json:"resourceId" gorm:"index"`     // 资源ID
	ResourceName string    `json:"resourceName" gorm:"size:256"` // 资源名称/描述
	IPAddress    string    `json:"ipAddress" gorm:"size:64"`    // IP地址
	UserAgent    string    `json:"userAgent" gorm:"size:512"`   // 浏览器UA
	Detail       string    `json:"detail" gorm:"type:text"`     // 详细信息JSON
	CreatedAt    time.Time `json:"createdAt" gorm:"index"`
}

func (OperationLog) TableName() string {
	return "operation_logs"
}

// LoginLog 登录日志
type LoginLog struct {
	BaseModel
	UserID       uint      `json:"userId" gorm:"index"`
	Username     string    `json:"username" gorm:"size:64;index"`
	Status       string    `json:"status" gorm:"size:16"`          // success, failed
	IPAddress    string    `json:"ipAddress" gorm:"size:64"`       // IP地址
	UserAgent    string    `json:"userAgent" gorm:"size:512"`      // 浏览器UA
	FailReason   string    `json:"failReason" gorm:"size:256"`     // 失败原因
	MFAUsed      bool      `json:"mfaUsed" gorm:"default:false"`   // 是否使用了MFA
	LoginMethod  string    `json:"loginMethod" gorm:"size:32"`     // login_method: password, mfa
	CreatedAt    time.Time `json:"createdAt" gorm:"index"`
}

func (LoginLog) TableName() string {
	return "login_logs"
}
