package dto

import "time"

// ============ 操作日志相关 ============

// OperationLogResp 操作日志响应
type OperationLogResp struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"userId"`
	Username    string    `json:"username"`
	MenuName    string    `json:"menuName"`
	Action      string    `json:"action"`
	ActionText  string    `json:"actionText"`
	ResourceType string   `json:"resourceType"`
	ResourceID  uint      `json:"resourceId"`
	ResourceName string   `json:"resourceName"`
	IPAddress   string    `json:"ipAddress"`
	Detail      string    `json:"detail"`
	CreatedAt   time.Time `json:"createdAt"`
}

// OperationLogListReq 操作日志列表请求
type OperationLogListReq struct {
	Page        int    `form:"page" validate:"min=1"`
	PageSize    int    `form:"pageSize" validate:"min=1,max=100"`
	Keyword     string `form:"keyword"`
	UserID      uint   `form:"userId"`
	MenuName    string `form:"menuName"`
	Action      string `form:"action"`
	ResourceType string `form:"resourceType"`
	StartDate   string `form:"startDate"`
	EndDate     string `form:"endDate"`
}

// ============ 登录日志相关 ============

// LoginLogResp 登录日志响应
type LoginLogResp struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"userId"`
	Username   string    `json:"username"`
	Status     string    `json:"status"`
	StatusText string    `json:"statusText"`
	IPAddress  string    `json:"ipAddress"`
	FailReason string    `json:"failReason"`
	MFAUsed    bool      `json:"mfaUsed"`
	LoginMethod string   `json:"loginMethod"`
	CreatedAt  time.Time `json:"createdAt"`
}

// LoginLogListReq 登录日志列表请求
type LoginLogListReq struct {
	Page      int    `form:"page" validate:"min=1"`
	PageSize  int    `form:"pageSize" validate:"min=1,max=100"`
	Keyword   string `form:"keyword"`
	UserID    uint   `form:"userId"`
	Status    string `form:"status"`
	StartDate string `form:"startDate"`
	EndDate   string `form:"endDate"`
}

// ExportResult 导出结果（含截断信息）
type ExportResult struct {
	Total     int64  `json:"total"`
	Exported  int    `json:"exported"`
	Truncated bool   `json:"truncated"`
	TruncMsg  string `json:"truncMsg,omitempty"`
}
