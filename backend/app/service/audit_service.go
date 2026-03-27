package service

import (
	"datauptwo/app/dto"
	"datauptwo/app/model"
	"datauptwo/app/repo"
	"encoding/json"
	"fmt"
)

type AuditService struct {
	opLogRepo  *repo.OperationLogRepo
	loginRepo  *repo.LoginLogRepo
}

func NewAuditService() *AuditService {
	return &AuditService{
		opLogRepo:  repo.NewOperationLogRepo(),
		loginRepo:  repo.NewLoginLogRepo(),
	}
}

// LogOperation 记录操作日志
func (s *AuditService) LogOperation(userID uint, username, menuName, action, resourceType string, resourceID uint, resourceName, ipAddress, userAgent string, detail interface{}) error {
	detailJSON, _ := json.Marshal(detail)
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
