package service

import (
	"datauptwo/app/dto"
	"datauptwo/app/model"
	"datauptwo/app/repo"
	"encoding/json"
	"errors"
)

type RoleService struct {
	roleRepo *repo.RoleRepo
}

func NewRoleService() *RoleService {
	return &RoleService{
		roleRepo: repo.NewRoleRepo(),
	}
}

// Create 创建角色
func (s *RoleService) Create(req dto.RoleCreateReq) (*dto.RoleResp, error) {
	// 检查编码是否存在
	count, _ := s.roleRepo.CountByCode(req.Code)
	if count > 0 {
		return nil, errors.New("角色编码已存在")
	}

	// 处理权限
	permissionsJSON := ""
	if len(req.Permissions) > 0 {
		b, _ := json.Marshal(req.Permissions)
		permissionsJSON = string(b)
	}

	role := &model.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Permissions: permissionsJSON,
		Sort:        req.Sort,
	}

	if err := s.roleRepo.Create(role); err != nil {
		return nil, err
	}

	return s.toRoleResp(role), nil
}

// Update 更新角色
func (s *RoleService) Update(req dto.RoleUpdateReq) (*dto.RoleResp, error) {
	role, err := s.roleRepo.GetByID(req.ID)
	if err != nil {
		return nil, errors.New("角色不存在")
	}

	// 检查编码是否被其他记录使用
	if req.Code != role.Code {
		count, _ := s.roleRepo.CountByCode(req.Code)
		if count > 0 {
			return nil, errors.New("角色编码已存在")
		}
	}

	// 处理权限
	permissionsJSON := ""
	if len(req.Permissions) > 0 {
		b, _ := json.Marshal(req.Permissions)
		permissionsJSON = string(b)
	}

	role.Name = req.Name
	role.Code = req.Code
	role.Description = req.Description
	role.Permissions = permissionsJSON
	role.Sort = req.Sort

	if err := s.roleRepo.Update(role); err != nil {
		return nil, err
	}

	return s.toRoleResp(role), nil
}

// Delete 删除角色
func (s *RoleService) Delete(id uint) error {
	_, err := s.roleRepo.GetByID(id)
	if err != nil {
		return errors.New("角色不存在")
	}
	return s.roleRepo.Delete(id)
}

// BatchDelete 批量删除角色
func (s *RoleService) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return s.roleRepo.BatchDelete(ids)
}

// GetByID 获取角色详情
func (s *RoleService) GetByID(id uint) (*dto.RoleResp, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("角色不存在")
	}
	return s.toRoleResp(role), nil
}

// List 分页列表
func (s *RoleService) List(req dto.RoleListReq) (*dto.PageResult, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}

	roles, total, err := s.roleRepo.List(req)
	if err != nil {
		return nil, err
	}

	items := make([]dto.RoleResp, len(roles))
	for i, role := range roles {
		items[i] = *s.toRoleResp(&role)
	}

	return &dto.PageResult{Total: total, Items: items}, nil
}

// GetAll 获取所有角色
func (s *RoleService) GetAll() ([]dto.RoleResp, error) {
	roles, err := s.roleRepo.GetAll()
	if err != nil {
		return nil, err
	}

	items := make([]dto.RoleResp, len(roles))
	for i, role := range roles {
		items[i] = *s.toRoleResp(&role)
	}

	return items, nil
}

// toRoleResp 转换为角色响应
func (s *RoleService) toRoleResp(role *model.Role) *dto.RoleResp {
	resp := &dto.RoleResp{
		ID:          role.ID,
		Name:        role.Name,
		Code:        role.Code,
		Description: role.Description,
		Sort:        role.Sort,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}

	// 解析权限
	if role.Permissions != "" {
		var permissions []string
		if err := json.Unmarshal([]byte(role.Permissions), &permissions); err == nil {
			resp.Permissions = permissions
		}
	}

	return resp
}
