package repo

import (
	"datauptwo/app/dto"
	"datauptwo/app/model"
	"datauptwo/global"
)

type ProjectRepo struct{}

func NewProjectRepo() *ProjectRepo {
	return &ProjectRepo{}
}

// Create 创建项目
func (r *ProjectRepo) Create(project *model.Project) error {
	return global.DB.Create(project).Error
}

// Update 更新项目
func (r *ProjectRepo) Update(project *model.Project) error {
	return global.DB.Save(project).Error
}

// Delete 软删除
func (r *ProjectRepo) Delete(id uint) error {
	return global.DB.Model(&model.Project{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// BatchDelete 批量软删除
func (r *ProjectRepo) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return global.DB.Model(&model.Project{}).Where("id IN ?", ids).Update("is_deleted", true).Error
}

// GetByID 获取单个项目
func (r *ProjectRepo) GetByID(id uint) (*model.Project, error) {
	var project model.Project
	err := global.DB.Where("id = ? AND is_deleted = ?", id, false).First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// GetByCode 根据编码获取
func (r *ProjectRepo) GetByCode(code string) (*model.Project, error) {
	var project model.Project
	err := global.DB.Where("code = ? AND is_deleted = ?", code, false).First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// ExistsByCode 检查编码是否存在
func (r *ProjectRepo) ExistsByCode(code string) (bool, error) {
	var count int64
	err := global.DB.Model(&model.Project{}).
		Where("code = ? AND is_deleted = ?", code, false).
		Count(&count).Error
	return count > 0, err
}

// ExistsByName 检查名称是否存在
func (r *ProjectRepo) ExistsByName(name string) (bool, error) {
	var count int64
	err := global.DB.Model(&model.Project{}).
		Where("name = ? AND is_deleted = ?", name, false).
		Count(&count).Error
	return count > 0, err
}

// List 分页列表
func (r *ProjectRepo) List(req dto.ProjectListReq) ([]model.Project, int64, error) {
	var projects []model.Project
	var total int64

	db := global.DB.Model(&model.Project{}).Where("is_deleted = ?", false)

	if req.Keyword != "" {
		db = db.Where("name LIKE ? OR code LIKE ? OR description LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	if req.Stage != "" {
		db = db.Where("stage = ?", req.Stage)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := db.Offset(offset).Limit(req.PageSize).Order("sort ASC, created_at DESC").Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

// ListAll 获取所有项目（用于下拉选择）
func (r *ProjectRepo) ListAll() ([]model.Project, error) {
	var projects []model.Project
	err := global.DB.Where("is_deleted = ? AND status = ?", false, "active").
		Order("sort ASC, created_at DESC").
		Find(&projects).Error
	return projects, err
}

// ListAllForKanban 获取所有项目（看板视图，不分页）
func (r *ProjectRepo) ListAllForKanban() ([]model.Project, error) {
	var projects []model.Project
	err := global.DB.Where("is_deleted = ?", false).
		Order("stage ASC, sort ASC, created_at DESC").
		Find(&projects).Error
	return projects, err
}

// ListAllSimple 获取所有项目（简单列表，不关联记录数）
func (r *ProjectRepo) ListAllSimple() ([]dto.ProjectSimpleResp, error) {
	var projects []dto.ProjectSimpleResp
	err := global.DB.Model(&model.Project{}).
		Select("id, name, code").
		Where("is_deleted = ? AND status = ?", false, "active").
		Order("sort ASC, created_at DESC").
		Scan(&projects).Error
	return projects, err
}

// CountRecordsByProjectID 统计某项目的上传记录数
func (r *ProjectRepo) CountRecordsByProjectID(projectName string) (int64, error) {
	var count int64
	err := global.DB.Model(&model.UploadRecord{}).
		Where("project_name = ? AND is_deleted = ?", projectName, false).
		Count(&count).Error
	return count, err
}

// SumFileSizeByProjectID 统计某项目的上传数据总量(字节)
func (r *ProjectRepo) SumFileSizeByProjectID(projectName string) (int64, error) {
	var size int64
	err := global.DB.Model(&model.UploadRecord{}).
		Where("project_name = ? AND is_deleted = ?", projectName, false).
		Select("COALESCE(SUM(file_size), 0)").
		Scan(&size).Error
	return size, err
}

// ProjectStats 项目统计（记录数 + 数据总量）
type ProjectStats struct {
	ProjectName string
	RecordCount int64
	TotalSize   int64
}

// BatchGetProjectStats 批量获取项目统计
func (r *ProjectRepo) BatchGetProjectStats(projectNames []string) (map[string]*ProjectStats, error) {
	result := make(map[string]*ProjectStats)
	if len(projectNames) == 0 {
		return result, nil
	}
	var rows []ProjectStats
	err := global.DB.Model(&model.UploadRecord{}).
		Select("project_name, COUNT(*) as record_count, COALESCE(SUM(file_size), 0) as total_size").
		Where("project_name IN ? AND is_deleted = ?", projectNames, false).
		Group("project_name").
		Scan(&rows).Error
	if err != nil {
		return result, err
	}
	for _, row := range rows {
		r := row
		result[row.ProjectName] = &r
	}
	return result, nil
}
