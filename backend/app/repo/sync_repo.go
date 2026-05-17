package repo

import (
	"datauptwo/app/dto"
	"datauptwo/app/model"
	"datauptwo/global"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ============ SyncStationRepo ============

type SyncStationRepo struct{}

func NewSyncStationRepo() *SyncStationRepo {
	return &SyncStationRepo{}
}

func (r *SyncStationRepo) Create(station *model.SyncStation) error {
	return global.DB.Create(station).Error
}

// UpsertByCode 原子插入或更新（code 冲突时更新 name/url），用于 Agent 自动注册
func (r *SyncStationRepo) UpsertByCode(station *model.SyncStation) error {
	return global.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "url", "api_key", "updated_at"}),
	}).Create(station).Error
}

func (r *SyncStationRepo) Update(station *model.SyncStation) error {
	return global.DB.Save(station).Error
}

func (r *SyncStationRepo) Delete(id uint) error {
	return global.DB.Model(&model.SyncStation{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// UpdateHeartbeatByCode 更新心跳时间（按 code 查找）
func (r *SyncStationRepo) UpdateHeartbeatByCode(code string) error {
	now := time.Now()
	return global.DB.Model(&model.SyncStation{}).
		Where("code = ? AND is_deleted = ?", code, false).
		Updates(map[string]interface{}{
			"last_heartbeat_at": now,
			"status":            "active",
		}).Error
}

func (r *SyncStationRepo) GetByID(id uint) (*model.SyncStation, error) {
	var station model.SyncStation
	err := global.DB.Where("id = ? AND is_deleted = ?", id, false).First(&station).Error
	if err != nil {
		return nil, err
	}
	return &station, nil
}

func (r *SyncStationRepo) GetByCode(code string) (*model.SyncStation, error) {
	var station model.SyncStation
	err := global.DB.Where("code = ? AND is_deleted = ?", code, false).First(&station).Error
	if err != nil {
		return nil, err
	}
	return &station, nil
}

func (r *SyncStationRepo) GetByAPIKey(apiKey string) (*model.SyncStation, error) {
	var station model.SyncStation
	err := global.DB.Where("api_key = ? AND is_deleted = ?", apiKey, false).First(&station).Error
	if err != nil {
		return nil, err
	}
	return &station, nil
}

func (r *SyncStationRepo) List(req dto.SyncStationListReq) ([]model.SyncStation, int64, error) {
	var stations []model.SyncStation
	var total int64

	db := global.DB.Model(&model.SyncStation{}).Where("is_deleted = ?", false)

	if req.Keyword != "" {
		db = db.Where("name LIKE ? OR code LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := db.Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&stations).Error; err != nil {
		return nil, 0, err
	}

	return stations, total, nil
}

func (r *SyncStationRepo) GetAll() ([]model.SyncStation, error) {
	var stations []model.SyncStation
	err := global.DB.Where("is_deleted = ? AND status = ?", false, "active").Order("id ASC").Find(&stations).Error
	return stations, err
}

// GetAllForProbe 获取所有非删除站点（用于主动探测，包含 inactive）
func (r *SyncStationRepo) GetAllForProbe() ([]model.SyncStation, error) {
	var stations []model.SyncStation
	err := global.DB.Where("is_deleted = ? AND is_center = ?", false, false).Order("id ASC").Find(&stations).Error
	return stations, err
}

func (r *SyncStationRepo) UpdateLastSync(id uint) error {
	now := time.Now()
	var m model.SyncStation
	if err := global.DB.Where("id = ?", id).First(&m).Error; err != nil {
		return err
	}
	m.LastSyncAt = &now
	m.SyncCount++
	return global.DB.Save(&m).Error
}

// UpdateCheckpoint 更新断点（lastSyncAt + lastSerialNo）
func (r *SyncStationRepo) UpdateCheckpoint(id uint, lastSyncAt *time.Time, lastSerialNo string) error {
	return global.DB.Model(&model.SyncStation{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_sync_at":    lastSyncAt,
		"last_serial_no":  lastSerialNo,
	}).Error
}

// UpdatePlainAPIKey 更新明文 API Key（仅 Agent 使用）
func (r *SyncStationRepo) UpdatePlainAPIKey(id uint, plainKey string) error {
	return global.DB.Model(&model.SyncStation{}).Where("id = ?", id).Updates(map[string]interface{}{
		"plain_api_key": plainKey,
	}).Error
}

// GetPlainAPIKey 获取明文 API Key
func (r *SyncStationRepo) GetPlainAPIKey(id uint) (string, error) {
	var m model.SyncStation
	if err := global.DB.Where("id = ? AND is_deleted = ?", id, false).First(&m).Error; err != nil {
		return "", err
	}
	return m.PlainAPIKey, nil
}

func (r *SyncStationRepo) GetCenterStation() (*model.SyncStation, error) {
	var station model.SyncStation
	err := global.DB.Where("is_center = ? AND is_deleted = ?", true, false).First(&station).Error
	if err != nil {
		// Center 模式下 Center 自身不在 sync_stations 表，视为正常
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &station, nil
}

// UpdateLastConnectedAt 更新站点最后探测成功时间
func (r *SyncStationRepo) UpdateLastConnectedAt(id uint, t *time.Time) error {
	return global.DB.Model(&model.SyncStation{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_connected_at": t,
	}).Error
}

// ============ SyncHistoryRepo ============

type SyncHistoryRepo struct{}

func NewSyncHistoryRepo() *SyncHistoryRepo {
	return &SyncHistoryRepo{}
}

func (r *SyncHistoryRepo) Create(history *model.SyncHistory) error {
	return global.DB.Create(history).Error
}

func (r *SyncHistoryRepo) Update(history *model.SyncHistory) error {
	return global.DB.Save(history).Error
}

func (r *SyncHistoryRepo) GetByID(id uint) (*model.SyncHistory, error) {
	var history model.SyncHistory
	err := global.DB.Preload("Station").Where("id = ?", id).First(&history).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

func (r *SyncHistoryRepo) List(req dto.SyncHistoryListReq) ([]model.SyncHistory, int64, error) {
	var histories []model.SyncHistory
	var total int64

	db := global.DB.Model(&model.SyncHistory{})

	if req.Keyword != "" {
		// 关键词搜索：备注 / 站点名称 / 流水号
		// 需要显式 JOIN sync_stations，因为 WHERE 中引用了它
		db = db.Joins("LEFT JOIN sync_stations ON sync_stations.id = sync_histories.station_id AND sync_stations.is_deleted = ?", false)
		keyword := "%" + req.Keyword + "%"
		db = db.Where(
			"sync_histories.remark LIKE ? OR sync_stations.name LIKE ? OR sync_histories.id IN (SELECT history_id FROM sync_details WHERE serial_no LIKE ?)",
			keyword, keyword, keyword,
		)
	} else {
		// 无关键词时用 Preload 加载站点（性能更好）
		db = db.Preload("Station")
	}

	if req.StationID > 0 {
		db = db.Where("station_id = ?", req.StationID)
	}
	if req.Direction != "" {
		db = db.Where("direction = ?", req.Direction)
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	if req.StartDate != "" {
		db = db.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		db = db.Where("created_at <= ?", req.EndDate+" 23:59:59")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := db.Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

func (r *SyncHistoryRepo) UpdateStatus(id uint, status string, errorMsg string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if errorMsg != "" {
		updates["error_msg"] = errorMsg
	}
	if status == "completed" || status == "failed" {
		now := time.Now()
		updates["completed_at"] = &now
	}
	return global.DB.Model(&model.SyncHistory{}).Where("id = ?", id).Updates(updates).Error
}

func (r *SyncHistoryRepo) IncrementCount(id uint, success, fail, conflict int) error {
	return global.DB.Model(&model.SyncHistory{}).Where("id = ?", id).Updates(map[string]interface{}{
		"success_count":    gorm.Expr("success_count + ?", success),
		"fail_count":       gorm.Expr("fail_count + ?", fail),
		"conflict_count":   gorm.Expr("conflict_count + ?", conflict),
	}).Error
}

func (r *SyncHistoryRepo) SetStarted(id uint) error {
	now := time.Now()
	return global.DB.Model(&model.SyncHistory{}).Where("id = ?", id).Updates(map[string]interface{}{
		"started_at": &now,
		"status":     "processing",
	}).Error
}

// GetStationSummaries 按站点聚合同步统计（支持按站点ID和日期范围过滤）
func (r *SyncHistoryRepo) GetStationSummaries(stationId uint, startDate, endDate string) ([]dto.SyncStationSummaryResp, error) {
	type row struct {
		StationID     uint
		StationName  string
		Code         string
		Status       string
		LastSyncAt   int64
		TotalSyncs   int
		TotalRecords int
		SuccessCount int
		FailCount    int
		ConflictCount int
	}

	// 构建数据库兼容的 MAX(completed_at) Unix ms 表达式
	var lastSyncExpr string
	switch global.DBType() {
	case "mysql":
		lastSyncExpr = "CAST(UNIX_TIMESTAMP(MAX(completed_at)) * 1000 AS BIGINT)"
	case "pgsql":
		lastSyncExpr = "CAST(EXTRACT(EPOCH FROM MAX(completed_at)) * 1000 AS BIGINT)"
	default: // sqlite
		lastSyncExpr = "CAST(MAX(CAST(strftime('%s', completed_at) AS INTEGER)) * 1000 AS INTEGER)"
	}

	selectClause := fmt.Sprintf(
		`sync_histories.station_id,
		sync_stations.name as station_name,
		sync_stations.code,
		sync_stations.status,
		%s as last_sync_at,
		COUNT(sync_histories.id) as total_syncs,
		COALESCE(SUM(sync_histories.total_records), 0) as total_records,
		COALESCE(SUM(sync_histories.success_count), 0) as success_count,
		COALESCE(SUM(sync_histories.fail_count), 0) as fail_count,
		COALESCE(SUM(sync_histories.conflict_count), 0) as conflict_count`,
		lastSyncExpr,
	)

	db := global.DB.Table("sync_histories").
		Select(selectClause).
		Joins("LEFT JOIN sync_stations ON sync_stations.id = sync_histories.station_id AND sync_stations.is_deleted = ?", false).
		Group("sync_histories.station_id, sync_stations.name, sync_stations.code, sync_stations.status").
		Order("last_sync_at DESC")

	if stationId > 0 {
		db = db.Where("sync_histories.station_id = ?", stationId)
	}
	if startDate != "" {
		db = db.Where("sync_histories.created_at >= ?", startDate)
	}
	if endDate != "" {
		db = db.Where("sync_histories.created_at <= ?", endDate+" 23:59:59")
	}

	var rows []row
	if err := db.Scan(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]dto.SyncStationSummaryResp, len(rows))
	for i, r := range rows {
		result[i] = dto.SyncStationSummaryResp{
			ID:            r.StationID,
			Name:          r.StationName,
			Code:          r.Code,
			Status:        r.Status,
			LastSyncAt:    r.LastSyncAt,
			TotalSyncs:    r.TotalSyncs,
			TotalRecords:  r.TotalRecords,
			SuccessCount:  r.SuccessCount,
			FailCount:     r.FailCount,
			ConflictCount: r.ConflictCount,
		}
	}
	return result, nil
}

// ============ SyncDetailRepo ============

type SyncDetailRepo struct{}

func NewSyncDetailRepo() *SyncDetailRepo {
	return &SyncDetailRepo{}
}

func (r *SyncDetailRepo) Create(detail *model.SyncDetail) error {
	return global.DB.Create(detail).Error
}

func (r *SyncDetailRepo) CreateBatch(details []model.SyncDetail) error {
	if len(details) == 0 {
		return nil
	}
	return global.DB.CreateInBatches(details, 100).Error
}

func (r *SyncDetailRepo) GetByID(id uint) (*model.SyncDetail, error) {
	var detail model.SyncDetail
	err := global.DB.Where("id = ?", id).First(&detail).Error
	if err != nil {
		return nil, err
	}
	return &detail, nil
}

func (r *SyncDetailRepo) ListByHistoryID(historyID uint) ([]model.SyncDetail, error) {
	var details []model.SyncDetail
	err := global.DB.Where("history_id = ?", historyID).Order("id ASC").Find(&details).Error
	return details, err
}

func (r *SyncDetailRepo) CountByHistoryID(historyID uint) (int64, error) {
	var count int64
	err := global.DB.Model(&model.SyncDetail{}).Where("history_id = ?", historyID).Count(&count).Error
	return count, err
}

// ListByHistoryIDFiltered 获取同步详情（支持分页和结果过滤）
func (r *SyncDetailRepo) ListByHistoryIDFiltered(historyID uint, result string, page, pageSize int) ([]model.SyncDetail, int64, error) {
	var details []model.SyncDetail
	var total int64

	db := global.DB.Model(&model.SyncDetail{}).Where("history_id = ?", historyID)
	if result != "" && result != "all" {
		db = db.Where("result = ?", result)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("id ASC").Find(&details).Error; err != nil {
		return nil, 0, err
	}
	return details, total, nil
}
