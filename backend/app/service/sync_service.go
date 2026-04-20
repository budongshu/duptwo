package service

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"datauptwo/app/dto"
	"datauptwo/app/model"
	"datauptwo/app/repo"
	"datauptwo/global"
)

// SyncService 同步服务
type SyncService struct {
	stationRepo  *repo.SyncStationRepo
	historyRepo  *repo.SyncHistoryRepo
	detailRepo   *repo.SyncDetailRepo
	recordRepo   *repo.UploadRecordRepo
}

func NewSyncService() *SyncService {
	return &SyncService{
		stationRepo:  repo.NewSyncStationRepo(),
		historyRepo:  repo.NewSyncHistoryRepo(),
		detailRepo:   repo.NewSyncDetailRepo(),
		recordRepo:   repo.NewUploadRecordRepo(),
	}
}

// ============ 站点管理 ============

// CreateStation 创建站点
func (s *SyncService) CreateStation(req dto.SyncStationCreateReq) (*dto.SyncStationResp, error) {
	// 检查代码是否已存在
	existing, _ := s.stationRepo.GetByCode(req.Code)
	if existing != nil {
		return nil, errors.New("站点代码已存在")
	}

	// 生成 API Key
	apiKey, apiKeyHash := s.generateAPIKey()

	station := &model.SyncStation{
		Name:        req.Name,
		Code:        req.Code,
		URL:         req.URL,
		Status:      "active",
		Description: req.Description,
		IsCenter:    req.IsCenter,
		APIKey:      apiKeyHash,
	}

	if err := s.stationRepo.Create(station); err != nil {
		return nil, err
	}

	resp := s.toStationResp(station)
	resp.APIKey = apiKey // 返回明文 API Key（仅创建时返回一次）
	return resp, nil
}

// UpdateStation 更新站点
func (s *SyncService) UpdateStation(req dto.SyncStationUpdateReq) (*dto.SyncStationResp, error) {
	station, err := s.stationRepo.GetByID(req.ID)
	if err != nil {
		return nil, errors.New("站点不存在")
	}

	// 检查新代码是否与其他站点冲突
	if req.Code != station.Code {
		existing, _ := s.stationRepo.GetByCode(req.Code)
		if existing != nil && existing.ID != req.ID {
			return nil, errors.New("站点代码已存在")
		}
	}

	station.Name = req.Name
	station.Code = req.Code
	station.URL = req.URL
	station.Description = req.Description
	station.IsCenter = req.IsCenter

	if req.Status != "" {
		station.Status = req.Status
	}

	if err := s.stationRepo.Update(station); err != nil {
		return nil, err
	}

	return s.toStationResp(station), nil
}

// DeleteStation 删除站点
func (s *SyncService) DeleteStation(id uint) error {
	_, err := s.stationRepo.GetByID(id)
	if err != nil {
		return errors.New("站点不存在")
	}
	return s.stationRepo.Delete(id)
}

// GetStation 获取站点详情
func (s *SyncService) GetStation(id uint) (*dto.SyncStationResp, error) {
	station, err := s.stationRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("站点不存在")
	}
	return s.toStationResp(station), nil
}

// ListStations 获取站点列表
func (s *SyncService) ListStations(req dto.SyncStationListReq) (*dto.PageResult, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}

	stations, total, err := s.stationRepo.List(req)
	if err != nil {
		return nil, err
	}

	items := make([]dto.SyncStationResp, len(stations))
	for i, station := range stations {
		items[i] = *s.toStationResp(&station)
	}

	return &dto.PageResult{Total: total, Items: items}, nil
}

// GetAllStations 获取所有站点
func (s *SyncService) GetAllStations() ([]dto.SyncStationResp, error) {
	stations, err := s.stationRepo.GetAll()
	if err != nil {
		return nil, err
	}

	items := make([]dto.SyncStationResp, len(stations))
	for i, station := range stations {
		items[i] = *s.toStationResp(&station)
	}

	return items, nil
}

// ============ 站点注册 ============

// RegisterStation 注册站点
func (s *SyncService) RegisterStation(req dto.SyncRegisterReq) (*dto.SyncRegisterResp, error) {
	// 查找站点（根据代码）
	station, err := s.stationRepo.GetByCode(req.StationCode)
	if err != nil {
		return nil, errors.New("站点不存在，请先在中心站点创建站点记录")
	}

	// 验证密码（使用SHA256哈希比较）
	passwordHash := s.hashPassword(req.Password)
	// TODO: 实际应该使用 bcrypt 比较，这里简化处理
	if station.APIKey != "" && station.APIKey != passwordHash {
		return nil, errors.New("站点密码错误")
	}

	// 生成新的 API Key
	apiKey, apiKeyHash := s.generateAPIKey()

	// 更新站点的 API Key
	station.APIKey = apiKeyHash
	station.URL = req.URL
	if station.Status == "" {
		station.Status = "active"
	}

	if err := s.stationRepo.Update(station); err != nil {
		return nil, err
	}

	return &dto.SyncRegisterResp{
		StationID:  station.ID,
		APIKey:     apiKey,
		APIKeyHash: apiKeyHash,
		Message:    "注册成功，请妥善保管API Key",
	}, nil
}

// ValidateAPIKey 验证API Key
func (s *SyncService) ValidateAPIKey(apiKey string) (*model.SyncStation, error) {
	// 查找所有活跃站点，逐一验证
	stations, err := s.stationRepo.GetAll()
	if err != nil {
		return nil, err
	}

	for _, station := range stations {
		if s.compareAPIKey(apiKey, station.APIKey) {
			return &station, nil
		}
	}

	return nil, errors.New("无效的API Key")
}

// ============ 记录上传同步 ============

// UploadRecords 处理上传记录同步
func (s *SyncService) UploadRecords(req dto.SyncUploadReq, stationID uint) (*dto.SyncUploadResp, error) {
	// 创建同步历史记录
	history := &model.SyncHistory{
		StationID:    stationID,
		Direction:    "upload",
		Status:       "pending",
		TotalRecords: len(req.Records),
	}

	if err := s.historyRepo.Create(history); err != nil {
		return nil, err
	}

	// 设置开始状态
	s.historyRepo.SetStarted(history.ID)

	resp := &dto.SyncUploadResp{
		TotalRecords: len(req.Records),
		Details:      make([]dto.SyncDetailItem, 0, len(req.Records)),
		HistoryID:    history.ID,
	}

	// 事务处理每条记录
	successCount := 0
	failCount := 0
	conflictCount := 0

	for _, record := range req.Records {
		detail, err := s.processRecord(history.ID, &record)
		if err != nil {
			resp.Details = append(resp.Details, *detail)
			failCount++
			continue
		}

		resp.Details = append(resp.Details, *detail)

		switch detail.Action {
		case "create", "update":
			successCount++
		case "conflict":
			conflictCount++
		case "skip":
			// skip 不计入成功也不计入冲突
		}
	}

	resp.SuccessCount = successCount
	resp.FailCount = failCount
	resp.ConflictCount = conflictCount

	// 更新历史记录状态
	status := "completed"
	if failCount == len(req.Records) {
		status = "failed"
	}
	s.historyRepo.UpdateStatus(history.ID, status, "")

	return resp, nil
}

// processRecord 处理单条记录（包含冲突处理逻辑）
func (s *SyncService) processRecord(historyID uint, record *dto.SyncRecordItem) (*dto.SyncDetailItem, error) {
	detail := &dto.SyncDetailItem{
		SerialNo:    record.SerialNo,
		ProjectName: record.ProjectName,
	}

	// 查找 serial_no 是否存在
	existing, err := s.recordRepo.GetBySerialNo(record.SerialNo)
	if err != nil {
		// 记录不存在，创建新记录
		return s.createRecord(historyID, detail, record)
	}

	// 记录存在，检查 project_name
	if existing.ProjectName == record.ProjectName {
		// project_name 相同，更新记录
		return s.updateRecord(historyID, detail, existing, record)
	}

	// project_name 不同，生成新 serial_no 创建
	return s.createConflictedRecord(historyID, detail, record, existing.SerialNo)
}

// createRecord 创建新记录
func (s *SyncService) createRecord(historyID uint, detail *dto.SyncDetailItem, record *dto.SyncRecordItem) (*dto.SyncDetailItem, error) {
	// 序列化 Data 字段
	dataJSON := ""
	if record.Data != nil {
		dataBytes, _ := json.Marshal(record.Data)
		dataJSON = string(dataBytes)
	}

	newRecord := &model.UploadRecord{
		SerialNo:    record.SerialNo,
		ProjectName: record.ProjectName,
		DiskLabel:   record.DiskLabel,
		DestPath:    record.DestPath,
		FileSize:    record.FileSize,
		Uploader:    record.Uploader,
		Status:      record.Status,
		Remark:      record.Remark,
		Data:        dataJSON,
	}

	if err := s.recordRepo.Create(newRecord); err != nil {
		detail.Action = "create"
		detail.Result = "failed"
		detail.ErrorMsg = err.Error()
		s.saveDetail(historyID, detail, record.SerialNo, "")
		return detail, err
	}

	detail.Action = "create"
	detail.Result = "success"
	s.saveDetail(historyID, detail, record.SerialNo, "")

	return detail, nil
}

// updateRecord 更新记录
func (s *SyncService) updateRecord(historyID uint, detail *dto.SyncDetailItem, existing *model.UploadRecord, record *dto.SyncRecordItem) (*dto.SyncDetailItem, error) {
	// 更新字段
	existing.DiskLabel = record.DiskLabel
	existing.DestPath = record.DestPath
	existing.FileSize = record.FileSize
	existing.Uploader = record.Uploader
	existing.Status = record.Status
	existing.Remark = record.Remark

	if record.Data != nil {
		dataBytes, _ := json.Marshal(record.Data)
		existing.Data = string(dataBytes)
	}

	if err := s.recordRepo.Update(existing); err != nil {
		detail.Action = "update"
		detail.Result = "failed"
		detail.ErrorMsg = err.Error()
		s.saveDetail(historyID, detail, record.SerialNo, "")
		return detail, err
	}

	detail.Action = "update"
	detail.Result = "success"
	s.saveDetail(historyID, detail, record.SerialNo, "")

	return detail, nil
}

// createConflictedRecord 创建冲突记录
func (s *SyncService) createConflictedRecord(historyID uint, detail *dto.SyncDetailItem, record *dto.SyncRecordItem, oldSerialNo string) (*dto.SyncDetailItem, error) {
	// 生成新 serial_no = 原serial_no + "_" + md5(project_name)[:6]
	hash := md5.Sum([]byte(record.ProjectName))
	newSerialNo := fmt.Sprintf("%s_%s", record.SerialNo, hex.EncodeToString(hash[:])[:6])

	// 检查新 serial_no 是否已存在，如果存在则追加随机后缀
	_, err := s.recordRepo.GetBySerialNo(newSerialNo)
	if err == nil {
		// 已存在，再追加随机后缀
		randomBytes := make([]byte, 4)
		rand.Read(randomBytes)
		newSerialNo = fmt.Sprintf("%s_%s", newSerialNo, hex.EncodeToString(randomBytes))
	}

	// 序列化 Data 字段
	dataJSON := ""
	if record.Data != nil {
		dataBytes, _ := json.Marshal(record.Data)
		dataJSON = string(dataBytes)
	}

	newRecord := &model.UploadRecord{
		SerialNo:    newSerialNo,
		ProjectName: record.ProjectName,
		DiskLabel:   record.DiskLabel,
		DestPath:    record.DestPath,
		FileSize:    record.FileSize,
		Uploader:    record.Uploader,
		Status:      record.Status,
		Remark:      record.Remark,
		Data:        dataJSON,
	}

	if err := s.recordRepo.Create(newRecord); err != nil {
		detail.Action = "conflict"
		detail.Result = "failed"
		detail.ErrorMsg = err.Error()
		detail.OldSerialNo = oldSerialNo
		s.saveDetail(historyID, detail, record.SerialNo, newSerialNo)
		return detail, err
	}

	detail.Action = "conflict"
	detail.Result = "success"
	detail.NewSerialNo = newSerialNo
	detail.OldSerialNo = oldSerialNo
	s.saveDetail(historyID, detail, record.SerialNo, newSerialNo)

	return detail, nil
}

// saveDetail 保存同步详情
func (s *SyncService) saveDetail(historyID uint, detail *dto.SyncDetailItem, oldSerialNo, newSerialNo string) {
	syncDetail := &model.SyncDetail{
		HistoryID:   historyID,
		SerialNo:    detail.SerialNo,
		ProjectName: detail.ProjectName,
		Action:      detail.Action,
		Result:      detail.Result,
		ErrorMsg:    detail.ErrorMsg,
		OldSerialNo: oldSerialNo,
		NewSerialNo: newSerialNo,
	}
	s.detailRepo.Create(syncDetail)
}

// ============ 同步历史 ============

// GetHistory 获取同步历史
func (s *SyncService) GetHistory(req dto.SyncHistoryListReq) (*dto.PageResult, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}

	histories, total, err := s.historyRepo.List(req)
	if err != nil {
		return nil, err
	}

	items := make([]dto.SyncHistoryResp, len(histories))
	for i, history := range histories {
		items[i] = *s.toHistoryResp(&history)
	}

	return &dto.PageResult{Total: total, Items: items}, nil
}

// GetHistoryDetails 获取同步历史详情
func (s *SyncService) GetHistoryDetails(historyID uint) (*dto.SyncHistoryDetailResp, error) {
	history, err := s.historyRepo.GetByID(historyID)
	if err != nil {
		return nil, errors.New("同步历史不存在")
	}

	details, err := s.detailRepo.ListByHistoryID(historyID)
	if err != nil {
		return nil, err
	}

	resp := &dto.SyncHistoryDetailResp{
		SyncHistoryResp: *s.toHistoryResp(history),
		Details:         make([]dto.SyncDetailResp, len(details)),
	}

	for i, detail := range details {
		resp.Details[i] = *s.toDetailResp(&detail)
	}

	return resp, nil
}

// GetSyncStatus 获取同步状态
func (s *SyncService) GetSyncStatus() *dto.SyncStatusResp {
	resp := &dto.SyncStatusResp{
		Enabled: global.CONF.Sync.Enabled,
		Mode:    global.CONF.Sync.Mode,
	}

	// 获取中心站点信息
	if global.CONF.Sync.Mode == "agent" && global.CONF.Sync.StationID != "" {
		// Agent模式，尝试获取本地站点信息
		resp.IsCenter = false
		resp.StationID = global.CONF.Sync.StationID
		resp.StationName = global.CONF.Sync.StationName
		resp.CenterURL = global.CONF.Sync.CenterURL

		// 从调度器获取注册状态
		if scheduler, ok := global.Scheduler.(*SyncScheduler); ok {
			status := scheduler.GetQueueStatus()
			resp.Registered, _ = status["registered"].(bool)
			if t, ok := status["lastErrorAt"].(time.Time); ok {
				resp.LastErrorAt = &t
			}
		}
	} else if global.CONF.Sync.Mode == "center" {
		// Center模式
		resp.IsCenter = true
		centerStation, err := s.stationRepo.GetCenterStation()
		if err == nil && centerStation != nil {
			resp.StationID = fmt.Sprintf("%d", centerStation.ID)
			resp.StationName = centerStation.Name
			resp.LastSyncAt = centerStation.LastSyncAt
		}
	}

	return resp
}

// ============ 辅助方法 ============

// toStationResp 转换为站点响应
func (s *SyncService) toStationResp(station *model.SyncStation) *dto.SyncStationResp {
	resp := &dto.SyncStationResp{
		ID:          station.ID,
		Name:        station.Name,
		Code:        station.Code,
		URL:         station.URL,
		Status:      station.Status,
		StatusText:  "正常",
		Description: station.Description,
		IsCenter:    station.IsCenter,
		LastSyncAt:  station.LastSyncAt,
		SyncCount:   station.SyncCount,
		Remark:      station.Remark,
		CreatedAt:   station.CreatedAt,
		UpdatedAt:   station.UpdatedAt,
	}

	if station.Status == "inactive" {
		resp.StatusText = "已禁用"
	}

	return resp
}

// toHistoryResp 转换为历史响应
func (s *SyncService) toHistoryResp(history *model.SyncHistory) *dto.SyncHistoryResp {
	resp := &dto.SyncHistoryResp{
		ID:            history.ID,
		StationID:     history.StationID,
		Direction:     history.Direction,
		DirectionText: s.getDirectionText(history.Direction),
		Status:        history.Status,
		StatusText:    s.getStatusText(history.Status),
		TotalRecords:  history.TotalRecords,
		SuccessCount:  history.SuccessCount,
		FailCount:     history.FailCount,
		ConflictCount: history.ConflictCount,
		StartedAt:     history.StartedAt,
		CompletedAt:   history.CompletedAt,
		ErrorMsg:      history.ErrorMsg,
		Remark:        history.Remark,
		CreatedAt:     history.CreatedAt,
	}

	if history.Station != nil {
		resp.StationName = history.Station.Name
	}

	return resp
}

// toDetailResp 转换为详情响应
func (s *SyncService) toDetailResp(detail *model.SyncDetail) *dto.SyncDetailResp {
	return &dto.SyncDetailResp{
		ID:          detail.ID,
		SerialNo:    detail.SerialNo,
		ProjectName: detail.ProjectName,
		Action:      detail.Action,
		ActionText:  s.getActionText(detail.Action),
		Result:      detail.Result,
		ErrorMsg:    detail.ErrorMsg,
		OldSerialNo: detail.OldSerialNo,
		NewSerialNo: detail.NewSerialNo,
		CreatedAt:   detail.CreatedAt,
	}
}

// getDirectionText 获取方向文本
func (s *SyncService) getDirectionText(direction string) string {
	switch direction {
	case "upload":
		return "上传"
	case "download":
		return "下载"
	default:
		return direction
	}
}

// getStatusText 获取状态文本
func (s *SyncService) getStatusText(status string) string {
	switch status {
	case "pending":
		return "等待中"
	case "processing":
		return "处理中"
	case "completed":
		return "已完成"
	case "failed":
		return "失败"
	default:
		return status
	}
}

// getActionText 获取操作文本
func (s *SyncService) getActionText(action string) string {
	switch action {
	case "create":
		return "创建"
	case "update":
		return "更新"
	case "skip":
		return "跳过"
	case "conflict":
		return "冲突处理"
	default:
		return action
	}
}

// hashPassword 密码哈希（SHA256）
func (s *SyncService) hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// generateAPIKey 生成32字节的十六进制API Key
func (s *SyncService) generateAPIKey() (string, string) {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	apiKey := hex.EncodeToString(bytes)

	// 存储哈希版本
	hash := sha256.Sum256([]byte(apiKey))
	apiKeyHash := hex.EncodeToString(hash[:])

	return apiKey, apiKeyHash
}

// compareAPIKey 比较API Key
func (s *SyncService) compareAPIKey(input, stored string) bool {
	hash := sha256.Sum256([]byte(input))
	inputHash := hex.EncodeToString(hash[:])
	return inputHash == stored
}
