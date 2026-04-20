package repo

import (
	"datauptwo/app/dto"
	"datauptwo/app/model"
	"datauptwo/global"
	"time"

	"gorm.io/gorm"
)

type UploadRecordRepo struct{}

func NewUploadRecordRepo() *UploadRecordRepo {
	return &UploadRecordRepo{}
}

func (r *UploadRecordRepo) Create(record *model.UploadRecord) error {
	return global.DB.Create(record).Error
}

func (r *UploadRecordRepo) Update(record *model.UploadRecord) error {
	return global.DB.Save(record).Error
}

func (r *UploadRecordRepo) Delete(id uint) error {
	return global.DB.Model(&model.UploadRecord{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *UploadRecordRepo) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return global.DB.Model(&model.UploadRecord{}).Where("id IN ?", ids).Update("is_deleted", true).Error
}

// BatchUpdateStatus 批量更新状态
func (r *UploadRecordRepo) BatchUpdateStatus(ids []uint, status string) error {
	if len(ids) == 0 {
		return nil
	}
	return global.DB.Model(&model.UploadRecord{}).Where("id IN ?", ids).Updates(map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}).Error
}

func (r *UploadRecordRepo) GetByID(id uint) (*model.UploadRecord, error) {
	var record model.UploadRecord
	err := global.DB.Where("id = ? AND is_deleted = ?", id, false).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *UploadRecordRepo) GetBySerialNo(serialNo string) (*model.UploadRecord, error) {
	var record model.UploadRecord
	err := global.DB.Where("serial_no = ? AND is_deleted = ?", serialNo, false).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// ExistsBySerialNo 检查流水号是否已存在
func (r *UploadRecordRepo) ExistsBySerialNo(serialNo string) (bool, error) {
	var count int64
	err := global.DB.Model(&model.UploadRecord{}).
		Where("serial_no = ? AND is_deleted = ?", serialNo, false).
		Count(&count).Error
	return count > 0, err
}

// StatisticsFilter 统计筛选条件
type StatisticsFilter struct {
	ProjectName string
	DiskLabel    string
	StartDate   string
	EndDate     string
	Status      string
	Uploader    string
}

// buildFilteredDB 构建带筛选条件的查询
func (r *UploadRecordRepo) buildFilteredDB(filters ...StatisticsFilter) *gorm.DB {
	db := global.DB.Model(&model.UploadRecord{}).Where("is_deleted = ?", false)

	if len(filters) > 0 {
		f := filters[0]
		if f.ProjectName != "" {
			db = db.Where("project_name = ?", f.ProjectName)
		}
		if f.DiskLabel != "" {
			db = db.Where("disk_label = ?", f.DiskLabel)
		}
		if f.Status != "" {
			db = db.Where("status = ?", f.Status)
		}
		if f.Uploader != "" {
			db = db.Where("uploader = ?", f.Uploader)
		}
		if f.StartDate != "" {
			db = db.Where("date(created_at) >= date(?)", f.StartDate)
		}
		if f.EndDate != "" {
			db = db.Where("date(created_at) <= date(?)", f.EndDate)
		}
	}

	return db
}

// GetUploaderList 获取所有上传者列表
func (r *UploadRecordRepo) GetUploaderList() ([]string, error) {
	var uploaders []string
	err := global.DB.Model(&model.UploadRecord{}).
		Where("is_deleted = ? AND uploader != ''", false).
		Distinct("uploader").
		Order("uploader ASC").
		Pluck("uploader", &uploaders).Error
	return uploaders, err
}

func (r *UploadRecordRepo) List(req dto.UploadRecordListReq) ([]model.UploadRecord, int64, error) {
	var records []model.UploadRecord
	var total int64

	db := global.DB.Model(&model.UploadRecord{}).Where("is_deleted = ?", false)

	if req.DiskLabel != "" {
		db = db.Where("disk_label = ?", req.DiskLabel)
	}

	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	if req.Uploader != "" {
		db = db.Where("uploader LIKE ?", "%"+req.Uploader+"%")
	}

	if req.ProjectName != "" {
		db = db.Where("project_name = ?", req.ProjectName)
	}

	if req.StartDate != "" {
		db = db.Where("date(created_at) >= date(?)", req.StartDate)
	}

	if req.EndDate != "" {
		db = db.Where("date(created_at) <= date(?)", req.EndDate)
	}

	if req.Keyword != "" {
		db = db.Where("dest_path LIKE ? OR remark LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 计数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := db.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// CountByDate 统计指定日期的记录数
func (r *UploadRecordRepo) CountByDate(date string, filters ...StatisticsFilter) (int64, error) {
	db := global.DB.Model(&model.UploadRecord{}).Where("is_deleted = ?", false)

	if len(filters) > 0 {
		f := filters[0]
		if f.ProjectName != "" {
			db = db.Where("project_name = ?", f.ProjectName)
		}
		if f.DiskLabel != "" {
			db = db.Where("disk_label = ?", f.DiskLabel)
		}
		if f.Status != "" {
			db = db.Where("status = ?", f.Status)
		}
		if f.Uploader != "" {
			db = db.Where("uploader = ?", f.Uploader)
		}
	}

	var count int64
	err := db.Where("date(created_at) = date(?)", date).
		Select("COUNT(*)").Scan(&count).Error
	return count, err
}
func (r *UploadRecordRepo) CountTotal(filters ...StatisticsFilter) (int64, error) {
	db := global.DB.Model(&model.UploadRecord{}).Where("is_deleted = ?", false)

	if len(filters) > 0 {
		f := filters[0]
		if f.ProjectName != "" {
			db = db.Where("project_name = ?", f.ProjectName)
		}
		if f.DiskLabel != "" {
			db = db.Where("disk_label = ?", f.DiskLabel)
		}
		if f.Status != "" {
			db = db.Where("status = ?", f.Status)
		}
		if f.Uploader != "" {
			db = db.Where("uploader = ?", f.Uploader)
		}
	}

	var count int64
	err := db.Select("COUNT(*)").Scan(&count).Error
	return count, err
}

// CountByDateRange 统计日期范围内的记录数
func (r *UploadRecordRepo) CountByDateRange(startDate, endDate string, filters ...StatisticsFilter) (int64, error) {
	db := global.DB.Model(&model.UploadRecord{}).Where("is_deleted = ?", false)

	if len(filters) > 0 {
		f := filters[0]
		if f.ProjectName != "" {
			db = db.Where("project_name = ?", f.ProjectName)
		}
		if f.DiskLabel != "" {
			db = db.Where("disk_label = ?", f.DiskLabel)
		}
		if f.Status != "" {
			db = db.Where("status = ?", f.Status)
		}
		if f.Uploader != "" {
			db = db.Where("uploader = ?", f.Uploader)
		}
	}

	var count int64
	err := db.Where("date(created_at) >= date(?) AND date(created_at) <= date(?)", startDate, endDate).
		Select("COUNT(*)").Scan(&count).Error
	return count, err
}

// SumFileSizeByDate 统计指定日期的文件大小
func (r *UploadRecordRepo) SumFileSizeByDate(date string, filters ...StatisticsFilter) (int64, error) {
	var total int64
	db := r.buildFilteredDB(filters...)
	err := db.Where("date(created_at) = date(?)", date).
		Select("COALESCE(SUM(file_size), 0)").
		Scan(&total).Error
	return total, err
}

// SumFileSizeByDateRange 统计日期范围内的文件大小
func (r *UploadRecordRepo) SumFileSizeByDateRange(startDate, endDate string, filters ...StatisticsFilter) (int64, error) {
	var total int64
	db := r.buildFilteredDB(filters...)
	err := db.Where("date(created_at) >= date(?) AND date(created_at) <= date(?)", startDate, endDate).
		Select("COALESCE(SUM(file_size), 0)").
		Scan(&total).Error
	return total, err
}

// CountByStatus 统计各状态的记录数
func (r *UploadRecordRepo) CountByStatus(filters ...StatisticsFilter) ([]dto.StatusCount, error) {
	var results []dto.StatusCount
	err := r.buildFilteredDB(filters...).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error
	return results, err
}

// CountByDiskLabel 统计各磁盘标签的记录数和数据量
func (r *UploadRecordRepo) CountByDiskLabel(filters ...StatisticsFilter) ([]dto.DiskLabelCount, error) {
	var results []dto.DiskLabelCount
	err := r.buildFilteredDB(filters...).
		Select("disk_label, COUNT(*) as count, COALESCE(SUM(file_size), 0) as total_size").
		Group("disk_label").
		Scan(&results).Error
	return results, err
}

// GetDiskLabelStatusAll 获取所有磁盘标签及其状态（支持日期范围）
func (r *UploadRecordRepo) GetDiskLabelStatusAll(startDate, endDate string) ([]dto.DiskLabelStatus, error) {
	var results []dto.DiskLabelStatus
	db := global.DB.Model(&model.UploadRecord{}).
		Select("disk_label, COUNT(*) as count").
		Where("is_deleted = 0 AND disk_label != '' AND disk_label IS NOT NULL")
	if startDate != "" {
		db = db.Where("created_at >= ?", startDate+" 00:00:00")
	}
	if endDate != "" {
		db = db.Where("created_at <= ?", endDate+" 23:59:59")
	}
	err := db.Group("disk_label").Order("count DESC").Scan(&results).Error
	return results, err
}

// GetDiskLabelStatusDetail 获取每个磁盘标签下各状态的详细数量（支持日期范围）
func (r *UploadRecordRepo) GetDiskLabelStatusDetail(startDate, endDate string) (map[string]map[string]int64, error) {
	type result struct {
		DiskLabel string `gorm:"column:disk_label"`
		Status    string `gorm:"column:status"`
		Count     int64  `gorm:"column:count"`
	}
	var rows []result
	db := global.DB.Model(&model.UploadRecord{}).
		Select("disk_label, status, COUNT(*) as count").
		Where("is_deleted = 0 AND disk_label != '' AND disk_label IS NOT NULL")
	if startDate != "" {
		db = db.Where("created_at >= ?", startDate+" 00:00:00")
	}
	if endDate != "" {
		db = db.Where("created_at <= ?", endDate+" 23:59:59")
	}
	err := db.Group("disk_label, status").Order("disk_label").Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	detail := make(map[string]map[string]int64)
	for _, row := range rows {
		if detail[row.DiskLabel] == nil {
			detail[row.DiskLabel] = make(map[string]int64)
		}
		detail[row.DiskLabel][row.Status] = row.Count
	}
	return detail, nil
}

// CountByProject 统计各项目的记录数和数据量（应用日期范围等全部筛选条件）
func (r *UploadRecordRepo) CountByProject(filters ...StatisticsFilter) ([]dto.ProjectCount, error) {
	var results []dto.ProjectCount
	db := r.buildFilteredDB(filters...)

	query := `
		SELECT COALESCE(NULLIF(TRIM(project_name), ''), '(空项目)') as project_name,
		       COUNT(*) as count,
		       COALESCE(SUM(file_size), 0) as total_size
		FROM upload_records
		WHERE is_deleted = 0
	`
	args := []interface{}{}

	if len(filters) > 0 {
		f := filters[0]
		if f.ProjectName != "" {
			query += " AND project_name = ?"
			args = append(args, f.ProjectName)
		}
		if f.DiskLabel != "" {
			query += " AND disk_label = ?"
			args = append(args, f.DiskLabel)
		}
		if f.Status != "" {
			query += " AND status = ?"
			args = append(args, f.Status)
		}
		if f.Uploader != "" {
			query += " AND uploader = ?"
			args = append(args, f.Uploader)
		}
		if f.StartDate != "" {
			query += " AND date(created_at) >= date(?)"
			args = append(args, f.StartDate)
		}
		if f.EndDate != "" {
			query += " AND date(created_at) <= date(?)"
			args = append(args, f.EndDate)
		}
	}

	query += " GROUP BY COALESCE(NULLIF(TRIM(project_name), ''), '(空项目)') ORDER BY count DESC"

	rows, err := db.Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item dto.ProjectCount
		if err := rows.Scan(&item.ProjectName, &item.Count, &item.TotalSize); err != nil {
			continue
		}
		results = append(results, item)
	}
	return results, nil
}


// GetTrend 获取每日趋势
func (r *UploadRecordRepo) GetTrend(startDate, endDate string, filters ...StatisticsFilter) ([]dto.DailyTrend, error) {
	// 去掉 is_deleted=0 过滤，累计曲线反映所有上传行为
	// 用子查询实现累计：外层取日期列表，内层累计该日期之前（含）所有记录
	filterSQL := ""
	filterArgs := []interface{}{}
	if len(filters) > 0 {
		f := filters[0]
		if f.ProjectName != "" {
			filterSQL += " AND project_name = ?"
			filterArgs = append(filterArgs, f.ProjectName)
		}
		if f.DiskLabel != "" {
			filterSQL += " AND disk_label = ?"
			filterArgs = append(filterArgs, f.DiskLabel)
		}
	}

	// 子查询 filterSQL 参数与日期参数交叉追加（3个子查询各需一份）
	allArgs := []interface{}{startDate, endDate + " 23:59:59"} // date subquery args
	allArgs = append(allArgs, filterArgs...)                  // date subquery filter
	allArgs = append(allArgs, filterArgs...)                  // count subquery filter
	allArgs = append(allArgs, filterArgs...)                  // size subquery filter

	query := `
		SELECT d.date,
		       (SELECT COUNT(*) FROM upload_records WHERE date(created_at) <= d.date` + filterSQL + `) as count,
		       (SELECT COALESCE(SUM(file_size), 0) FROM upload_records WHERE date(created_at) <= d.date` + filterSQL + `) as total_size
		FROM (SELECT DISTINCT date(created_at) as date FROM upload_records WHERE date(created_at) >= date(?) AND date(created_at) <= date(?)` + filterSQL + `) d
		ORDER BY d.date ASC
	`

	rows, err := global.DB.Raw(query, allArgs...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trends []dto.DailyTrend
	for rows.Next() {
		var item dto.DailyTrend
		if err := rows.Scan(&item.Date, &item.Count, &item.TotalSize); err != nil {
			continue
		}
		trends = append(trends, item)
	}
	return trends, nil
}


// GetTodayCount 获取今日记录数
func (r *UploadRecordRepo) GetTodayCount(filters ...StatisticsFilter) (int64, error) {
	today := time.Now().Format("2006-01-02")
	return r.CountByDate(today, filters...)
}

// GetTodayFileSize 获取今日文件大小
func (r *UploadRecordRepo) GetTodayFileSize(filters ...StatisticsFilter) (int64, error) {
	today := time.Now().Format("2006-01-02")
	return r.SumFileSizeByDate(today, filters...)
}

// GetRecentRecords 获取最近上传记录
func (r *UploadRecordRepo) GetRecentRecords(limit int, filters ...StatisticsFilter) ([]model.UploadRecord, error) {
	var records []model.UploadRecord
	db := r.buildFilteredDB(filters...)
	err := db.Order("created_at DESC").Limit(limit).Find(&records).Error
	return records, err
}

// ListAllForExport 获取所有符合条件的记录（用于导出，不分页）
func (r *UploadRecordRepo) ListAllForExport(req dto.UploadRecordListReq) ([]model.UploadRecord, error) {
	var records []model.UploadRecord

	db := global.DB.Model(&model.UploadRecord{}).Where("is_deleted = ?", false)

	if req.DiskLabel != "" {
		db = db.Where("disk_label = ?", req.DiskLabel)
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	if req.Uploader != "" {
		db = db.Where("uploader LIKE ?", "%"+req.Uploader+"%")
	}
	if req.ProjectName != "" {
		db = db.Where("project_name = ?", req.ProjectName)
	}
	if req.StartDate != "" {
		db = db.Where("date(created_at) >= date(?)", req.StartDate)
	}
	if req.EndDate != "" {
		db = db.Where("date(created_at) <= date(?)", req.EndDate)
	}
	if req.Keyword != "" {
		db = db.Where("dest_path LIKE ? OR remark LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if err := db.Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

// GetRecordsSince 获取指定时间之后或指定SerialNo之后的记录（用于增量同步）
func (r *UploadRecordRepo) GetRecordsSince(lastSyncAt *time.Time, lastSerialNo string, projectNames []string, limit int) ([]model.UploadRecord, error) {
	var records []model.UploadRecord

	db := global.DB.Model(&model.UploadRecord{}).Where("is_deleted = ?", false)

	// 项目过滤
	if len(projectNames) > 0 {
		db = db.Where("project_name IN ?", projectNames)
	}

	// 时间过滤（断点恢复）
	if lastSyncAt != nil {
		db = db.Where("created_at > ?", lastSyncAt)
	}

	// 或者 SerialNo 过滤（用于同时间的断点）
	if lastSerialNo != "" {
		db = db.Where("serial_no > ?", lastSerialNo)
	}

	err := db.Order("created_at ASC, serial_no ASC").Limit(limit).Find(&records).Error
	return records, err
}
