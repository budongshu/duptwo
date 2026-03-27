package service

import (
	"encoding/json"
	"datauptwo/app/dto"
	"datauptwo/app/model"
	"datauptwo/app/repo"
	"errors"
)

type ProjectService struct {
	projectRepo *repo.ProjectRepo
}

func NewProjectService() *ProjectService {
	return &ProjectService{
		projectRepo: repo.NewProjectRepo(),
	}
}

// Create 创建项目
func (s *ProjectService) Create(req dto.ProjectCreateReq) (*model.Project, error) {
	// 检查编码唯一性
	exists, _ := s.projectRepo.ExistsByCode(req.Code)
	if exists {
		return nil, errors.New("项目编码已存在")
	}

	// 检查名称唯一性
	exists, _ = s.projectRepo.ExistsByName(req.Name)
	if exists {
		return nil, errors.New("项目名称已存在")
	}

	// 序列化驻场点信息
	onsiteJSON := ""
	if len(req.OnsiteStations) > 0 {
		data, _ := json.Marshal(req.OnsiteStations)
		onsiteJSON = string(data)
	}

	project := &model.Project{
		Name:           req.Name,
		Code:          req.Code,
		Description:   req.Description,
		Status:        req.Status,
		Stage:         req.Stage,
		Sort:          req.Sort,
		ProjectPerson: req.ProjectPerson,
		OpsPerson:     req.OpsPerson,
		CompanyAddr:   req.CompanyAddr,
		ProjectPeriod: req.ProjectPeriod,
		OnsiteStations: onsiteJSON,
	}

	if project.Status == "" {
		project.Status = "active"
	}
	if project.Stage == "" {
		project.Stage = "planning"
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, err
	}

	return project, nil
}

// Update 更新项目
func (s *ProjectService) Update(req dto.ProjectUpdateReq) error {
	project, err := s.projectRepo.GetByID(req.ID)
	if err != nil {
		return errors.New("项目不存在")
	}

	// 检查编码唯一性（排除自身）
	if project.Code != req.Code {
		exists, _ := s.projectRepo.ExistsByCode(req.Code)
		if exists {
			return errors.New("项目编码已存在")
		}
	}

	// 检查名称唯一性（排除自身）
	if project.Name != req.Name {
		exists, _ := s.projectRepo.ExistsByName(req.Name)
		if exists {
			return errors.New("项目名称已存在")
		}
	}

	// 序列化驻场点信息
	onsiteJSON := ""
	if len(req.OnsiteStations) > 0 {
		data, _ := json.Marshal(req.OnsiteStations)
		onsiteJSON = string(data)
	}

	project.Name = req.Name
	project.Code = req.Code
	project.Description = req.Description
	project.Status = req.Status
	project.Stage = req.Stage
	project.Sort = req.Sort
	project.ProjectPerson = req.ProjectPerson
	project.OpsPerson = req.OpsPerson
	project.CompanyAddr = req.CompanyAddr
	project.ProjectPeriod = req.ProjectPeriod
	project.OnsiteStations = onsiteJSON

	return s.projectRepo.Update(project)
}

// Delete 删除项目
func (s *ProjectService) Delete(id uint) error {
	return s.projectRepo.Delete(id)
}

// BatchDelete 批量删除
func (s *ProjectService) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return s.projectRepo.BatchDelete(ids)
}

// GetByID 获取详情
func (s *ProjectService) GetByID(id uint) (*dto.ProjectResp, error) {
	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("项目不存在")
	}

	resp := s.toProjectRespSimple(project)
	return resp, nil
}

// List 分页列表
func (s *ProjectService) List(req dto.ProjectListReq) (*dto.PageResult, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}

	projects, total, err := s.projectRepo.List(req)
	if err != nil {
		return nil, err
	}

	// 批量获取统计数据
	names := make([]string, len(projects))
	for i, p := range projects {
		names[i] = p.Name
	}
	statsMap, _ := s.projectRepo.BatchGetProjectStats(names)

	items := make([]dto.ProjectResp, len(projects))
	for i, p := range projects {
		items[i] = *s.toProjectResp(&p, statsMap[p.Name])
	}

	return &dto.PageResult{Total: total, Items: items}, nil
}

// ListAllSimple 获取所有项目（用于下拉选择）
func (s *ProjectService) ListAllSimple() ([]dto.ProjectSimpleResp, error) {
	return s.projectRepo.ListAllSimple()
}

// ListAllForKanban 获取所有项目（看板视图）
func (s *ProjectService) ListAllForKanban() ([]dto.ProjectResp, error) {
	projects, err := s.projectRepo.ListAllForKanban()
	if err != nil {
		return nil, err
	}

	// 批量获取统计数据
	names := make([]string, len(projects))
	for i, p := range projects {
		names[i] = p.Name
	}
	statsMap, _ := s.projectRepo.BatchGetProjectStats(names)

	items := make([]dto.ProjectResp, len(projects))
	for i, p := range projects {
		items[i] = *s.toProjectResp(&p, statsMap[p.Name])
	}
	return items, nil
}

// toProjectResp 转换为响应（使用预取的统计数据）
func (s *ProjectService) toProjectResp(p *model.Project, stats *repo.ProjectStats) *dto.ProjectResp {
	recordCount := int64(0)
	totalDataSize := int64(0)
	if stats != nil {
		recordCount = stats.RecordCount
		totalDataSize = stats.TotalSize
	}

	// 反序列化驻场点信息
	var onsiteStations []dto.OnSiteStation
	if p.OnsiteStations != "" {
		json.Unmarshal([]byte(p.OnsiteStations), &onsiteStations)
	}

	stage := p.Stage
	if stage == "" {
		stage = "planning"
	}

	return &dto.ProjectResp{
		ID:            p.ID,
		Name:          p.Name,
		Code:          p.Code,
		Description:   p.Description,
		Status:        p.Status,
		Stage:         stage,
		Sort:          p.Sort,
		ProjectPerson: p.ProjectPerson,
		OpsPerson:     p.OpsPerson,
		CompanyAddr:   p.CompanyAddr,
		ProjectPeriod: p.ProjectPeriod,
		OnsiteStations: onsiteStations,
		RecordCount:   recordCount,
		TotalDataSize: totalDataSize,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

// toProjectRespSimple 简单转换（不查统计，用于单条记录详情）
func (s *ProjectService) toProjectRespSimple(p *model.Project) *dto.ProjectResp {
	var onsiteStations []dto.OnSiteStation
	if p.OnsiteStations != "" {
		json.Unmarshal([]byte(p.OnsiteStations), &onsiteStations)
	}
	stage := p.Stage
	if stage == "" {
		stage = "planning"
	}
	return &dto.ProjectResp{
		ID:            p.ID,
		Name:          p.Name,
		Code:          p.Code,
		Description:   p.Description,
		Status:        p.Status,
		Stage:         stage,
		Sort:          p.Sort,
		ProjectPerson: p.ProjectPerson,
		OpsPerson:     p.OpsPerson,
		CompanyAddr:   p.CompanyAddr,
		ProjectPeriod: p.ProjectPeriod,
		OnsiteStations: onsiteStations,
		RecordCount:   0,
		TotalDataSize: 0,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}
