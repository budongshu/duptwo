package dto

import "time"

// ============ 字段配置相关 ============

// FieldConfigCreateReq 创建字段配置请求
type FieldConfigCreateReq struct {
	Name        string   `json:"name" validate:"required,max=64"`
	Code        string   `json:"code" validate:"required,max=64,alphanum"`
	Type        string   `json:"type" validate:"required,oneof=text number select date datetime"`
	Required    bool     `json:"required"`
	Options     []string `json:"options"`
	DefaultValue string  `json:"defaultValue" validate:"max=512"`
	Placeholder string   `json:"placeholder" validate:"max=256"`
	Sort        int      `json:"sort"`
	Enabled     bool     `json:"enabled"`
}

// FieldConfigUpdateReq 更新字段配置请求
type FieldConfigUpdateReq struct {
	ID          uint     `json:"id" validate:"required"`
	Name        string   `json:"name" validate:"required,max=64"`
	Code        string   `json:"code" validate:"required,max=64"`
	Type        string   `json:"type" validate:"required,oneof=text number select date datetime"`
	Required    bool     `json:"required"`
	Options     []string `json:"options"`
	DefaultValue string  `json:"defaultValue" validate:"max=512"`
	Placeholder string   `json:"placeholder" validate:"max=256"`
	Sort        int      `json:"sort"`
	Enabled     bool     `json:"enabled"`
}

// FieldConfigListReq 字段配置列表请求
type FieldConfigListReq struct {
	Page     int    `form:"page" validate:"min=1"`
	PageSize int    `form:"pageSize" validate:"min=1,max=100"`
	Keyword  string `form:"keyword"`
	Enabled  *bool  `form:"enabled"`
}

// FieldConfigResp 字段配置响应
type FieldConfigResp struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Type        string    `json:"type"`
	Required    bool      `json:"required"`
	Options     []string  `json:"options"`
	DefaultValue string   `json:"defaultValue"`
	Placeholder string    `json:"placeholder"`
	Sort        int       `json:"sort"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
