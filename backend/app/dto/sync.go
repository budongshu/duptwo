package dto

import "time"

// ============ 同步站点相关 ============

// SyncStationCreateReq 创建站点请求
type SyncStationCreateReq struct {
	Name        string `json:"name" validate:"required,max=128"`
	Code        string `json:"code" validate:"required,max=64"`
	URL         string `json:"url" validate:"required,max=512"`
	Description string `json:"description" validate:"max=512"`
	IsCenter    bool   `json:"isCenter"`
}

// SyncStationUpdateReq 更新站点请求
type SyncStationUpdateReq struct {
	ID          uint   `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required,max=128"`
	Code        string `json:"code" validate:"required,max=64"`
	URL         string `json:"url" validate:"required,max=512"`
	Status      string `json:"status" validate:"omitempty,oneof=active inactive"`
	Description string `json:"description" validate:"max=512"`
	IsCenter    bool   `json:"isCenter"`
}

// SyncStationResp 站点响应
type SyncStationResp struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	URL         string     `json:"url"`
	Status      string     `json:"status"`
	StatusText  string     `json:"statusText"`
	Description string     `json:"description"`
	IsCenter    bool       `json:"isCenter"`
	APIKey      string     `json:"apiKey,omitempty"` // 创建时返回API Key
	LastSyncAt  *time.Time `json:"lastSyncAt"`
	SyncCount   int64      `json:"syncCount"`
	Remark      string     `json:"remark"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// SyncStationListReq 站点列表请求
type SyncStationListReq struct {
	Page     int    `form:"page" validate:"min=1"`
	PageSize int    `form:"pageSize" validate:"min=1,max=100"`
	Keyword  string `form:"keyword"`
	Status   string `form:"status"`
}

// ============ 同步注册相关 ============

// SyncRegisterReq 站点注册请求
type SyncRegisterReq struct {
	StationCode string `json:"stationCode" validate:"required,max=64"` // 站点代码
	StationName string `json:"stationName" validate:"required,max=128"` // 站点名称
	URL         string `json:"url" validate:"required,max=512"`         // 站点URL
	Password    string `json:"password" validate:"required,max=128"`    // 站点密码
}

// SyncRegisterResp 站点注册响应
type SyncRegisterResp struct {
	StationID  uint   `json:"stationId"`
	APIKey     string `json:"apiKey"`      // 新生成的API Key
	APIKeyHash string `json:"-"`          // API Key 哈希（存储用）
	Message    string `json:"message"`
}

// ============ 同步记录上传相关 ============

// SyncUploadReq 上传记录同步请求
type SyncUploadReq struct {
	Records []SyncRecordItem `json:"records" validate:"required,min=1"`
}

// SyncRecordItem 同步记录项
type SyncRecordItem struct {
	SerialNo    string                 `json:"serialNo" validate:"required,max=64"`    // 流水号
	ProjectName string                 `json:"projectName" validate:"required,max=128"`  // 项目名称
	DiskLabel   string                 `json:"diskLabel" validate:"max=64"`             // 磁盘标签
	DestPath    string                 `json:"destPath" validate:"max=512"`             // 目标路径
	FileSize    int64                  `json:"fileSize"`                               // 文件大小
	Uploader    string                 `json:"uploader" validate:"max=64"`             // 上传人
	Status      string                 `json:"status" validate:"max=16"`               // 状态
	Remark      string                 `json:"remark" validate:"max=512"`              // 备注
	Data        map[string]interface{} `json:"data"`                                  // 动态字段数据
}

// SyncUploadResp 同步响应
type SyncUploadResp struct {
	TotalRecords   int                `json:"totalRecords"`   // 总记录数
	SuccessCount   int                `json:"successCount"`   // 成功数
	FailCount      int                `json:"failCount"`      // 失败数
	ConflictCount  int                `json:"conflictCount"` // 冲突数
	Details        []SyncDetailItem   `json:"details"`        // 每条记录的处理结果
	HistoryID      uint               `json:"historyId"`      // 同步历史ID
}

// SyncDetailItem 单条同步详情
type SyncDetailItem struct {
	SerialNo    string `json:"serialNo"`     // 原始流水号
	ProjectName string `json:"projectName"`  // 项目名称
	Action      string `json:"action"`       // 操作: create/update/skip/conflict
	Result      string `json:"result"`       // 结果: success/failed
	ErrorMsg    string `json:"errorMsg"`     // 错误信息
	NewSerialNo string `json:"newSerialNo"` // 新流水号（冲突时返回）
	OldSerialNo string `json:"oldSerialNo"` // 原流水号（冲突时）
}

// ============ 同步历史相关 ============

// SyncHistoryListReq 同步历史列表请求
type SyncHistoryListReq struct {
	Page       int    `form:"page" validate:"min=1"`
	PageSize   int    `form:"pageSize" validate:"min=1,max=100"`
	StationID  uint   `form:"stationId"`
	Direction  string `form:"direction"`
	Status     string `form:"status"`
	StartDate  string `form:"startDate"`
	EndDate    string `form:"endDate"`
}

// SyncHistoryResp 同步历史响应
type SyncHistoryResp struct {
	ID           uint        `json:"id"`
	StationID    uint        `json:"stationId"`
	StationName  string      `json:"stationName"`
	Direction    string      `json:"direction"`
	DirectionText string     `json:"directionText"`
	Status       string      `json:"status"`
	StatusText   string      `json:"statusText"`
	TotalRecords int         `json:"totalRecords"`
	SuccessCount int         `json:"successCount"`
	FailCount    int         `json:"failCount"`
	ConflictCount int        `json:"conflictCount"`
	StartedAt    *time.Time  `json:"startedAt"`
	CompletedAt  *time.Time  `json:"completedAt"`
	ErrorMsg     string      `json:"errorMsg"`
	Remark       string      `json:"remark"`
	CreatedAt    time.Time   `json:"createdAt"`
}

// SyncHistoryDetailResp 同步历史详情响应
type SyncHistoryDetailResp struct {
	SyncHistoryResp
	Details []SyncDetailResp `json:"details"`
}

// SyncDetailResp 同步详情响应
type SyncDetailResp struct {
	ID          uint   `json:"id"`
	SerialNo     string `json:"serialNo"`
	ProjectName  string `json:"projectName"`
	Action       string `json:"action"`
	ActionText   string `json:"actionText"`
	Result       string `json:"result"`
	ErrorMsg     string `json:"errorMsg"`
	OldSerialNo  string `json:"oldSerialNo"`
	NewSerialNo  string `json:"newSerialNo"`
	CreatedAt   time.Time `json:"createdAt"`
}

// ============ 同步状态相关 ============

// SyncStatusResp 同步状态响应
type SyncStatusResp struct {
	Enabled     bool    `json:"enabled"`
	Mode        string  `json:"mode"`         // center/agent
	IsCenter    bool    `json:"isCenter"`
	StationID   string  `json:"stationId"`
	StationName string  `json:"stationName"`
	CenterURL   string  `json:"centerUrl"`
	LastSyncAt  *time.Time `json:"lastSyncAt"`
	SyncQueueCount int   `json:"syncQueueCount"` // 队列中的同步任务数
}

// SyncQueueItem 同步队列项
type SyncQueueItem struct {
	ID          uint      `json:"id"`
	SerialNo    string    `json:"serialNo"`
	ProjectName string    `json:"projectName"`
	Status      string    `json:"status"`
	RetryCount  int       `json:"retryCount"`
	NextRetryAt *time.Time `json:"nextRetryAt"`
	CreatedAt   time.Time `json:"createdAt"`
	ErrorMsg    string    `json:"errorMsg"`
}
