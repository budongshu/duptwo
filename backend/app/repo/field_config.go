package repo

import (
	"datauptwo/app/dto"
	"datauptwo/app/model"
	"datauptwo/global"
)

type FieldConfigRepo struct{}

func NewFieldConfigRepo() *FieldConfigRepo {
	return &FieldConfigRepo{}
}

func (r *FieldConfigRepo) Create(config *model.FieldConfig) error {
	return global.DB.Create(config).Error
}

func (r *FieldConfigRepo) Update(config *model.FieldConfig) error {
	return global.DB.Save(config).Error
}

func (r *FieldConfigRepo) Delete(id uint) error {
	return global.DB.Model(&model.FieldConfig{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *FieldConfigRepo) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return global.DB.Model(&model.FieldConfig{}).Where("id IN ?", ids).Update("is_deleted", true).Error
}

func (r *FieldConfigRepo) GetByID(id uint) (*model.FieldConfig, error) {
	var config model.FieldConfig
	err := global.DB.Where("id = ? AND is_deleted = ?", id, false).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *FieldConfigRepo) GetByCode(code string) (*model.FieldConfig, error) {
	var config model.FieldConfig
	err := global.DB.Where("code = ? AND is_deleted = ?", code, false).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *FieldConfigRepo) List(req dto.FieldConfigListReq) ([]model.FieldConfig, int64, error) {
	var configs []model.FieldConfig
	var total int64

	db := global.DB.Model(&model.FieldConfig{}).Where("is_deleted = ?", false)

	if req.Keyword != "" {
		db = db.Where("name LIKE ? OR code LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if req.Enabled != nil {
		db = db.Where("enabled = ?", *req.Enabled)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := db.Offset(offset).Limit(req.PageSize).Order("sort ASC, id ASC").Find(&configs).Error; err != nil {
		return nil, 0, err
	}

	return configs, total, nil
}

func (r *FieldConfigRepo) GetAllEnabled() ([]model.FieldConfig, error) {
	var configs []model.FieldConfig
	err := global.DB.Where("is_deleted = ? AND enabled = ?", false, true).Order("sort ASC, id ASC").Find(&configs).Error
	return configs, err
}

func (r *FieldConfigRepo) CountByCode(code string) (int64, error) {
	var count int64
	err := global.DB.Model(&model.FieldConfig{}).Where("code = ? AND is_deleted = ?", code, false).Count(&count).Error
	return count, err
}
