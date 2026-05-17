package model

import (
	"time"
)

// UploadRecord 上传记录
type UploadRecord struct {
	BaseModel
	SerialNo    string  `json:"serialNo" gorm:"size:64;index"`                    // 流水号
	DiskLabel   string  `json:"diskLabel" gorm:"column:disk_label;size:64;index"`  // 磁盘标签
	ProjectID   *uint   `json:"projectId" gorm:"index"`                            // 项目ID（外键， nullable）
	ProjectName string  `json:"projectName" gorm:"size:128;index"`                 // 项目名称（保留用于显示/兼容性）
	DestPath    string  `json:"destPath" gorm:"size:512"`                          // 目标路径
	FileSize    int64   `json:"fileSize" gorm:"default:0"`                        // 文件大小(bytes)
	Uploader    string  `json:"uploader" gorm:"size:64;index"`                    // 上传人
	Status      string  `json:"status" gorm:"size:16;default:pending"`            // 状态: pending/processing/completed/failed
	Remark      string  `json:"remark" gorm:"size:512"`                           // 备注
	Data        string  `json:"data" gorm:"type:text"`                             // 动态字段数据(JSON)
	IsDeleted   bool    `json:"-" gorm:"default:false"`                           // 软删除标记

	// 关联项目（可选，用于预加载）
	Project *Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
}

func (UploadRecord) TableName() string {
	return "upload_records"
}

// UploadRecordResp 上传记录响应
type UploadRecordResp struct {
	ID          uint      `json:"id"`
	SerialNo    string    `json:"serialNo"`     // 流水号
	DiskLabel   string    `json:"diskLabel"`    // 磁盘标签
	ProjectID   *uint     `json:"projectId"`     // 项目ID
	ProjectName string    `json:"projectName"`  // 项目名称
	DestPath    string    `json:"destPath"`
	FileSize    int64     `json:"fileSize"`
	FileSizeStr string    `json:"fileSizeStr"` // 格式化大小
	Uploader    string    `json:"uploader"`
	Status      string    `json:"status"`
	StatusText  string    `json:"statusText"`
	Remark      string    `json:"remark"`
	Data        string    `json:"data"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
