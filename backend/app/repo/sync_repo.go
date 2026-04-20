package repo

import (
	"datauptwo/app/dto"
	"datauptwo/app/model"
	"datauptwo/global"
	"time"
)

// ============ SyncStationRepo ============

type SyncStationRepo struct{}

func NewSyncStationRepo() *SyncStationRepo {
	return &SyncStationRepo{}
}

func (r *SyncStationRepo) Create(station *model.SyncStation) error {
	return global.DB.Create(station).Error
}

func (r *SyncStationRepo) Update(station *model.SyncStation) error {
	return global.DB.Save(station).Error
}

func (r *SyncStationRepo) Delete(id uint) error {
	return global.DB.Model(&model.SyncStation{}).Where("id = ?", id).Update("is_deleted", true).Error
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

func (r *SyncStationRepo) UpdateLastSync(id uint) error {
	now := time.Now()
	return global.DB.Model(&model.SyncStation{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_sync_at": now,
		"sync_count":   global.DB.Raw("sync_count + 1"),
	}).Error
}

func (r *SyncStationRepo) GetCenterStation() (*model.SyncStation, error) {
	var station model.SyncStation
	err := global.DB.Where("is_center = ? AND is_deleted = ?", true, false).First(&station).Error
	if err != nil {
		return nil, err
	}
	return &station, nil
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
	if err := db.Preload("Station").Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&histories).Error; err != nil {
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
		"success_count":    global.DB.Raw("success_count + ?", success),
		"fail_count":       global.DB.Raw("fail_count + ?", fail),
		"conflict_count":   global.DB.Raw("conflict_count + ?", conflict),
	}).Error
}

func (r *SyncHistoryRepo) SetStarted(id uint) error {
	now := time.Now()
	return global.DB.Model(&model.SyncHistory{}).Where("id = ?", id).Updates(map[string]interface{}{
		"started_at": &now,
		"status":     "processing",
	}).Error
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
