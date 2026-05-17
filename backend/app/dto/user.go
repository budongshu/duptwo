package dto

import "time"

// ============ 用户相关 ============

// RegisterReq 注册请求
type RegisterReq struct {
	Username string `json:"username" validate:"required,max=64"`
	Password string `json:"password" validate:"required,min=6,max=64"`
	Nickname string `json:"nickname" validate:"max=64"`
	Email    string `json:"email" validate:"max=128"`
}

// UserListReq 用户列表请求
type UserListReq struct {
	Page      int    `form:"page" validate:"min=1"`
	PageSize  int    `form:"pageSize" validate:"min=1,max=100"`
	Keyword   string `form:"keyword"`
	Status    string `form:"status"`
	RoleID    uint   `form:"roleId"`
	GroupID   uint   `form:"groupId"`
	SortField string `form:"sortField"`
	SortOrder string `form:"sortOrder"`
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
	Locked      bool       `json:"locked"`
	LastLoginAt *time.Time `json:"lastLoginAt"`
	LastLoginIP string     `json:"lastLoginIP"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// ============ 角色相关 ============

// ProfileUpdateReq 个人资料更新请求
type ProfileUpdateReq struct {
	Nickname string `json:"nickname" validate:"max=64"`
	Email    string `json:"email" validate:"max=128"`
	Phone    string `json:"phone" validate:"max=32"`
}

// RoleCreateReq 创建角色请求
type RoleCreateReq struct {
	Name        string   `json:"name" validate:"required,max=64"`
	Code        string   `json:"code" validate:"required,max=64"`
	Description string   `json:"description" validate:"max=256"`
	Permissions []string `json:"permissions"`
	Sort        int      `json:"sort"`
}

// RoleUpdateReq 更新角色请求
type RoleUpdateReq struct {
	ID          uint     `json:"id" validate:"required"`
	Name        string   `json:"name" validate:"required,max=64"`
	Code        string   `json:"code" validate:"required,max=64"`
	Description string   `json:"description" validate:"max=256"`
	Permissions []string `json:"permissions"`
	Sort        int      `json:"sort"`
}

// RoleListReq 角色列表请求
type RoleListReq struct {
	Page     int    `form:"page" validate:"min=1"`
	PageSize int    `form:"pageSize" validate:"min=1,max=100"`
	Keyword  string `form:"keyword"`
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

// ============ 用户组相关 ============

// UserGroupCreateReq 创建用户组请求
type UserGroupCreateReq struct {
	Name        string `json:"name" validate:"required,max=64"`
	Code        string `json:"code" validate:"required,max=64"`
	Description string `json:"description" validate:"max=256"`
	RoleID      uint   `json:"roleId"`
	Sort        int    `json:"sort"`
}

// UserGroupUpdateReq 更新用户组请求
type UserGroupUpdateReq struct {
	ID          uint   `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required,max=64"`
	Code        string `json:"code" validate:"required,max=64"`
	Description string `json:"description" validate:"max=256"`
	RoleID      uint   `json:"roleId"`
	Sort        int    `json:"sort"`
}

// UserGroupListReq 用户组列表请求
type UserGroupListReq struct {
	Page     int    `form:"page" validate:"min=1"`
	PageSize int    `form:"pageSize" validate:"min=1,max=100"`
	Keyword  string `form:"keyword"`
}

// UserGroupResp 用户组响应
type UserGroupResp struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	RoleID      uint      `json:"roleId"`
	RoleName    string    `json:"roleName"`
	Sort        int       `json:"sort"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ============ 管理操作相关 ============

// ResetPasswordReq 重置密码请求
type ResetPasswordReq struct {
	UserID      uint   `json:"userId" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=6,max=128"`
}

// ResetMFAReq 重置MFA请求
type ResetMFAReq struct {
	UserID uint `json:"userId" validate:"required"`
}

// AdminEnableMFAReq 管理员启用MFA请求
type AdminEnableMFAReq struct {
	UserID uint   `json:"userId" validate:"required"`
	Code   string `json:"code" validate:"required,len=6"`
}

// GenerateMFASecretResp 生成MFA密钥响应
type GenerateMFASecretResp struct {
	Secret string `json:"secret"`
	QRCode string `json:"qrCode"` // otpauth:// URL
}

// BatchUpdateRoleReq 批量更新用户角色请求
type BatchUpdateRoleReq struct {
	IDs    []uint `json:"ids" validate:"required,min=1"`
	RoleID uint   `json:"roleId" validate:"required"`
}
