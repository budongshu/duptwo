package repo

import (
	"datauptwo/app/dto"
	"datauptwo/app/model"
	"datauptwo/global"
)

type PersonnelRepo struct{}

func NewPersonnelRepo() *PersonnelRepo {
	return &PersonnelRepo{}
}

// Create 创建人员
func (r *PersonnelRepo) Create(personnel *model.Personnel) error {
	return global.DB.Create(personnel).Error
}

// Update 更新人员
func (r *PersonnelRepo) Update(personnel *model.Personnel) error {
	return global.DB.Save(personnel).Error
}

// Delete 软删除
func (r *PersonnelRepo) Delete(id uint) error {
	return global.DB.Model(&model.Personnel{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// BatchDelete 批量软删除
func (r *PersonnelRepo) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return global.DB.Model(&model.Personnel{}).Where("id IN ?", ids).Update("is_deleted", true).Error
}

// GetByID 获取单个人员
func (r *PersonnelRepo) GetByID(id uint) (*model.Personnel, error) {
	var personnel model.Personnel
	err := global.DB.Where("id = ? AND is_deleted = ?", id, false).First(&personnel).Error
	if err != nil {
		return nil, err
	}
	return &personnel, nil
}

// ExistsByName 检查姓名是否存在
func (r *PersonnelRepo) ExistsByName(name string) (bool, error) {
	var count int64
	err := global.DB.Model(&model.Personnel{}).
		Where("name = ? AND is_deleted = ?", name, false).
		Count(&count).Error
	return count > 0, err
}

// ExistsByNameExcludeId 检查姓名是否存在（排除指定ID）
func (r *PersonnelRepo) ExistsByNameExcludeId(name string, excludeId uint) (bool, error) {
	var count int64
	err := global.DB.Model(&model.Personnel{}).
		Where("name = ? AND is_deleted = ? AND id != ?", name, false, excludeId).
		Count(&count).Error
	return count > 0, err
}

// List 分页列表
func (r *PersonnelRepo) List(req dto.PersonnelListReq) ([]model.Personnel, int64, error) {
	var personnels []model.Personnel
	var total int64

	db := global.DB.Model(&model.Personnel{}).Where("is_deleted = ?", false)

	if req.Keyword != "" {
		db = db.Where("name LIKE ? OR phone LIKE ? OR company LIKE ? OR position LIKE ? OR location LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	if req.OnProject != "" {
		db = db.Where("on_project_status = ?", req.OnProject)
	}
	if req.Position != "" {
		db = db.Where("position = ?", req.Position)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := db.Offset(offset).Limit(req.PageSize).Order("sort ASC, created_at DESC").Find(&personnels).Error; err != nil {
		return nil, 0, err
	}

	return personnels, total, nil
}

// ListAll 获取所有人员（用于下拉选择）
func (r *PersonnelRepo) ListAll() ([]dto.PersonnelResp, error) {
	var personnels []model.Personnel
	err := global.DB.Where("is_deleted = ? AND status = ?", false, "active").
		Order("sort ASC, name ASC").
		Find(&personnels).Error
	if err != nil {
		return nil, err
	}

	items := make([]dto.PersonnelResp, len(personnels))
	for i, p := range personnels {
		items[i] = dto.PersonnelResp{
			ID:               p.ID,
			Name:             p.Name,
			Phone:            p.Phone,
			Email:            p.Email,
			Company:          p.Company,
			Position:         p.Position,
			WorkExperience:   p.WorkExperience,
			EntryDate:        p.EntryDate,
			ProjectStartDate: p.ProjectStartDate,
			OnProjectStatus: p.OnProjectStatus,
			Salary:           p.Salary,
			Location:         p.Location,
			Remark:           p.Remark,
			Status:           p.Status,
			Sort:             p.Sort,
			CreatedAt:        p.CreatedAt,
			UpdatedAt:        p.UpdatedAt,
		}
	}
	return items, nil
}

// ListForExport 获取所有人员（不分页，用于导出）
func (r *PersonnelRepo) ListForExport(req dto.PersonnelListReq) ([]model.Personnel, error) {
	var personnels []model.Personnel
	db := global.DB.Model(&model.Personnel{}).Where("is_deleted = ?", false)

	if req.Keyword != "" {
		db = db.Where("name LIKE ? OR phone LIKE ? OR company LIKE ? OR position LIKE ? OR location LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	if req.OnProject != "" {
		db = db.Where("on_project_status = ?", req.OnProject)
	}

	err := db.Order("sort ASC, created_at DESC").Find(&personnels).Error
	return personnels, err
}

// Statistics 职位统计（后端计算，避免前端加载全量数据）
func (r *PersonnelRepo) Statistics(req dto.PersonnelListReq) (*dto.PersonnelStatisticsResp, error) {
	db := global.DB.Model(&model.Personnel{}).Where("is_deleted = ?", false)

	if req.Keyword != "" {
		db = db.Where("name LIKE ? OR phone LIKE ? OR company LIKE ? OR position LIKE ? OR location LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	if req.OnProject != "" {
		db = db.Where("on_project_status = ?", req.OnProject)
	}
	if req.Position != "" {
		db = db.Where("position = ?", req.Position)
	}

	// 总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 如果指定了职位，只统计该职位
	if req.Position != "" {
		return &dto.PersonnelStatisticsResp{
			Total: total,
			ByPosition: []dto.PositionCount{{Position: req.Position, Count: total}},
		}, nil
	}

	// 按职位统计
	type positionRow struct {
		Position string
		Count    int64
	}
	var rows []positionRow
	err := db.Select("COALESCE(NULLIF(TRIM(position), ''), '(空职位)') as position, COUNT(*) as count").
		Group("COALESCE(NULLIF(TRIM(position), ''), '(空职位)')").
		Order("count DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	byPosition := make([]dto.PositionCount, len(rows))
	for i, row := range rows {
		byPosition[i] = dto.PositionCount{Position: row.Position, Count: row.Count}
	}

	return &dto.PersonnelStatisticsResp{Total: total, ByPosition: byPosition}, nil
}
