package api

import (
	"fmt"
	"net/http"
	"strconv"

	"datauptwo/app/dto"
	"datauptwo/app/service"

	"github.com/gin-gonic/gin"
)

type FieldConfigApi struct {
	fieldConfigService *service.FieldConfigService
	auditService       *service.AuditService
}

func NewFieldConfigApi() *FieldConfigApi {
	return &FieldConfigApi{
		fieldConfigService: service.NewFieldConfigService(),
		auditService:       service.NewAuditService(),
	}
}

func (api *FieldConfigApi) getUserID(c *gin.Context) uint {
	if id, exists := c.Get("userId"); exists {
		return id.(uint)
	}
	return 0
}

func (api *FieldConfigApi) getUsername(c *gin.Context) string {
	if username, exists := c.Get("username"); exists {
		return username.(string)
	}
	return ""
}

// Create 创建字段配置
func (api *FieldConfigApi) Create(c *gin.Context) {
	var req dto.FieldConfigCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	config, err := api.fieldConfigService.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"字段配置", "create", "FieldConfig", config.ID,
		config.Name, c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "创建成功", Data: config})
}

// List 字段配置列表
func (api *FieldConfigApi) List(c *gin.Context) {
	var req dto.FieldConfigListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PageSize = 20
	}

	result, err := api.fieldConfigService.List(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: result})
}

// GetByID 获取字段配置详情
func (api *FieldConfigApi) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	resp, err := api.fieldConfigService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{Code: 404, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: resp})
}

// Update 更新字段配置
func (api *FieldConfigApi) Update(c *gin.Context) {
	var req dto.FieldConfigUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	// 先获取字段配置信息用于日志
	oldConfig, _ := api.fieldConfigService.GetByID(req.ID)
	configName := ""
	if oldConfig != nil {
		configName = oldConfig.Name
	}

	updatedConfig, err := api.fieldConfigService.Update(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"字段配置", "update", "FieldConfig", req.ID,
		configName, c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "更新成功", Data: updatedConfig})
}

// Delete 删除字段配置
func (api *FieldConfigApi) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := api.fieldConfigService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"字段配置", "delete", "FieldConfig", uint(id),
		"", c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "删除成功"})
}

// BatchDelete 批量删除字段配置
func (api *FieldConfigApi) BatchDelete(c *gin.Context) {
	var req dto.BatchDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "请选择要删除的字段"})
		return
	}

	if err := api.fieldConfigService.BatchDelete(req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"字段配置", "batch_delete", "FieldConfig", 0,
		fmt.Sprintf("批量删除 %d 个字段", len(req.IDs)), c.ClientIP(), c.GetHeader("User-Agent"),
		map[string]interface{}{"ids": req.IDs},
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: fmt.Sprintf("成功删除 %d 个字段", len(req.IDs))})
}

// GetAllEnabled 获取所有启用的字段配置
func (api *FieldConfigApi) GetAllEnabled(c *gin.Context) {
	configs, err := api.fieldConfigService.GetAllEnabled()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: configs})
}
