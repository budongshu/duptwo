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
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Code            string `json:"code"`
	URL             string `json:"url"`
	Status          string `json:"status"`
	StatusText      string `json:"statusText"`
	Description     string `json:"description"`
	IsCenter        bool   `json:"isCenter"`
	APIKey          string `json:"apiKey,omitempty"`        // 创建时返回API Key
	LastSyncAt      int64  `json:"lastSyncAt"`             // Unix ms
	SyncCount       int64  `json:"syncCount"`
	Remark          string `json:"remark"`
	CreatedAt       int64  `json:"createdAt"`   // Unix ms
	UpdatedAt       int64  `json:"updatedAt"`   // Unix ms
	LastHeartbeatAt int64  `json:"lastHeartbeatAt"` // Unix ms
	LastConnectedAt int64  `json:"lastConnectedAt"` // Unix ms，最后探测成功时间
	IsConnected     bool   `json:"isConnected"`      // 主动探测结果
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
	StationCode string `json:"stationCode" validate:"required,max=64"`   // 站点代码
	StationName string `json:"stationName" validate:"required,max=128"`  // 站点名称
	URL         string `json:"url" validate:"required,max=512"`          // 站点URL
	Password    string `json:"password"`                                 // 站点密码（可选，用于安全验证）
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
	Keyword    string `form:"keyword"` // 搜索流水号/备注
}

// SyncStationSummaryResp 站点同步汇总（用于 Center 首页卡片）
type SyncStationSummaryResp struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Code          string `json:"code"`
	Status        string `json:"status"`
	LastSyncAt    int64  `json:"lastSyncAt"`    // Unix ms，DB存为INTEGER
	TotalSyncs    int    `json:"totalSyncs"`    // 历史总同步次数
	TotalRecords  int    `json:"totalRecords"`   // 历史累计记录数
	SuccessCount  int    `json:"successCount"`   // 历史累计成功数
	FailCount     int    `json:"failCount"`      // 历史累计失败数
	ConflictCount int    `json:"conflictCount"`  // 历史累计冲突数
}

// SyncStationSummaryReq 站点汇总查询请求
type SyncStationSummaryReq struct {
	StationID uint   `form:"stationId"`
	StartDate string `form:"startDate"`
	EndDate   string `form:"endDate"`
}

// SyncHistoryResp 同步历史响应
type SyncHistoryResp struct {
	ID            uint   `json:"id"`
	StationID     uint   `json:"stationId"`
	StationName   string `json:"stationName"`
	StationCode   string `json:"stationCode"`
	Direction     string `json:"direction"`
	DirectionText string `json:"directionText"`
	Status        string `json:"status"`
	StatusText    string `json:"statusText"`
	TotalRecords  int    `json:"totalRecords"`
	SuccessCount  int    `json:"successCount"`
	FailCount     int    `json:"failCount"`
	ConflictCount int    `json:"conflictCount"`
	StartedAt     int64  `json:"startedAt"`    // Unix ms，DB存为string
	CompletedAt   int64  `json:"completedAt"`  // Unix ms
	ErrorMsg      string `json:"errorMsg"`
	Remark        string `json:"remark"`
	CreatedAt     int64  `json:"createdAt"`   // Unix ms
}

// SyncHistoryDetailResp 同步历史详情响应
type SyncHistoryDetailResp struct {
	SyncHistoryResp
	Total   int64             `json:"total"`    // 详情总条数（分页用）
	Details []SyncDetailResp   `json:"details"`
}

// SyncDetailResp 同步详情响应
type SyncDetailResp struct {
	ID          uint   `json:"id"`
	HistoryID   uint   `json:"historyId"`
	SerialNo    string `json:"serialNo"`
	ProjectName string `json:"projectName"`
	Action      string `json:"action"`
	ActionText  string `json:"actionText"`
	Result      string `json:"result"`
	ErrorMsg    string `json:"errorMsg"`
	OldSerialNo string `json:"oldSerialNo"`
	NewSerialNo string `json:"newSerialNo"`
	CreatedAt   int64  `json:"createdAt"` // unix timestamp 毫秒
}

// ============ 同步状态相关 ============

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

// SyncStatusResp 同步状态响应（增强版，包含调度器完整状态）
type SyncStatusResp struct {
	Enabled        bool       `json:"enabled"`
	Mode           string     `json:"mode"`            // center/agent
	IsCenter       bool       `json:"isCenter"`
	StationID      string     `json:"stationId"`
	StationName    string     `json:"stationName"`
	CenterURL      string     `json:"centerUrl"`
	Registered     bool       `json:"registered"`        // Agent是否已注册到Center
	Interval          string     `json:"interval"`          // 同步间隔
	HeartbeatInterval string     `json:"heartbeatInterval"` // 心跳间隔
	HeartbeatTimeout int        `json:"heartbeatTimeout"`  // 心跳超时（秒）
	BatchSize        int        `json:"batchSize"`        // 每批数量
	Filter         *SyncFilterResp `json:"filter"`      // 同步过滤器
	LastSyncAt     *time.Time `json:"lastSyncAt"`       // 上次同步时间
	LastSerialNo   string     `json:"lastSerialNo"`     // 上次同步最后SerialNo
	SyncQueueCount int        `json:"syncQueueCount"`   // 队列中任务数
	QueueTotal     int        `json:"queueTotal"`       // 队列总任务数
	QueuePending   int        `json:"queuePending"`     // 等待中
	QueueCompleted int        `json:"queueCompleted"`    // 已完成
	QueueFailed    int        `json:"queueFailed"`      // 失败
	LastErrorAt    *time.Time `json:"lastErrorAt"`      // 最后错误时间
	LastError      string     `json:"lastError"`        // 最后错误信息
}

// SyncFilterResp 同步过滤器响应
type SyncFilterResp struct {
	ProjectNames []string `json:"projectNames"` // 只同步这些项目
}

// SyncResetKeyResp 重置API Key响应
type SyncResetKeyResp struct {
	ID      uint   `json:"id"`
	APIKey  string `json:"apiKey"` // 新API Key（仅返回一次）
}
