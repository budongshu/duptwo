package service

import (
	"datauptwo/app/dto"
	"datauptwo/app/model"
	"datauptwo/app/repo"
	"errors"
)

type PersonnelService struct {
	personnelRepo *repo.PersonnelRepo
}

func NewPersonnelService() *PersonnelService {
	return &PersonnelService{
		personnelRepo: repo.NewPersonnelRepo(),
	}
}

// Create 创建人员
func (s *PersonnelService) Create(req dto.PersonnelCreateReq) (*model.Personnel, error) {
	// 检查姓名唯一性
	exists, _ := s.personnelRepo.ExistsByName(req.Name)
	if exists {
		return nil, errors.New("人员姓名已存在")
	}

	if req.Status == "" {
		req.Status = "active"
	}
	if req.OnProjectStatus == "" {
		req.OnProjectStatus = "离项"
	}

	personnel := &model.Personnel{
		Name:             req.Name,
		Phone:            req.Phone,
		Email:            req.Email,
		Company:          req.Company,
		Position:         req.Position,
		WorkExperience:   req.WorkExperience,
		EntryDate:        req.EntryDate,
		ProjectStartDate: req.ProjectStartDate,
		OnProjectStatus:  req.OnProjectStatus,
		Salary:           req.Salary,
		Location:         req.Location,
		Remark:           req.Remark,
		Status:           req.Status,
		Sort:             req.Sort,
	}

	if err := s.personnelRepo.Create(personnel); err != nil {
		return nil, err
	}

	return personnel, nil
}

// Update 更新人员
func (s *PersonnelService) Update(req dto.PersonnelUpdateReq) error {
	personnel, err := s.personnelRepo.GetByID(req.ID)
	if err != nil {
		return errors.New("人员不存在")
	}

	// 检查姓名唯一性（排除自身）
	if personnel.Name != req.Name {
		exists, _ := s.personnelRepo.ExistsByNameExcludeId(req.Name, req.ID)
		if exists {
			return errors.New("人员姓名已存在")
		}
	}

	personnel.Name = req.Name
	personnel.Phone = req.Phone
	personnel.Email = req.Email
	personnel.Company = req.Company
	personnel.Position = req.Position
	personnel.WorkExperience = req.WorkExperience
	personnel.EntryDate = req.EntryDate
	personnel.ProjectStartDate = req.ProjectStartDate
	personnel.OnProjectStatus = req.OnProjectStatus
	personnel.Salary = req.Salary
	personnel.Location = req.Location
	personnel.Remark = req.Remark
	personnel.Status = req.Status
	personnel.Sort = req.Sort

	return s.personnelRepo.Update(personnel)
}

// Delete 删除人员
func (s *PersonnelService) Delete(id uint) error {
	return s.personnelRepo.Delete(id)
}

// BatchDelete 批量删除
func (s *PersonnelService) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return s.personnelRepo.BatchDelete(ids)
}

// GetByID 获取详情
func (s *PersonnelService) GetByID(id uint) (*dto.PersonnelResp, error) {
	personnel, err := s.personnelRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("人员不存在")
	}

	return &dto.PersonnelResp{
		ID:               personnel.ID,
		Name:             personnel.Name,
		Phone:            personnel.Phone,
		Email:            personnel.Email,
		Company:          personnel.Company,
		Position:         personnel.Position,
		WorkExperience:   personnel.WorkExperience,
		EntryDate:        personnel.EntryDate,
		ProjectStartDate: personnel.ProjectStartDate,
		OnProjectStatus:  personnel.OnProjectStatus,
		Salary:           personnel.Salary,
		Location:         personnel.Location,
		Remark:           personnel.Remark,
		Status:           personnel.Status,
		Sort:             personnel.Sort,
		CreatedAt:        personnel.CreatedAt,
		UpdatedAt:        personnel.UpdatedAt,
	}, nil
}

// List 分页列表
func (s *PersonnelService) List(req dto.PersonnelListReq) (*dto.PageResult, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}

	personnels, total, err := s.personnelRepo.List(req)
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
			OnProjectStatus:  p.OnProjectStatus,
			Salary:           p.Salary,
			Location:         p.Location,
			Remark:           p.Remark,
			Status:           p.Status,
			Sort:             p.Sort,
			CreatedAt:        p.CreatedAt,
			UpdatedAt:        p.UpdatedAt,
		}
	}

	return &dto.PageResult{Total: total, Items: items}, nil
}

// ListForExport 获取所有人员（用于导出）
func (s *PersonnelService) ListForExport(req dto.PersonnelListReq) ([]dto.PersonnelResp, error) {
	personnels, err := s.personnelRepo.ListForExport(req)
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
			OnProjectStatus:  p.OnProjectStatus,
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

// ListAll 获取所有人员（用于下拉选择）
func (s *PersonnelService) ListAll() ([]dto.PersonnelResp, error) {
	return s.personnelRepo.ListAll()
}
