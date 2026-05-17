package repo

import (
	"datauptwo/app/dto"
	"datauptwo/app/model"
	"datauptwo/global"
)

type UserGroupRepo struct{}

func NewUserGroupRepo() *UserGroupRepo {
	return &UserGroupRepo{}
}

func (r *UserGroupRepo) Create(group *model.UserGroup) error {
	return global.DB.Create(group).Error
}

func (r *UserGroupRepo) Update(group *model.UserGroup) error {
	return global.DB.Save(group).Error
}

func (r *UserGroupRepo) Delete(id uint) error {
	return global.DB.Model(&model.UserGroup{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *UserGroupRepo) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return global.DB.Model(&model.UserGroup{}).Where("id IN ?", ids).Update("is_deleted", true).Error
}

func (r *UserGroupRepo) GetByID(id uint) (*model.UserGroup, error) {
	var group model.UserGroup
	err := global.DB.Where("id = ? AND is_deleted = ?", id, false).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *UserGroupRepo) GetByCode(code string) (*model.UserGroup, error) {
	var group model.UserGroup
	err := global.DB.Where("code = ? AND is_deleted = ?", code, false).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *UserGroupRepo) List(req dto.UserGroupListReq) ([]model.UserGroup, int64, error) {
	var groups []model.UserGroup
	var total int64

	db := global.DB.Model(&model.UserGroup{}).Where("is_deleted = ?", false)

	if req.Keyword != "" {
		db = db.Where("name LIKE ? OR code LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := db.Offset(offset).Limit(req.PageSize).Order("sort ASC, id ASC").Find(&groups).Error; err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}

func (r *UserGroupRepo) GetAll() ([]model.UserGroup, error) {
	var groups []model.UserGroup
	err := global.DB.Where("is_deleted = ?", false).Order("sort ASC, id ASC").Find(&groups).Error
	return groups, err
}

func (r *UserGroupRepo) CountByCode(code string) (int64, error) {
	var count int64
	err := global.DB.Model(&model.UserGroup{}).Where("code = ? AND is_deleted = ?", code, false).Count(&count).Error
	return count, err
}
