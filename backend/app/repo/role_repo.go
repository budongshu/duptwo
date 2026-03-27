package repo

import (
	"datauptwo/app/dto"
	"datauptwo/app/model"
	"datauptwo/global"
	"encoding/json"
)

type RoleRepo struct{}

func NewRoleRepo() *RoleRepo {
	return &RoleRepo{}
}

func (r *RoleRepo) Create(role *model.Role) error {
	return global.DB.Create(role).Error
}

func (r *RoleRepo) Update(role *model.Role) error {
	return global.DB.Save(role).Error
}

func (r *RoleRepo) Delete(id uint) error {
	return global.DB.Model(&model.Role{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *RoleRepo) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return global.DB.Model(&model.Role{}).Where("id IN ?", ids).Update("is_deleted", true).Error
}

func (r *RoleRepo) GetByID(id uint) (*model.Role, error) {
	var role model.Role
	err := global.DB.Where("id = ? AND is_deleted = ?", id, false).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepo) GetByCode(code string) (*model.Role, error) {
	var role model.Role
	err := global.DB.Where("code = ? AND is_deleted = ?", code, false).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepo) List(req dto.RoleListReq) ([]model.Role, int64, error) {
	var roles []model.Role
	var total int64

	db := global.DB.Model(&model.Role{}).Where("is_deleted = ?", false)

	if req.Keyword != "" {
		db = db.Where("name LIKE ? OR code LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := db.Offset(offset).Limit(req.PageSize).Order("sort ASC, id ASC").Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (r *RoleRepo) GetAll() ([]model.Role, error) {
	var roles []model.Role
	err := global.DB.Where("is_deleted = ?", false).Order("sort ASC, id ASC").Find(&roles).Error
	return roles, err
}

func (r *RoleRepo) CountByCode(code string) (int64, error) {
	var count int64
	err := global.DB.Model(&model.Role{}).Where("code = ? AND is_deleted = ?", code, false).Count(&count).Error
	return count, err
}

// GetAllPermissions 获取所有角色拥有的权限
func (r *RoleRepo) GetAllPermissions() ([]string, error) {
	var roles []model.Role
	err := global.DB.Where("is_deleted = ?", false).Find(&roles).Error
	if err != nil {
		return nil, err
	}

	permSet := make(map[string]bool)
	for _, role := range roles {
		if role.Permissions == "" {
			continue
		}
		var perms []string
		if err := json.Unmarshal([]byte(role.Permissions), &perms); err != nil {
			// 权限数据损坏，跳过该角色
			continue
		}
		for _, p := range perms {
			permSet[p] = true
		}
	}

	var perms []string
	for k := range permSet {
		perms = append(perms, k)
	}
	return perms, nil
}
