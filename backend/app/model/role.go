package model

import (
	"time"
)

// Role 角色
type Role struct {
	BaseModel
	Name        string `json:"name" gorm:"size:64;not null"`
	Code        string `json:"code" gorm:"size:64;uniqueIndex;not null"`
	Description string `json:"description" gorm:"size:256"`
	Permissions string `json:"permissions" gorm:"type:text"` // JSON array of permission codes
	IsDeleted   bool   `json:"-" gorm:"default:false"`
	Sort        int    `json:"sort" gorm:"default:0"`
}

func (Role) TableName() string {
	return "roles"
}

// RoleResp 角色响应
type RoleResp struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Permissions []string  `json:"permissions"`
	Sort        int       `json:"sort"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// UserGroup 用户组
type UserGroup struct {
	BaseModel
	Name        string `json:"name" gorm:"size:64;not null"`
	Code        string `json:"code" gorm:"size:64;uniqueIndex;not null"`
	Description string `json:"description" gorm:"size:256"`
	IsDeleted   bool   `json:"-" gorm:"default:false"`
	Sort        int    `json:"sort" gorm:"default:0"`
}

func (UserGroup) TableName() string {
	return "user_groups"
}

// UserGroupResp 用户组响应
type UserGroupResp struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Sort        int       `json:"sort"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
