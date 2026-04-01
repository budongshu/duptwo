package dto

import "time"

// ============ 通用响应结构 ============

// Response API统一响应
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageResult 分页结果
type PageResult struct {
	Total int64       `json:"total"`
	Items interface{} `json:"items"`
}

// PageReq 分页请求
type PageReq struct {
	Page     int    `json:"page" form:"page" validate:"required,min=1"`
	PageSize int    `json:"pageSize" form:"pageSize" validate:"required,min=1,max=100"`
	Keyword  string `json:"keyword" form:"keyword"`
}

// ============ 用户相关 ============

// LoginReq 登录请求
type LoginReq struct {
	Username  string `json:"username" validate:"required,min=3,max=32"`
	Password  string `json:"password" validate:"required,min=6,max=128"`
	CaptchaID string `json:"captchaId"`
	Captcha   string `json:"captcha"`
}

// ADLoginReq AD登录请求
type ADLoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	CaptchaID string `json:"captchaId"`
	Captcha   string `json:"captcha"`
}

// LoginResp 登录响应
type LoginResp struct {
	Token        string    `json:"token"`
	ExpireAt     int64     `json:"expireAt"`
	User         UserInfo  `json:"user"`
	MFARequired  bool      `json:"mfaRequired"` // 是否需要MFA验证
	CaptchaEnabled bool    `json:"captchaEnabled"` // 是否需要验证码
}

// MFAVerifyReq MFA验证请求
type MFAVerifyReq struct {
	UserID   uint   `json:"userId" validate:"required"`
	Code     string `json:"code" validate:"required,len=6"`
	TmpToken string `json:"tmpToken"` // 临时token，用于MFA验证
}

// MFAEnableReq 启用MFA请求
type MFAEnableReq struct {
	Code string `json:"code" validate:"required,len=6"`
}

// MFAStatusResp MFA状态响应
type MFAStatusResp struct {
	Enabled    bool   `json:"enabled"`
	Secret     string `json:"secret,omitempty"` // 启用时返回密钥用于生成二维码
	BackupCodes []string `json:"backupCodes,omitempty"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID          uint     `json:"id"`
	Username    string   `json:"username"`
	Nickname    string   `json:"nickname"`
	Email       string   `json:"email"`
	Phone       string   `json:"phone"`
	Avatar      string   `json:"avatar"`
	RoleID      uint     `json:"roleId"`
	RoleName    string   `json:"roleName"`
	RoleCode    string   `json:"roleCode"`
	Status      string   `json:"status"`
	MFAEnabled  bool     `json:"mfaEnabled"`
	Source      string   `json:"source"` // LOCAL / AD
	Permissions []string `json:"permissions"` // 用户权限列表，用于前端菜单控制
}

// UserCreateReq 创建用户请求
type UserCreateReq struct {
	Username   string `json:"username" validate:"required,min=3,max=32"`
	Password   string `json:"password" validate:"required,min=6,max=128"`
	Nickname   string `json:"nickname" validate:"max=64"`
	Email      string `json:"email" validate:"email"`
	Phone      string `json:"phone" validate:"max=32"`
	RoleID     uint   `json:"roleId"`
	GroupID    uint   `json:"groupId"`
	Status     string `json:"status"`
	MFAEnabled bool   `json:"mfaEnabled"`
}

// UserUpdateReq 更新用户请求
type UserUpdateReq struct {
	ID         uint   `json:"id" validate:"required"`
	Nickname   string `json:"nickname" validate:"max=64"`
	Email      string `json:"email" validate:"email"`
	Phone      string `json:"phone" validate:"max=32"`
	RoleID     uint   `json:"roleId"`
	GroupID    uint   `json:"groupId"`
	Status     string `json:"status"`
	MFAEnabled bool   `json:"mfaEnabled"`
}

// ChangePasswordReq 修改密码请求
type ChangePasswordReq struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=6,max=128"`
}

// ============ 表单相关 ============

// FormCreateReq 创建表单请求
type FormCreateReq struct {
	Name        string      `json:"name" validate:"required,min=2,max=128"`
	Code        string      `json:"code" validate:"required,min=2,max=64"`
	Description string      `json:"description" validate:"max=512"`
	Fields      []FormField `json:"fields" validate:"required,min=1"`
	Settings    FormSettings `json:"settings"`
}

// FormUpdateReq 更新表单请求
type FormUpdateReq struct {
	ID          uint        `json:"id" validate:"required"`
	Name        string      `json:"name" validate:"required,min=2,max=128"`
	Description string      `json:"description" validate:"max=512"`
	Fields      []FormField `json:"fields" validate:"required,min=1"`
	Settings    FormSettings `json:"settings"`
	ChangeNote  string      `json:"changeNote" validate:"max=256"`
}

// FormPublishReq 发布表单请求
type FormPublishReq struct {
	ID         uint   `json:"id" validate:"required"`
	ChangeNote string `json:"changeNote" validate:"max=256"`
}

// FormField 表单字段
type FormField struct {
	Name       string      `json:"name" validate:"required"`
	Code       string      `json:"code" validate:"required"`
	Type       string      `json:"type" validate:"required,oneof=text number date datetime select multiselect radio checkbox textarea file image"`
	Required   bool        `json:"required"`
	Default    interface{} `json:"default"`
	Options    []string    `json:"options,omitempty"`
	Validation string      `json:"validation,omitempty"`
	Sort       int         `json:"sort"`
	Width      int         `json:"width,omitempty"`
}

// FormSettings 表单设置
type FormSettings struct {
	AllowAnonymous bool `json:"allowAnonymous"`
	AllowEdit      bool `json:"allowEdit"`
	AllowDelete    bool `json:"allowDelete"`
	NotifyOnSubmit bool `json:"notifyOnSubmit"`
	SaveHistories  bool `json:"saveHistories"`
}

// FormResp 表单响应
type FormResp struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Code        string         `json:"code"`
	Description string         `json:"description"`
	Fields      []FormField    `json:"fields"`
	Settings    FormSettings   `json:"settings"`
	Status      string         `json:"status"`
	Version     int            `json:"version"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	RecordCount int64          `json:"recordCount"` // 记录数量
}

// FormListReq 表单列表请求
type FormListReq struct {
	Page     int    `form:"page" validate:"min=1"`
	PageSize int    `form:"pageSize" validate:"min=1,max=100"`
	Keyword  string `form:"keyword"`
	Status   string `form:"status"`
}

// ============ 数据记录相关 ============

// RecordSubmitReq 提交记录请求（公开接口）
type RecordSubmitReq struct {
	FormCode     string                 `json:"formCode" validate:"required"`
	Data         map[string]interface{} `json:"data" validate:"required"`
	Submitter    string                 `json:"submitter" validate:"required"`
	SubmitSource string                 `json:"submitSource"`
}

// RecordCreateReq 创建记录请求
type RecordCreateReq struct {
	FormCode    string                 `json:"formCode" validate:"required"`
	Data        map[string]interface{} `json:"data" validate:"required"`
	Submitter   string                  `json:"submitter"`
	SubmitSource string                 `json:"submitSource"`
}

// RecordUpdateReq 更新记录请求
type RecordUpdateReq struct {
	ID       uint                   `json:"id" validate:"required"`
	FormCode string                  `json:"formCode" validate:"required"`
	Data     map[string]interface{}  `json:"data" validate:"required"`
	ChangeNote string                `json:"changeNote"`
}

// RecordApproveReq 审核记录请求
type RecordApproveReq struct {
	ID       uint   `json:"id" validate:"required"`
	Status   string `json:"status" validate:"required,oneof=approved rejected"`
	ApproveNote string `json:"approveNote"`
}

// RecordStatusUpdateReq 更新记录状态请求
type RecordStatusUpdateReq struct {
	ID       uint   `json:"id" validate:"required"`
	Status   string `json:"status" validate:"required"`
	Approver string `json:"approver"`
}

// RecordBatchStatusReq 批量更新记录状态请求
type RecordBatchStatusReq struct {
	IDs     []uint `json:"ids" validate:"required,min=1"`
	Status  string `json:"status" validate:"required"`
	Message string `json:"message"`
}

// RecordBatchDeleteReq 批量删除记录请求
type RecordBatchDeleteReq struct {
	IDs     []uint `json:"ids" validate:"required,min=1"`
	Message string `json:"message"`
}

// BatchDeleteReq 批量删除请求
type BatchDeleteReq struct {
	IDs []uint `json:"ids" validate:"required,min=1"`
}

// BatchOperationResult 批量操作结果
type BatchOperationResult struct {
	Total     int      `json:"total"`
	Success   int      `json:"success"`
	Failed    int      `json:"failed"`
	FailItems []BatchFailItem `json:"failItems"`
}

// BatchFailItem 批量操作失败项
type BatchFailItem struct {
	ID     uint   `json:"id"`
	SerialNo string `json:"serialNo"`
	Reason string `json:"reason"`
}

// RecordListReq 记录列表请求
type RecordListReq struct {
	Page        int    `form:"page" validate:"min=1"`
	PageSize    int    `form:"pageSize" validate:"min=1,max=100"`
	FormCode    string `form:"formCode"`
	Status      string `form:"status"`
	StartDate   string `form:"startDate"`
	EndDate     string `form:"endDate"`
	Submitter   string `form:"submitter"`
	Keyword     string `form:"keyword"` // 搜索Data字段
}

// RecordResp 记录响应
type RecordResp struct {
	ID           uint                  `json:"id"`
	FormID       uint                  `json:"formID"`
	FormCode     string                `json:"formCode"`
	FormName     string                `json:"formName"`
	SerialNo     string                `json:"serialNo"`
	Data         map[string]interface{} `json:"data"`
	Status       string                `json:"status"`
	Submitter    string                `json:"submitter"`
	SubmitSource string                `json:"submitSource"`
	Approver     string                `json:"approver"`
	ApprovedAt   *time.Time            `json:"approvedAt"`
	ApproveNote  string                `json:"approveNote"`
	Version      int                   `json:"version"`
	CreatedAt    time.Time             `json:"createdAt"`
	UpdatedAt    time.Time             `json:"updatedAt"`
}

// RecordDetailResp 记录详情响应
type RecordDetailResp struct {
	RecordResp
	Histories []RecordHistoryResp `json:"histories,omitempty"`
	Comments  []RecordCommentResp `json:"comments,omitempty"`
}

// RecordHistoryResp 记录历史响应
type RecordHistoryResp struct {
	ID         uint      `json:"id"`
	Version    int       `json:"version"`
	Data       map[string]interface{} `json:"data"`
	ChangedBy  string    `json:"changedBy"`
	ChangeType string    `json:"changeType"`
	ChangeNote string    `json:"changeNote"`
	CreatedAt  time.Time `json:"createdAt"`
}

// RecordCommentResp 记录评论响应
type RecordCommentResp struct {
	ID       uint                  `json:"id"`
	UserID   uint                  `json:"userID"`
	Username string                `json:"username"`
	Content  string                `json:"content"`
	ParentID *uint                 `json:"parentID"`
	Replies  []RecordCommentResp   `json:"replies,omitempty"`
	CreatedAt time.Time           `json:"createdAt"`
}

// ============ 统计报表相关 ============

// StatisticsReq 统计请求
type StatisticsReq struct {
	FormCode    string `form:"formCode"`
	StartDate   string `form:"startDate"`
	EndDate     string `form:"endDate"`
	GroupBy     string `form:"groupBy"` // day/week/month
	ProjectName string `form:"projectName"`
	DataType    string `form:"dataType"`
	Status      string `form:"status"`
	Uploader    string `form:"uploader"`
}

// StatisticsResp 统计响应
type StatisticsResp struct {
	Total      int64                  `json:"total"`
	Pending    int64                  `json:"pending"`
	Approved   int64                  `json:"approved"`
	Rejected   int64                  `json:"rejected"`
	Today      int64                  `json:"today"`
	ThisWeek   int64                  `json:"thisWeek"`
	ThisMonth  int64                  `json:"thisMonth"`
	Trend      []TrendItem            `json:"trend"`
	TopForms   []FormStatItem         `json:"topForms"`
	// 数据传输统计
	DataTransfer *DataTransferStats   `json:"dataTransfer"`
}

// DataTransferStats 数据传输统计
type DataTransferStats struct {
	TodayVolume    float64           `json:"todayVolume"`    // 今日传输量(GB)
	WeekVolume     float64           `json:"weekVolume"`     // 本周传输量(GB)
	MonthVolume    float64           `json:"monthVolume"`    // 本月传输量(GB)
	TotalVolume    float64           `json:"totalVolume"`    // 总传输量(GB)
	AvgSpeed       float64           `json:"avgSpeed"`       // 平均传输速度(MB/s)
	VolumeTrend    []VolumeTrendItem `json:"volumeTrend"`    // 传输量趋势
}

// VolumeTrendItem 传输量趋势项
type VolumeTrendItem struct {
	Date   string  `json:"date"`
	Volume float64 `json:"volume"` // GB
	Count  int64   `json:"count"`
}

// TrendItem 趋势数据项
type TrendItem struct {
	Date   string `json:"date"`
	Count  int64  `json:"count"`
}

// FormStatItem 表单统计项
type FormStatItem struct {
	FormCode string `json:"formCode"`
	FormName string `json:"formName"`
	Count    int64  `json:"count"`
}

// ReportGenerateReq 生成报表请求
type ReportGenerateReq struct {
	FormCode   string   `json:"formCode" validate:"required"`
	StartDate  string   `json:"startDate" validate:"required"`
	EndDate    string   `json:"endDate" validate:"required"`
	Title      string   `json:"title"`
	Fields     []string `json:"fields"` // 导出的字段
	Format     string   `json:"format" validate:"required,oneof=excel csv pdf"`
}

// ============ 系统设置相关 ============

// SystemSettings 系统设置
type SystemSettings struct {
	SiteName        string `json:"siteName"`
	SiteLogo        string `json:"siteLogo"`
	AllowRegister   bool   `json:"allowRegister"`
	DefaultRoleID   uint   `json:"defaultRoleID"`
	SessionTimeout  int    `json:"sessionTimeout"`
	PasswordMinLen  int    `json:"passwordMinLen"`
	PasswordRequire bool   `json:"passwordRequire"`
}

// ============ API Key 相关 ============

// APIKeyCreateReq 创建API Key请求
type APIKeyCreateReq struct {
	Name        string   `json:"name" validate:"required,max=64"`
	Permissions []string `json:"permissions"` // 允许的API权限
	ExpiresAt   *int64   `json:"expiresAt"`  // 过期时间
}

// APIKeyResp API Key响应
type APIKeyResp struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Key       string    `json:"key"` // 仅创建时返回
	Secret    string    `json:"secret"` // 仅创建时返回
	Permissions []string `json:"permissions"`
	ExpiresAt *int64    `json:"expiresAt"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	LastUsedAt *int64   `json:"lastUsedAt"`
}

// ============ 上传记录相关 ============

// UploadRecordCreateReq 创建上传记录请求
type UploadRecordCreateReq struct {
	DataType    string                 `json:"dataType" validate:"required,max=64"`
	ProjectID   *uint                  `json:"projectId"`   // 项目ID（可选，可通过 projectName 查找）
	ProjectName string                 `json:"projectName" validate:"max=128"`
	DestPath    string                 `json:"destPath" validate:"required,max=512"`
	FileSize    int64                  `json:"fileSize" validate:"required,min=0"`
	Uploader    string                 `json:"uploader" validate:"required,max=64"`
	Status      string                 `json:"status" validate:"required,oneof=pending processing completed failed"`
	Remark      string                 `json:"remark" validate:"max=512"`
	Data        map[string]interface{} `json:"data"` // 动态字段数据
	CreatedAt   string                 `json:"createdAt"` // 可选，格式：2006-01-02 或 2006-01-02T15:04:05
}

// UploadRecordUpdateReq 更新上传记录请求
type UploadRecordUpdateReq struct {
	ID          uint                   `json:"id" validate:"required"`
	Status      string                 `json:"status" validate:"required,oneof=pending processing completed failed"`
	Remark      string                 `json:"remark" validate:"max=512"`
	Data        map[string]interface{} `json:"data"` // 动态字段数据
}

// UploadRecordListReq 上传记录列表请求
type UploadRecordListReq struct {
	Page       int    `form:"page" validate:"min=1"`
	PageSize   int    `form:"pageSize" validate:"min=1,max=100"`
	DataType   string `form:"dataType"`
	ProjectName string `form:"projectName"`
	Status     string `form:"status"`
	Uploader   string `form:"uploader"`
	StartDate  string `form:"startDate"`
	EndDate    string `form:"endDate"`
	Keyword    string `form:"keyword"`
}

// UploadRecordResp 上传记录响应
type UploadRecordResp struct {
	ID          uint                  `json:"id"`
	SerialNo    string                `json:"serialNo"`
	DataType    string                `json:"dataType"`
	ProjectID   *uint                 `json:"projectId"`
	ProjectName string                `json:"projectName"`
	DestPath    string                `json:"destPath"`
	FileSize    int64                 `json:"fileSize"`
	FileSizeStr string                `json:"fileSizeStr"`
	Uploader    string                `json:"uploader"`
	Status      string                `json:"status"`
	StatusText  string                `json:"statusText"`
	Remark      string                `json:"remark"`
	Data        map[string]interface{} `json:"data"`
	CreatedAt   string                `json:"createdAt"`
	UpdatedAt   string                `json:"updatedAt"`
}

// UploadRecordStatisticsResp 上传记录统计响应
type UploadRecordStatisticsResp struct {
	TodayCount   int64            `json:"todayCount"`
	TodaySize    int64            `json:"todaySize"`
	TodaySizeStr string           `json:"todaySizeStr"`
	WeekCount    int64            `json:"weekCount"`
	WeekSize     int64            `json:"weekSize"`
	WeekSizeStr  string           `json:"weekSizeStr"`
	MonthCount   int64            `json:"monthCount"`
	MonthSize    int64            `json:"monthSize"`
	MonthSizeStr string           `json:"monthSizeStr"`
	TotalCount   int64            `json:"totalCount"`
	TotalSize    int64            `json:"totalSize"`
	TotalSizeStr string           `json:"totalSizeStr"`
	Trend        []DailyTrend     `json:"trend"`
	ByStatus     []StatusCount    `json:"byStatus"`
	ByDataType   []DataTypeCount  `json:"byDataType"`
	ByProject    []ProjectCount   `json:"byProject"`
}

// StatusCount 状态统计
type StatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

// DataTypeCount 数据标签统计
type DataTypeCount struct {
	DataType  string `json:"dataType"`
	Count     int64  `json:"count"`
	TotalSize int64  `json:"totalSize"`
}

// ProjectCount 项目统计
type ProjectCount struct {
	ProjectName string `json:"projectName"`
	Count       int64  `json:"count"`
	TotalSize   int64  `json:"totalSize"`
}

// DailyTrend 每日趋势
type DailyTrend struct {
	Date      string `json:"date"`
	Count     int64  `json:"count"`
	TotalSize int64  `json:"totalSize"`
}

// ============ 批量导入相关 ============

// ImportTemplateField 导入模板字段定义
type ImportTemplateField struct {
	Field     string `json:"field"`      // Excel列名（中文）
	Code      string `json:"code"`       // 字段代码
	Required  bool   `json:"required"`   // 是否必填
	Type      string `json:"type"`       // 字段类型：text/number/select/date
	Options   string `json:"options"`    // 下拉选项（逗号分隔）
	MaxLength int    `json:"maxLength"`  // 最大长度
	Example   string `json:"example"`    // 填写示例
}

// ImportTemplateResp 导入模板信息响应
type ImportTemplateResp struct {
	Fields    []ImportTemplateField `json:"fields"`    // 所有字段定义
	SheetName string                `json:"sheetName"` // 工作表名称
	Title     string                `json:"title"`     // 标题行文字
}

// ImportReq 导入请求（内部使用，文件通过multipart上传）
type ImportReq struct {
	Overwrite bool `json:"overwrite"` // 是否覆盖已存在记录
}

// ImportFailRow 导入失败行
type ImportFailRow struct {
	Row     int    `json:"row"`     // Excel行号（从2开始，第1行是表头）
	Data    string `json:"data"`   // 该行原始数据
	Reason  string `json:"reason"`  // 失败原因
}

// ImportResultResp 导入结果响应
type ImportResultResp struct {
	Total      int             `json:"total"`       // 总行数
	Success    int             `json:"success"`     // 成功行数
	Failed     int             `json:"failed"`      // 失败行数
	FailRows   []ImportFailRow `json:"failRows"`    // 失败行明细
}

// ============ 安全设置相关 ============

// SecuritySettingsUpdateReq 更新安全设置请求
type SecuritySettingsUpdateReq struct {
	// 登录验证码
	CaptchaEnabled          bool `json:"captchaEnabled"`
	CaptchaMinLen           int  `json:"captchaMinLen" validate:"min=1,max=10"`

	// 不活跃用户自动禁用
	InactiveAutoDisable    bool `json:"inactiveAutoDisable"`
	InactiveDaysThreshold  int  `json:"inactiveDaysThreshold" validate:"min=1,max=365"`

	// 用户级别登录限制
	UserLoginMaxAttempts int `json:"userLoginMaxAttempts" validate:"min=1,max=20"`
	UserLoginLockMinutes int `json:"userLoginLockMinutes" validate:"min=1,max:10080"`

	// IP级别登录限制
	IPLoginMaxAttempts int `json:"ipLoginMaxAttempts" validate:"min=1,max=100"`
	IPLoginLockMinutes int `json:"ipLoginLockMinutes" validate:"min=1,max:10080"`

	// IP白名单（逗号分隔）
	IPWhitelist string `json:"ipWhitelist" validate:"max=2048"`

	// IP黑名单（逗号分隔）
	IPBlacklist string `json:"ipBlacklist" validate:"max=2048"`

	// 密码安全
	PasswordExpiryDays       int  `json:"passwordExpiryDays" validate:"min=0,max=365"`
	PasswordMinLength        int  `json:"passwordMinLength" validate:"min=6,max=128"`
	PasswordRequireUppercase bool `json:"passwordRequireUppercase"`
	PasswordRequireLowercase bool `json:"passwordRequireLowercase"`
	PasswordRequireDigit     bool `json:"passwordRequireDigit"`
	PasswordRequireSpecial   bool `json:"passwordRequireSpecial"`

	// 会话配置
	SessionTimeoutHours int `json:"sessionTimeoutHours" validate:"min=1,max=168"`
}

// IPBlacklistUpdateReq IP黑名单更新请求
type IPBlacklistUpdateReq struct {
	IP       string `json:"ip" validate:"required,ip"`
	Reason   string `json:"reason" validate:"max=256"`
	Operator string `json:"operator" validate:"required"` // 操作人
}

// IPBlacklistResp IP黑名单记录响应
type IPBlacklistResp struct {
	ID        uint      `json:"id"`
	IP        string    `json:"ip"`
	Reason    string    `json:"reason"`
	AddedBy   string    `json:"addedBy"`
	AddedAt   time.Time `json:"addedAt"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	IsActive  bool      `json:"isActive"`
}

// SecurityOverviewResp 安全总览响应
type SecurityOverviewResp struct {
	TotalLockedUsers   int `json:"totalLockedUsers"`
	TotalLockedIPs     int `json:"totalLockedIPs"`
	WhitelistCount     int `json:"whitelistCount"`
	BlacklistCount     int `json:"blacklistCount"`
	PasswordExpiryWarn  int `json:"passwordExpiryWarn"`  // 即将过期用户数
	InactiveUsersWarn   int `json:"inactiveUsersWarn"`  // 不活跃用户数
}

