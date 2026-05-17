package service

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"datauptwo/app/dto"
	"datauptwo/app/model"
	"datauptwo/app/repo"
	"datauptwo/global"
	"datauptwo/middleware"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
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

// ============ Center 主动探测 ============

// ProbeAllAgents 并发探测所有 Agent 站点的连通性（使用 errgroup 限制最多 10 个并发）
func (s *SyncService) ProbeAllAgents() {
	stations, err := s.stationRepo.GetAllForProbe()
	if err != nil {
		global.AppLogger.Error("[CenterProbe] 获取站点列表失败: %v", err)
		return
	}

	now := time.Now()
	global.AppLogger.Debug("[CenterProbe] 开始探测，当前时间: %s (unix=%d)，站点数: %d",
		now.Format("2006-01-02 15:04:05"), now.Unix(), len(stations))

	results := make(map[uint]bool)
	var mu sync.Mutex
	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(10) // 最多 10 个并发探测

	for _, station := range stations {
		if station.IsCenter || station.Status != "active" || station.URL == "" {
			continue // 跳过 Center / 非活跃 / 无 URL 的站点
		}

		station := station // 捕获循环变量
		g.Go(func() error {
			probeURL := strings.TrimSuffix(station.URL, "/") + "/health"
			err := s.pingURL(probeURL, 10*time.Second)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				global.SetProbeConnected(station.ID, false)
				results[station.ID] = false
				global.AppLogger.Warn("[CenterProbe] %s [%s] 不可达: %v", station.Name, station.Code, err)
			} else {
				global.SetProbeConnected(station.ID, true)
				results[station.ID] = true
				s.stationRepo.UpdateLastConnectedAt(station.ID, &now)
				global.AppLogger.Info("[CenterProbe] %s [%s] 在线", station.Name, station.Code)
			}
			return nil
		})
	}

	_ = ctx // 忽略 ctx 取消（所有探测均执行）
	g.Wait() // 等待所有探测完成

	global.SetProbeConnectedBatch(results)
}

// pingURL 发送 HTTP GET 请求探测 URL 可达性
func (s *SyncService) pingURL(targetURL string, timeout time.Duration) error {
	client := &http.Client{Timeout: timeout}
	resp, err := client.Get(targetURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return nil
}

// ============ 站点管理 ============

// CreateStation 创建站点
func (s *SyncService) CreateStation(req dto.SyncStationCreateReq) (*dto.SyncStationResp, error) {
	// 检查代码是否已存在
	existing, _ := s.stationRepo.GetByCode(req.Code)
	if existing != nil {
		return nil, fmt.Errorf("站点代码已存在")
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
		return nil, fmt.Errorf("站点不存在")
	}

	// 检查新代码是否与其他站点冲突
	if req.Code != station.Code {
		existing, _ := s.stationRepo.GetByCode(req.Code)
		if existing != nil && existing.ID != req.ID {
			return nil, fmt.Errorf("站点代码已存在")
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
		return fmt.Errorf("站点不存在")
	}
	return s.stationRepo.Delete(id)
}

// GetStation 获取站点详情
func (s *SyncService) GetStation(id uint) (*dto.SyncStationResp, error) {
	station, err := s.stationRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("站点不存在")
	}
	return s.toStationResp(station), nil
}

// ResetApiKey 重置站点API Key（不中断服务，生成新Key后旧Key立即失效）
func (s *SyncService) ResetApiKey(id uint) (*dto.SyncResetKeyResp, error) {
	station, err := s.stationRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("站点不存在")
	}

	// 生成新 API Key
	apiKey, apiKeyHash := s.generateAPIKey()

	// 更新站点的 API Key（包括明文 Key）
	station.APIKey = apiKeyHash
	station.PlainAPIKey = apiKey
	if err := s.stationRepo.Update(station); err != nil {
		return nil, err
	}

	// 刷新缓存，使旧 Key 立即失效，Agent 会自动重新注册
	middleware.RefreshCache()

	global.AppLogger.Info("站点 %s 的 API Key 已重置", station.Name)

	return &dto.SyncResetKeyResp{
		ID:     station.ID,
		APIKey: apiKey,
	}, nil
}

// UpdateHeartbeat 更新站点心跳时间
func (s *SyncService) UpdateHeartbeat(stationID uint) error {
	station, err := s.stationRepo.GetByID(stationID)
	if err != nil {
		return fmt.Errorf("站点不存在")
	}
	return s.stationRepo.UpdateHeartbeatByCode(station.Code)
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

// RegisterStation 注册站点（幂等：相同 code 只注册一次，后续复用已有 key）
// 完全幂等：先查是否已存在，已存在直接返回已有 key，避免 DB 并发写入冲突
func (s *SyncService) RegisterStation(req dto.SyncRegisterReq) (*dto.SyncRegisterResp, error) {
	// 1. 先查是否已存在（幂等核心：已有 key 直接返回，不生成新 key）
	if existing, err := s.stationRepo.GetByCode(req.StationCode); err == nil && existing != nil && existing.ID > 0 {
		// 已有记录，重新激活
		if existing.Status != "active" {
			existing.Status = "active"
			s.stationRepo.Update(existing)
			middleware.RefreshCache()
		}
		// 从数据库读取 plain key 返回（幂等注册必须返回明文，Agent 需要用明文计算 hash）
		return &dto.SyncRegisterResp{
			StationID:  existing.ID,
			APIKey:     existing.PlainAPIKey, // 明文 Key（Agent 用它计算 hash）
			APIKeyHash: existing.APIKey,     // hash（供参考）
			Message:    "注册成功（复用已有 Key）",
		}, nil
	}

	// 2. 不存在，生成新 key 并插入
	apiKey, apiKeyHash := s.generateAPIKey()
	station := &model.SyncStation{
		Name:         req.StationName,
		Code:        req.StationCode,
		URL:         req.URL,
		Status:      "active",
		Description: "自动注册",
		APIKey:      apiKeyHash,
		PlainAPIKey: apiKey, // 存储明文 Key
	}

	if err := s.stationRepo.UpsertByCode(station); err != nil {
		return nil, fmt.Errorf("注册站点失败: %w", err)
	}

	// 刷新缓存
	middleware.RefreshCache()

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

	return nil, fmt.Errorf("无效的API Key")
}

// ============ 记录上传同步 ============

// UploadRecords 处理上传记录同步（事务保证原子性：全部成功或全部回滚）
func (s *SyncService) UploadRecords(req dto.SyncUploadReq, stationID uint) (*dto.SyncUploadResp, error) {
	resp := &dto.SyncUploadResp{
		TotalRecords: len(req.Records),
		Details:      make([]dto.SyncDetailItem, 0, len(req.Records)),
	}

	// 开启事务：整批记录要么全部写入，要么全部回滚
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		// 创建同步历史记录（事务内）
		history := &model.SyncHistory{
			StationID:    stationID,
			Direction:    "upload",
			Status:       "pending",
			TotalRecords: len(req.Records),
		}
		if err := tx.Create(history).Error; err != nil {
			return fmt.Errorf("创建历史记录失败: %w", err)
		}

		resp.HistoryID = history.ID

		// 设置开始状态
		now := time.Now()
		if err := tx.Model(history).Updates(map[string]interface{}{
			"started_at": now,
			"status":     "processing",
		}).Error; err != nil {
			return fmt.Errorf("设置开始状态失败: %w", err)
		}

		// 事务内处理每条记录（使用同一个 tx）
		successCount, failCount, conflictCount := 0, 0, 0

		for i := range req.Records {
			detail, err := s.processRecordTx(tx, history.ID, &req.Records[i])
			resp.Details = append(resp.Details, *detail)

			if err != nil {
				failCount++
				continue
			}

			switch detail.Action {
			case "create", "update":
				successCount++
			case "conflict":
				conflictCount++
			}
		}

		resp.SuccessCount = successCount
		resp.FailCount = failCount
		resp.ConflictCount = conflictCount

		// 更新历史计数
		if err := tx.Model(history).Updates(map[string]interface{}{
			"status":          "completed",
			"completed_at":    now,
			"success_count":   gorm.Expr("success_count + ?", successCount),
			"fail_count":      gorm.Expr("fail_count + ?", failCount),
			"conflict_count":  gorm.Expr("conflict_count + ?", conflictCount),
		}).Error; err != nil {
			return fmt.Errorf("更新历史计数失败: %w", err)
		}

		global.AppLogger.Info("[UploadRecords] 完成 historyID=%d, stationID=%d, success=%d, fail=%d, conflict=%d",
			history.ID, stationID, successCount, failCount, conflictCount)

		// 如果有失败记录，整批返回错误，Agent 会重试失败的记录
		if failCount > 0 {
			return fmt.Errorf("部分记录处理失败: success=%d, fail=%d", successCount, failCount)
		}
		return nil
	})

	if err != nil {
		global.AppLogger.Warn("[UploadRecords] stationID=%d, %v", stationID, err)
	}
	return resp, nil
}

// processRecordTx 事务内处理单条记录
func (s *SyncService) processRecordTx(tx *gorm.DB, historyID uint, record *dto.SyncRecordItem) (*dto.SyncDetailItem, error) {
	detail := &dto.SyncDetailItem{
		SerialNo:    record.SerialNo,
		ProjectName: record.ProjectName,
	}

	// 查找 serial_no 是否存在
	var existing model.UploadRecord
	findErr := tx.Where("serial_no = ? AND is_deleted = ?", record.SerialNo, false).First(&existing).Error

	if findErr != nil {
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			return s.createRecordTx(tx, historyID, detail, record)
		}
		detail.Action = "create"
		detail.Result = "failed"
		detail.ErrorMsg = fmt.Sprintf("查询记录失败: %v", findErr)
		s.saveDetailTx(tx, historyID, detail, record.SerialNo, "")
		return detail, findErr
	}

	if existing.ProjectName == record.ProjectName {
		return s.updateRecordTx(tx, historyID, detail, &existing, record)
	}
	return s.createConflictedRecordTx(tx, historyID, detail, record, existing.SerialNo)
}

// createRecordTx 事务内创建记录
func (s *SyncService) createRecordTx(tx *gorm.DB, historyID uint, detail *dto.SyncDetailItem, record *dto.SyncRecordItem) (*dto.SyncDetailItem, error) {
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
	if err := tx.Create(newRecord).Error; err != nil {
		detail.Action = "create"
		detail.Result = "failed"
		detail.ErrorMsg = err.Error()
		s.saveDetailTx(tx, historyID, detail, record.SerialNo, "")
		return detail, err
	}
	detail.Action = "create"
	detail.Result = "success"
	s.saveDetailTx(tx, historyID, detail, record.SerialNo, "")
	return detail, nil
}

// updateRecordTx 事务内更新记录
func (s *SyncService) updateRecordTx(tx *gorm.DB, historyID uint, detail *dto.SyncDetailItem, existing *model.UploadRecord, record *dto.SyncRecordItem) (*dto.SyncDetailItem, error) {
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
	if err := tx.Save(existing).Error; err != nil {
		detail.Action = "update"
		detail.Result = "failed"
		detail.ErrorMsg = err.Error()
		s.saveDetailTx(tx, historyID, detail, record.SerialNo, "")
		return detail, err
	}
	detail.Action = "update"
	detail.Result = "success"
	s.saveDetailTx(tx, historyID, detail, record.SerialNo, "")
	return detail, nil
}

// createConflictedRecordTx 事务内创建冲突记录
func (s *SyncService) createConflictedRecordTx(tx *gorm.DB, historyID uint, detail *dto.SyncDetailItem, record *dto.SyncRecordItem, oldSerialNo string) (*dto.SyncDetailItem, error) {
	hash := md5.Sum([]byte(record.ProjectName))
	newSerialNo := fmt.Sprintf("%s_%s", record.SerialNo, hex.EncodeToString(hash[:])[:6])
	var existing model.UploadRecord
	if err := tx.Where("serial_no = ? AND is_deleted = ?", newSerialNo, false).First(&existing).Error; err == nil {
		randomBytes := make([]byte, 4)
		rand.Read(randomBytes)
		newSerialNo = fmt.Sprintf("%s_%s", newSerialNo, hex.EncodeToString(randomBytes))
	}
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
	if err := tx.Create(newRecord).Error; err != nil {
		detail.Action = "conflict"
		detail.Result = "failed"
		detail.ErrorMsg = err.Error()
		s.saveDetailTx(tx, historyID, detail, record.SerialNo, newSerialNo)
		return detail, err
	}
	detail.Action = "conflict"
	detail.Result = "success"
	detail.OldSerialNo = oldSerialNo
	detail.NewSerialNo = newSerialNo
	s.saveDetailTx(tx, historyID, detail, record.SerialNo, newSerialNo)
	return detail, nil
}

// saveDetailTx 事务内保存详情
func (s *SyncService) saveDetailTx(tx *gorm.DB, historyID uint, detail *dto.SyncDetailItem, serialNo, newSerialNo string) {
	syncDetail := &model.SyncDetail{
		HistoryID:   historyID,
		SerialNo:    serialNo,
		ProjectName: detail.ProjectName,
		Action:      detail.Action,
		Result:      detail.Result,
		ErrorMsg:    detail.ErrorMsg,
		OldSerialNo: detail.OldSerialNo,
		NewSerialNo: detail.NewSerialNo,
	}
	tx.Create(syncDetail)
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

// GetStationSummaries 获取各站点的同步汇总统计
func (s *SyncService) GetStationSummaries(req dto.SyncStationSummaryReq) ([]dto.SyncStationSummaryResp, error) {
	return s.historyRepo.GetStationSummaries(req.StationID, req.StartDate, req.EndDate)
}

// GetHistoryDetails 获取同步历史详情（支持分页和过滤）
func (s *SyncService) GetHistoryDetails(historyID uint, detailResult string, page, pageSize int) (*dto.SyncHistoryDetailResp, error) {
	history, err := s.historyRepo.GetByID(historyID)
	if err != nil {
		return nil, fmt.Errorf("同步历史不存在")
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	details, total, err := s.detailRepo.ListByHistoryIDFiltered(historyID, detailResult, page, pageSize)
	if err != nil {
		return nil, err
	}

	items := make([]dto.SyncDetailResp, len(details))
	for i, d := range details {
		items[i] = *s.toDetailResp(&d)
	}

	resp := &dto.SyncHistoryDetailResp{
		SyncHistoryResp: *s.toHistoryResp(history),
		Total:          total,
		Details:        items,
	}

	return resp, nil
}

// GetSyncStatus 获取同步状态（增强版，包含调度器完整状态）
func (s *SyncService) GetSyncStatus() *dto.SyncStatusResp {
	mode := global.CONF.Sync.Mode

	hbInterval := global.CONF.Sync.Agent.HeartbeatInterval
	if hbInterval == "" {
		hbInterval = "60s"
	}

	resp := &dto.SyncStatusResp{
		Enabled:           mode != "",
		Mode:              mode,
		Interval:          global.CONF.Sync.Agent.Interval,
		HeartbeatInterval: hbInterval,
		BatchSize:         global.CONF.Sync.Agent.BatchSize,
	}

	if mode == "agent" && global.CONF.Sync.Agent.StationID != "" {
		// Agent 模式
		resp.IsCenter = false
		resp.StationID = global.CONF.Sync.Agent.StationID
		resp.StationName = global.CONF.Sync.Agent.StationName
		resp.CenterURL = global.CONF.Sync.Agent.CenterURL

		// 从调度器获取注册状态
		if scheduler, ok := global.Scheduler.(*SyncScheduler); ok {
			status := scheduler.GetQueueStatus()
			resp.Registered, _ = status["registered"].(bool)
			resp.LastSyncAt, _ = status["lastSyncAt"].(*time.Time)
			resp.LastSerialNo, _ = status["lastSerialNo"].(string)
			resp.SyncQueueCount, _ = status["pending"].(int)
			resp.QueueTotal, _ = status["total"].(int)
			resp.QueuePending, _ = status["pending"].(int)
			resp.QueueCompleted, _ = status["completed"].(int)
			resp.QueueFailed, _ = status["failed"].(int)
			if t, ok := status["lastErrorAt"].(time.Time); ok {
				resp.LastErrorAt = &t
			}
			if filter := scheduler.GetFilter(); filter != nil && len(filter.ProjectNames) > 0 {
				resp.Filter = &dto.SyncFilterResp{ProjectNames: filter.ProjectNames}
			}
		}
	} else if mode == "center" {
		// Center 模式
		resp.IsCenter = true
		timeout := global.CONF.Sync.Center.HeartbeatTimeout
		if timeout <= 0 {
			timeout = 60
		}
		resp.HeartbeatTimeout = timeout
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
	// 主动探测结果优先（内存状态），其次用心跳时间保底
	connected := global.GetProbeConnected(station.ID)
	if !connected && station.LastHeartbeatAt != nil {
		timeout := global.CONF.Sync.Center.HeartbeatTimeout
		if timeout <= 0 {
			timeout = 60
		}
		connected = time.Since(*station.LastHeartbeatAt) < time.Duration(timeout)*time.Second
	}

	resp := &dto.SyncStationResp{
		ID:              station.ID,
		Name:            station.Name,
		Code:            station.Code,
		URL:             station.URL,
		Status:          station.Status,
		StatusText:      "正常",
		Description:     station.Description,
		IsCenter:        station.IsCenter,
		SyncCount:       station.SyncCount,
		Remark:          station.Remark,
		IsConnected:     connected,
	}
	if station.CreatedAt.Unix() > 0 {
		resp.CreatedAt = station.CreatedAt.UnixMilli()
	}
	if station.UpdatedAt.Unix() > 0 {
		resp.UpdatedAt = station.UpdatedAt.UnixMilli()
	}
	if station.LastSyncAt != nil && station.LastSyncAt.Unix() > 0 {
		resp.LastSyncAt = station.LastSyncAt.UnixMilli()
	}
	if station.LastHeartbeatAt != nil && station.LastHeartbeatAt.Unix() > 0 {
		resp.LastHeartbeatAt = station.LastHeartbeatAt.UnixMilli()
	}
	if station.LastConnectedAt != nil && station.LastConnectedAt.Unix() > 0 {
		resp.LastConnectedAt = station.LastConnectedAt.UnixMilli()
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
		ErrorMsg:      history.ErrorMsg,
		Remark:        history.Remark,
	}
	if history.CreatedAt.Unix() > 0 {
		resp.CreatedAt = history.CreatedAt.UnixMilli()
	}
	if history.StartedAt != nil && history.StartedAt.Unix() > 0 {
		resp.StartedAt = history.StartedAt.UnixMilli()
	}
	if history.CompletedAt != nil && history.CompletedAt.Unix() > 0 {
		resp.CompletedAt = history.CompletedAt.UnixMilli()
	}

	if history.Station != nil {
		resp.StationName = history.Station.Name
		resp.StationCode = history.Station.Code
	}

	return resp
}

// toDetailResp 转换为详情响应
func (s *SyncService) toDetailResp(detail *model.SyncDetail) *dto.SyncDetailResp {
	return &dto.SyncDetailResp{
		ID:          detail.ID,
		HistoryID:   detail.HistoryID,
		SerialNo:    detail.SerialNo,
		ProjectName: detail.ProjectName,
		Action:      detail.Action,
		ActionText:  s.getActionText(detail.Action),
		Result:      detail.Result,
		ErrorMsg:    detail.ErrorMsg,
		OldSerialNo: detail.OldSerialNo,
		NewSerialNo: detail.NewSerialNo,
		CreatedAt:   detail.CreatedAt.UnixMilli(),
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
