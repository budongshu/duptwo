package model

import (
	"time"
)

// FieldConfig 字段配置
type FieldConfig struct {
	BaseModel
	Name        string `json:"name" gorm:"size:64;not null"`       // 字段名称
	Code        string `json:"code" gorm:"size:64;not null;uniqueIndex"` // 字段编码
	Type        string `json:"type" gorm:"size:32;default:text"`    // 字段类型: text/number/select/date/datetime
	Required    bool   `json:"required" gorm:"default:false"`       // 是否必填
	Options     string `json:"options" gorm:"type:text"`            // 选项(JSON数组,select类型用)
	DefaultValue string `json:"defaultValue" gorm:"size:512"`       // 默认值
	Placeholder string `json:"placeholder" gorm:"size:256"`        // 占位提示
	Sort        int    `json:"sort" gorm:"default:0"`               // 排序
	Enabled     bool   `json:"enabled" gorm:"default:true"`        // 是否启用
	IsDeleted   bool   `json:"-" gorm:"default:false"`             // 软删除
}

func (FieldConfig) TableName() string {
	return "field_configs"
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
