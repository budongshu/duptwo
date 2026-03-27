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
	DataType    string
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
		if f.DataType != "" {
			db = db.Where("data_type = ?", f.DataType)
		}
		if f.Status != "" {
			db = db.Where("status = ?", f.Status)
		}
		if f.Uploader != "" {
			db = db.Where("uploader = ?", f.Uploader)
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

	if req.DataType != "" {
		db = db.Where("data_type = ?", req.DataType)
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
		db = db.Where("created_at >= ?", req.StartDate)
	}

	if req.EndDate != "" {
		db = db.Where("created_at <= ?", req.EndDate+" 23:59:59")
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
	var count int64
	db := r.buildFilteredDB(filters...)
	err := db.Where("DATE(created_at) = ?", date).Count(&count).Error
	return count, err
}

// CountTotal 统计总记录数
func (r *UploadRecordRepo) CountTotal(filters ...StatisticsFilter) (int64, error) {
	var count int64
	err := r.buildFilteredDB(filters...).Count(&count).Error
	return count, err
}

// CountByDateRange 统计日期范围内的记录数
func (r *UploadRecordRepo) CountByDateRange(startDate, endDate string, filters ...StatisticsFilter) (int64, error) {
	var count int64
	db := r.buildFilteredDB(filters...)
	err := db.Where("DATE(created_at) >= ? AND DATE(created_at) <= ?", startDate, endDate).Count(&count).Error
	return count, err
}

// SumFileSizeByDate 统计指定日期的文件大小
func (r *UploadRecordRepo) SumFileSizeByDate(date string, filters ...StatisticsFilter) (int64, error) {
	var total int64
	db := r.buildFilteredDB(filters...)
	err := db.Where("DATE(created_at) = ?", date).
		Select("COALESCE(SUM(file_size), 0)").
		Scan(&total).Error
	return total, err
}

// SumFileSizeByDateRange 统计日期范围内的文件大小
func (r *UploadRecordRepo) SumFileSizeByDateRange(startDate, endDate string, filters ...StatisticsFilter) (int64, error) {
	var total int64
	db := r.buildFilteredDB(filters...)
	err := db.Where("DATE(created_at) >= ? AND DATE(created_at) <= ?", startDate, endDate).
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

// CountByDataType 统计各数据类型的记录数
func (r *UploadRecordRepo) CountByDataType(filters ...StatisticsFilter) ([]dto.DataTypeCount, error) {
	var results []dto.DataTypeCount
	err := r.buildFilteredDB(filters...).
		Select("data_type, COUNT(*) as count").
		Group("data_type").
		Scan(&results).Error
	return results, err
}

// CountByProject 统计各项目的记录数和数据量
func (r *UploadRecordRepo) CountByProject(filters ...StatisticsFilter) ([]dto.ProjectCount, error) {
	var results []dto.ProjectCount
	db := r.buildFilteredDB(filters...)

	// 处理 project_name 为空的情况（空白项目也算一个分组）
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
		if f.DataType != "" {
			query += " AND data_type = ?"
			args = append(args, f.DataType)
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

	return results, err
}

// GetTrend 获取每日趋势
func (r *UploadRecordRepo) GetTrend(startDate, endDate string, filters ...StatisticsFilter) ([]dto.DailyTrend, error) {
	var trends []dto.DailyTrend

	query := `
		SELECT DATE(created_at) as date,
		       COUNT(*) as count,
		       COALESCE(SUM(file_size), 0) as total_size
		FROM upload_records
		WHERE is_deleted = 0
		  AND created_at >= ?
		  AND created_at <= ?
	`
	args := []interface{}{startDate, endDate + " 23:59:59"}

	if len(filters) > 0 {
		f := filters[0]
		if f.ProjectName != "" {
			query += " AND project_name = ?"
			args = append(args, f.ProjectName)
		}
		if f.DataType != "" {
			query += " AND data_type = ?"
			args = append(args, f.DataType)
		}
	}

	query += " GROUP BY DATE(created_at) ORDER BY date ASC"

	rows, err := global.DB.Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

	if req.DataType != "" {
		db = db.Where("data_type = ?", req.DataType)
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
		db = db.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		db = db.Where("created_at <= ?", req.EndDate+" 23:59:59")
	}
	if req.Keyword != "" {
		db = db.Where("dest_path LIKE ? OR remark LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if err := db.Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}
