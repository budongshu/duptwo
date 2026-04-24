package repo

import (
	"datauptwo/app/dto"
	"datauptwo/app/model"
	"datauptwo/global"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (r *UserRepo) Create(user *model.User) error {
	return global.DB.Create(user).Error
}

func (r *UserRepo) Update(user *model.User) error {
	return global.DB.Save(user).Error
}

func (r *UserRepo) Delete(id uint) error {
	return global.DB.Model(&model.User{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *UserRepo) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return global.DB.Model(&model.User{}).Where("id IN ?", ids).Update("is_deleted", true).Error
}

func (r *UserRepo) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := global.DB.Where("id = ? AND is_deleted = ?", id, false).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := global.DB.Where("username = ? AND is_deleted = ?", username, false).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) List(req dto.UserListReq) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	db := global.DB.Model(&model.User{}).Where("is_deleted = ?", false)

	if req.Keyword != "" {
		db = db.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	if req.RoleID > 0 {
		db = db.Where("role_id = ?", req.RoleID)
	}

	if req.GroupID > 0 {
		db = db.Where("group_id = ?", req.GroupID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := db.Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepo) CountByUsername(username string) (int64, error) {
	var count int64
	err := global.DB.Model(&model.User{}).Where("username = ? AND is_deleted = ?", username, false).Count(&count).Error
	return count, err
}

func (r *UserRepo) UpdateLastLogin(id uint, ip string) error {
	return global.DB.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_login_at": global.DB.NowFunc(),
		"last_login_ip": ip,
	}).Error
}

// ListForExport 获取所有用户（用于导出）
func (r *UserRepo) ListForExport(req dto.UserListReq) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	db := global.DB.Model(&model.User{}).Where("is_deleted = ?", false)

	if req.Keyword != "" {
		db = db.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
