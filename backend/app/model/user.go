package model

import (
	"time"
)

// User 用户
type User struct {
	BaseModel
	Username    string `json:"username" gorm:"size:64;not null"`
	Password    string `json:"-" gorm:"size:256"` // bcrypt hash (LOCAL user only)
	Nickname    string `json:"nickname" gorm:"size:64"`
	Email       string `json:"email" gorm:"size:128;index"`
	Phone       string `json:"phone" gorm:"size:32"`
	Avatar      string `json:"avatar" gorm:"size:512"`
	Status      string `json:"status" gorm:"size:16;default:active"` // active/inactive
	RoleID      uint   `json:"roleId" gorm:"index"`
	GroupID     uint   `json:"groupId" gorm:"index"`
	LastLoginAt *time.Time `json:"lastLoginAt" gorm:"size:64"`
	LastLoginIP string `json:"lastLoginIP" gorm:"size:64"`
	IsDeleted   bool   `json:"-" gorm:"default:false"`
	// 来源：LOCAL=本地用户，AD=Active Directory用户
	Source string `json:"source" gorm:"size:16;default:LOCAL"`
	// AD用户的DN，用于AD认证后关联
	ADDN    string `json:"-" gorm:"size:512"`
	// AD 扩展字段
	Department string `json:"department" gorm:"size:128"` // 部门
	Title      string `json:"title" gorm:"size:128"`     // 职位/职称
	Company    string `json:"company" gorm:"size:128"`    // 公司
	// MFA 相关
	MFAEnabled  bool   `json:"mfaEnabled" gorm:"default:false"`
	MFASecret   string `json:"-" gorm:"size:128"` // TOTP secret, not exposed to client
}

func (User) TableName() string {
	return "users"
}

// UserResp 用户响应
type UserResp struct {
	ID          uint       `json:"id"`
	Username    string     `json:"username"`
	Nickname    string     `json:"nickname"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	Avatar      string     `json:"avatar"`
	Status      string     `json:"status"`
	StatusText  string     `json:"statusText"`
	RoleID      uint       `json:"roleId"`
	RoleName    string     `json:"roleName"`
	GroupID     uint       `json:"groupId"`
	GroupName   string     `json:"groupName"`
	Department  string     `json:"department"`
	Title       string     `json:"title"`
	Company     string     `json:"company"`
	Source      string     `json:"source"` // LOCAL / AD
	MFAEnabled  bool       `json:"mfaEnabled"`
	LastLoginAt *time.Time `json:"lastLoginAt"`
	LastLoginIP string     `json:"lastLoginIP"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}
