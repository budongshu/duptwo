package model

import (
	"time"
)

// SyncStation 同步站点表
type SyncStation struct {
	BaseModel
	Name        string `json:"name" gorm:"size:128;not null"`           // 站点名称
	Code        string `json:"code" gorm:"size:64;uniqueIndex"`          // 站点代码
	URL         string `json:"url" gorm:"size:512"`                      // 站点URL
	APIKey      string `json:"-" gorm:"size:256"`                        // API Key (加密存储)
	Status      string `json:"status" gorm:"size:16;default:active"`    // 状态: active/inactive
	Description string `json:"description" gorm:"size:512"`              // 描述
	IsCenter    bool   `json:"isCenter" gorm:"default:false"`            // 是否为中心站点
	LastSyncAt  *time.Time `json:"lastSyncAt" gorm:"size:64"`            // 最后同步时间
	SyncCount   int64  `json:"syncCount" gorm:"default:0"`               // 同步次数
	Remark      string `json:"remark" gorm:"size:512"`                   // 备注
	IsDeleted   bool   `json:"-" gorm:"default:false"`                   // 软删除标记
}

func (SyncStation) TableName() string {
	return "sync_stations"
}

// SyncHistory 同步历史表
type SyncHistory struct {
	BaseModel
	StationID    uint       `json:"stationId" gorm:"index"`              // 站点ID
	Direction    string     `json:"direction" gorm:"size:16"`            // 方向: upload/download
	Status       string     `json:"status" gorm:"size:16"`               // 状态: pending/processing/completed/failed
	TotalRecords int        `json:"totalRecords" gorm:"default:0"`       // 总记录数
	SuccessCount int        `json:"successCount" gorm:"default:0"`       // 成功数
	FailCount    int        `json:"failCount" gorm:"default:0"`           // 失败数
	ConflictCount int       `json:"conflictCount" gorm:"default:0"`      // 冲突数
	StartedAt    *time.Time `json:"startedAt" gorm:"size:64"`            // 开始时间
	CompletedAt  *time.Time `json:"completedAt" gorm:"size:64"`          // 完成时间
	ErrorMsg     string     `json:"errorMsg" gorm:"size:1024"`           // 错误信息
	Remark       string     `json:"remark" gorm:"size:512"`              // 备注

	// 关联
	Station *SyncStation `json:"station,omitempty" gorm:"foreignKey:StationID"`
}

func (SyncHistory) TableName() string {
	return "sync_histories"
}

// SyncDetail 同步详情表
type SyncDetail struct {
	BaseModel
	HistoryID    uint   `json:"historyId" gorm:"index"`                  // 同步历史ID
	SerialNo     string `json:"serialNo" gorm:"size:64;index"`           // 流水号
	ProjectName  string `json:"projectName" gorm:"size:128;index"`        // 项目名称
	Action       string `json:"action" gorm:"size:16"`                   // 操作: create/update/skip/conflict
	Result       string `json:"result" gorm:"size:16"`                   // 结果: success/failed
	ErrorMsg     string `json:"errorMsg" gorm:"size:512"`                // 错误信息
	OldSerialNo  string `json:"oldSerialNo" gorm:"size:64"`               // 原流水号（冲突时）
	NewSerialNo  string `json:"newSerialNo" gorm:"size:64"`              // 新流水号（冲突时）

	// 关联
	History *SyncHistory `json:"history,omitempty" gorm:"foreignKey:HistoryID"`
}

func (SyncDetail) TableName() string {
	return "sync_details"
}
