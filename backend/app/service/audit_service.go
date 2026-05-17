package service

import (
	"datauptwo/app/dto"
	"datauptwo/app/model"
	"datauptwo/app/repo"
	"encoding/json"
	"fmt"
)

type AuditService struct {
	opLogRepo *repo.OperationLogRepo
	loginRepo *repo.LoginLogRepo
}

func NewAuditService() *AuditService {
	return &AuditService{
		opLogRepo: repo.NewOperationLogRepo(),
		loginRepo: repo.NewLoginLogRepo(),
	}
}

// FieldChange 字段变更结构
type FieldChange struct {
	Field    string      `json:"field"`    // 字段名
	Label    string      `json:"label"`    // 字段中文名
	OldValue interface{} `json:"oldValue"` // 变更前的值
	NewValue interface{} `json:"newValue"` // 变更后的值
}

// OperationDetail 操作明细
type OperationDetail struct {
	Before interface{}     `json:"before,omitempty"` // 操作前的完整数据
	After  interface{}     `json:"after,omitempty"`  // 操作后的完整数据
	Changes []FieldChange  `json:"changes,omitempty"` // 变更的字段列表
}

// LogOperation 记录操作日志
func (s *AuditService) LogOperation(userID uint, username, menuName, action, resourceType string, resourceID uint, resourceName, ipAddress, userAgent string, detail interface{}) error {
	var detailStr string
	if detail != nil {
		detailJSON, _ := json.Marshal(detail)
		detailStr = string(detailJSON)
	}
	log := &model.OperationLog{
		UserID:       userID,
		Username:     username,
		MenuName:     menuName,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		ResourceName: resourceName,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		Detail:       detailStr,
	}
	return s.opLogRepo.Create(log)
}

// LogOperationWithDetail 记录操作日志（带变更明细）
func (s *AuditService) LogOperationWithDetail(userID uint, username, menuName, action, resourceType string, resourceID uint, resourceName, ipAddress, userAgent string, operationDetail *OperationDetail) error {
	detailJSON, _ := json.Marshal(operationDetail)
	log := &model.OperationLog{
		UserID:       userID,
		Username:     username,
		MenuName:     menuName,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		ResourceName: resourceName,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		Detail:       string(detailJSON),
	}
	return s.opLogRepo.Create(log)
}

// GetChanges 计算两个对象之间的字段变更
func GetChanges[T any](oldObj, newObj T, fieldLabels map[string]string) []FieldChange {
	var changes []FieldChange
	oldJSON, _ := json.Marshal(oldObj)
	newJSON, _ := json.Marshal(newObj)

	// 如果完全相同，不记录变更
	if string(oldJSON) == string(newJSON) {
		return changes
	}

	oldMap := make(map[string]interface{})
	newMap := make(map[string]interface{})
	json.Unmarshal(oldJSON, &oldMap)
	json.Unmarshal(newJSON, &newMap)

	for key, newVal := range newMap {
		oldVal := oldMap[key]
		// 比较值是否变化（简单比较，复杂类型可能不准确）
		if fmt.Sprintf("%v", oldVal) != fmt.Sprintf("%v", newVal) {
			label := fieldLabels[key]
			if label == "" {
				label = key
			}
			changes = append(changes, FieldChange{
				Field:    key,
				Label:    label,
				OldValue: oldVal,
				NewValue: newVal,
			})
		}
	}
	return changes
}

// ListOperationLogs 获取操作日志列表
func (s *AuditService) ListOperationLogs(req dto.OperationLogListReq) (*dto.PageResult, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	logs, total, err := s.opLogRepo.List(req.Page, req.PageSize, req.UserID, req.Keyword, req.MenuName, req.Action, req.ResourceType, req.Keyword, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	// 转换响应
	items := make([]dto.OperationLogResp, len(logs))
	actionMap := map[string]string{
		"view":    "查看",
		"create":  "创建",
		"update":  "更新",
		"delete":  "删除",
		"export":  "导出",
		"login":   "登录",
		"logout":  "登出",
		"import":  "导入",
	}
	for i, log := range logs {
		items[i] = dto.OperationLogResp{
			ID:           log.ID,
			UserID:       log.UserID,
			Username:     log.Username,
			MenuName:     log.MenuName,
			Action:       log.Action,
			ActionText:   actionMap[log.Action],
			ResourceType: log.ResourceType,
			ResourceID:   log.ResourceID,
			ResourceName: log.ResourceName,
			IPAddress:   log.IPAddress,
			Detail:      log.Detail,
			CreatedAt:   log.CreatedAt,
		}
	}

	return &dto.PageResult{
		Total: total,
		Items: items,
	}, nil
}

// ListOperationLogsForExport 获取操作日志列表用于导出
func (s *AuditService) ListOperationLogsForExport(req dto.OperationLogListReq) ([]dto.OperationLogResp, *dto.ExportResult, error) {
	logs, total, truncated, err := s.opLogRepo.ListForExport(req.UserID, req.Keyword, req.MenuName, req.Action, req.ResourceType, req.Keyword, req.StartDate, req.EndDate)
	if err != nil {
		return nil, nil, err
	}

	actionMap := map[string]string{
		"view":   "查看",
		"create": "创建",
		"update": "更新",
		"delete": "删除",
		"export": "导出",
		"login":  "登录",
		"logout": "登出",
		"import": "导入",
	}

	items := make([]dto.OperationLogResp, len(logs))
	for i, log := range logs {
		items[i] = dto.OperationLogResp{
			ID:           log.ID,
			UserID:       log.UserID,
			Username:     log.Username,
			MenuName:     log.MenuName,
			Action:       log.Action,
			ActionText:   actionMap[log.Action],
			ResourceType: log.ResourceType,
			ResourceID:   log.ResourceID,
			ResourceName: log.ResourceName,
			IPAddress:    log.IPAddress,
			Detail:       log.Detail,
			CreatedAt:    log.CreatedAt,
		}
	}

	exportResult := &dto.ExportResult{
		Total:     total,
		Exported:  len(items),
		Truncated: truncated,
	}
	if truncated {
		exportResult.TruncMsg = fmt.Sprintf("导出已截断，仅展示前 %d 条（共 %d 条）", 10000, total)
	}

	return items, exportResult, nil
}

// LogLogin 记录登录日志
func (s *AuditService) LogLogin(userID uint, username, status, ipAddress, userAgent, failReason string, mfaUsed bool, loginMethod string) error {
	log := &model.LoginLog{
		UserID:      userID,
		Username:    username,
		Status:      status,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		FailReason:  failReason,
		MFAUsed:     mfaUsed,
		LoginMethod: loginMethod,
	}
	return s.loginRepo.Create(log)
}

// ListLoginLogs 获取登录日志列表
func (s *AuditService) ListLoginLogs(req dto.LoginLogListReq) (*dto.PageResult, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	logs, total, err := s.loginRepo.List(req.Page, req.PageSize, req.UserID, "", req.Status, req.Keyword, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	// 转换响应
	items := make([]dto.LoginLogResp, len(logs))
	for i, log := range logs {
		statusText := "成功"
		if log.Status == "failed" {
			statusText = "失败"
		}
		items[i] = dto.LoginLogResp{
			ID:          log.ID,
			UserID:      log.UserID,
			Username:    log.Username,
			Status:      log.Status,
			StatusText:  statusText,
			IPAddress:   log.IPAddress,
			FailReason:  log.FailReason,
			MFAUsed:     log.MFAUsed,
			LoginMethod: log.LoginMethod,
			CreatedAt:  log.CreatedAt,
		}
	}

	return &dto.PageResult{
		Total: total,
		Items: items,
	}, nil
}

// ListLoginLogsForExport 获取登录日志列表用于导出
// ListLoginLogsForExport 获取登录日志列表用于导出
func (s *AuditService) ListLoginLogsForExport(req dto.LoginLogListReq) ([]dto.LoginLogResp, *dto.ExportResult, error) {
	logs, total, truncated, err := s.loginRepo.ListForExport(req.UserID, "", req.Status, req.Keyword, req.StartDate, req.EndDate)
	if err != nil {
		return nil, nil, err
	}

	items := make([]dto.LoginLogResp, len(logs))
	for i, log := range logs {
		statusText := "成功"
		if log.Status == "failed" {
			statusText = "失败"
		}
		items[i] = dto.LoginLogResp{
			ID:          log.ID,
			UserID:      log.UserID,
			Username:    log.Username,
			Status:      log.Status,
			StatusText:  statusText,
			IPAddress:   log.IPAddress,
			FailReason:  log.FailReason,
			MFAUsed:     log.MFAUsed,
			LoginMethod: log.LoginMethod,
			CreatedAt:   log.CreatedAt,
		}
	}

	exportResult := &dto.ExportResult{
		Total:     total,
		Exported:  len(items),
		Truncated: truncated,
	}
	if truncated {
		exportResult.TruncMsg = fmt.Sprintf("导出已截断，仅展示前 %d 条（共 %d 条）", 10000, total)
	}

	return items, exportResult, nil
}
