package service

import (
	"datauptwo/app/dto"
	"datauptwo/app/model"
	"datauptwo/app/repo"
	"errors"
)

type UserGroupService struct {
	groupRepo *repo.UserGroupRepo
	roleRepo  *repo.RoleRepo
}

func NewUserGroupService() *UserGroupService {
	return &UserGroupService{
		groupRepo: repo.NewUserGroupRepo(),
		roleRepo:  repo.NewRoleRepo(),
	}
}

// Create 创建用户组
func (s *UserGroupService) Create(req dto.UserGroupCreateReq) (*dto.UserGroupResp, error) {
	// 检查编码是否存在
	count, _ := s.groupRepo.CountByCode(req.Code)
	if count > 0 {
		return nil, errors.New("用户组编码已存在")
	}

	group := &model.UserGroup{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		RoleID:      req.RoleID,
		Sort:        req.Sort,
	}

	if err := s.groupRepo.Create(group); err != nil {
		return nil, err
	}

	return s.toGroupResp(group), nil
}

// Update 更新用户组
func (s *UserGroupService) Update(req dto.UserGroupUpdateReq) (*dto.UserGroupResp, error) {
	group, err := s.groupRepo.GetByID(req.ID)
	if err != nil {
		return nil, errors.New("用户组不存在")
	}

	// 检查编码是否被其他记录使用
	if req.Code != group.Code {
		count, _ := s.groupRepo.CountByCode(req.Code)
		if count > 0 {
			return nil, errors.New("用户组编码已存在")
		}
	}

	group.Name = req.Name
	group.Code = req.Code
	group.Description = req.Description
	group.RoleID = req.RoleID
	group.Sort = req.Sort

	if err := s.groupRepo.Update(group); err != nil {
		return nil, err
	}

	return s.toGroupResp(group), nil
}

// Delete 删除用户组
func (s *UserGroupService) Delete(id uint) error {
	_, err := s.groupRepo.GetByID(id)
	if err != nil {
		return errors.New("用户组不存在")
	}
	return s.groupRepo.Delete(id)
}

// BatchDelete 批量删除用户组
func (s *UserGroupService) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return s.groupRepo.BatchDelete(ids)
}

// GetByID 获取用户组详情
func (s *UserGroupService) GetByID(id uint) (*dto.UserGroupResp, error) {
	group, err := s.groupRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("用户组不存在")
	}
	return s.toGroupResp(group), nil
}

// List 分页列表
func (s *UserGroupService) List(req dto.UserGroupListReq) (*dto.PageResult, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}

	groups, total, err := s.groupRepo.List(req)
	if err != nil {
		return nil, err
	}

	items := make([]dto.UserGroupResp, len(groups))
	for i, group := range groups {
		items[i] = *s.toGroupResp(&group)
	}

	return &dto.PageResult{Total: total, Items: items}, nil
}

// GetAll 获取所有用户组
func (s *UserGroupService) GetAll() ([]dto.UserGroupResp, error) {
	groups, err := s.groupRepo.GetAll()
	if err != nil {
		return nil, err
	}

	items := make([]dto.UserGroupResp, len(groups))
	for i, group := range groups {
		items[i] = *s.toGroupResp(&group)
	}

	return items, nil
}

// toGroupResp 转换为用户组响应
func (s *UserGroupService) toGroupResp(group *model.UserGroup) *dto.UserGroupResp {
	resp := &dto.UserGroupResp{
		ID:          group.ID,
		Name:        group.Name,
		Code:        group.Code,
		Description: group.Description,
		RoleID:      group.RoleID,
		Sort:        group.Sort,
		CreatedAt:   group.CreatedAt,
		UpdatedAt:   group.UpdatedAt,
	}
	if group.RoleID > 0 {
		if role, err := s.roleRepo.GetByID(group.RoleID); err == nil {
			resp.RoleName = role.Name
		}
	}
	return resp
}
