package service

import (
	"datauptwo/app/dto"
	"datauptwo/app/model"
	"datauptwo/app/repo"
	"encoding/json"
	"errors"
)

type FieldConfigService struct {
	fieldConfigRepo *repo.FieldConfigRepo
}

func NewFieldConfigService() *FieldConfigService {
	return &FieldConfigService{
		fieldConfigRepo: repo.NewFieldConfigRepo(),
	}
}

// Create 创建字段配置
func (s *FieldConfigService) Create(req dto.FieldConfigCreateReq) (*model.FieldConfig, error) {
	// 检查编码是否已存在
	count, _ := s.fieldConfigRepo.CountByCode(req.Code)
	if count > 0 {
		return nil, errors.New("字段编码已存在")
	}

	// 处理选项
	optionsJSON := ""
	if len(req.Options) > 0 {
		b, _ := json.Marshal(req.Options)
		optionsJSON = string(b)
	}

	config := &model.FieldConfig{
		Name:         req.Name,
		Code:         req.Code,
		Type:         req.Type,
		Required:     req.Required,
		Options:      optionsJSON,
		DefaultValue: req.DefaultValue,
		Placeholder:  req.Placeholder,
		Sort:         req.Sort,
		Enabled:      req.Enabled,
	}

	if err := s.fieldConfigRepo.Create(config); err != nil {
		return nil, err
	}

	return config, nil
}

// Update 更新字段配置
func (s *FieldConfigService) Update(req dto.FieldConfigUpdateReq) (*model.FieldConfig, error) {
	config, err := s.fieldConfigRepo.GetByID(req.ID)
	if err != nil {
		return nil, errors.New("字段配置不存在")
	}

	// 检查编码是否被其他记录使用
	if req.Code != config.Code {
		count, _ := s.fieldConfigRepo.CountByCode(req.Code)
		if count > 0 {
			return nil, errors.New("字段编码已存在")
		}
	}

	// 处理选项
	optionsJSON := ""
	if len(req.Options) > 0 {
		b, _ := json.Marshal(req.Options)
		optionsJSON = string(b)
	}

	config.Name = req.Name
	config.Code = req.Code
	config.Type = req.Type
	config.Required = req.Required
	config.Options = optionsJSON
	config.DefaultValue = req.DefaultValue
	config.Placeholder = req.Placeholder
	config.Sort = req.Sort
	config.Enabled = req.Enabled

	if err := s.fieldConfigRepo.Update(config); err != nil {
		return nil, err
	}

	return config, nil
}

// Delete 删除字段配置
func (s *FieldConfigService) Delete(id uint) error {
	_, err := s.fieldConfigRepo.GetByID(id)
	if err != nil {
		return errors.New("字段配置不存在")
	}
	return s.fieldConfigRepo.Delete(id)
}

// BatchDelete 批量删除字段配置
func (s *FieldConfigService) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return s.fieldConfigRepo.BatchDelete(ids)
}

// GetByID 获取详情
func (s *FieldConfigService) GetByID(id uint) (*dto.FieldConfigResp, error) {
	config, err := s.fieldConfigRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("字段配置不存在")
	}
	return s.toResp(config), nil
}

// List 分页列表
func (s *FieldConfigService) List(req dto.FieldConfigListReq) (*dto.PageResult, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}

	configs, total, err := s.fieldConfigRepo.List(req)
	if err != nil {
		return nil, err
	}

	items := make([]dto.FieldConfigResp, len(configs))
	for i, config := range configs {
		items[i] = *s.toResp(&config)
	}

	return &dto.PageResult{Total: total, Items: items}, nil
}

// GetAllEnabled 获取所有启用的字段配置
func (s *FieldConfigService) GetAllEnabled() ([]dto.FieldConfigResp, error) {
	configs, err := s.fieldConfigRepo.GetAllEnabled()
	if err != nil {
		return nil, err
	}

	items := make([]dto.FieldConfigResp, len(configs))
	for i, config := range configs {
		items[i] = *s.toResp(&config)
	}

	return items, nil
}

// toResp 转换为响应格式
func (s *FieldConfigService) toResp(config *model.FieldConfig) *dto.FieldConfigResp {
	resp := &dto.FieldConfigResp{
		ID:          config.ID,
		Name:        config.Name,
		Code:        config.Code,
		Type:        config.Type,
		Required:    config.Required,
		DefaultValue: config.DefaultValue,
		Placeholder: config.Placeholder,
		Sort:        config.Sort,
		Enabled:     config.Enabled,
		CreatedAt:   config.CreatedAt,
		UpdatedAt:   config.UpdatedAt,
	}

	// 解析选项
	if config.Options != "" {
		var options []string
		if err := json.Unmarshal([]byte(config.Options), &options); err == nil {
			resp.Options = options
		}
	}

	return resp
}
