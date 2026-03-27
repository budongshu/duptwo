package role

import (
	"datauptwo/app/dto"
	"datauptwo/app/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleApi struct {
	roleService *service.RoleService
}

func NewRoleApi() *RoleApi {
	return &RoleApi{
		roleService: service.NewRoleService(),
	}
}

// Create 创建角色
// @Summary 创建角色
// @Tags Role
// @Security Bearer
// @Accept json
// @Param request body dto.RoleCreateReq true "角色信息"
// @Success 200 {object} dto.Response
// @Router /api/roles [post]
func (api *RoleApi) Create(c *gin.Context) {
	var req dto.RoleCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	role, err := api.roleService.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "创建成功", Data: role})
}

// List 角色列表
// @Summary 角色列表
// @Tags Role
// @Security Bearer
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} dto.Response
// @Router /api/roles [get]
func (api *RoleApi) List(c *gin.Context) {
	var req dto.RoleListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PageSize = 20
	}

	result, err := api.roleService.List(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: result})
}

// GetAll 获取所有角色
// @Summary 获取所有角色
// @Tags Role
// @Security Bearer
// @Success 200 {object} dto.Response
// @Router /api/roles/all [get]
func (api *RoleApi) GetAll(c *gin.Context) {
	roles, err := api.roleService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: roles})
}

// GetByID 获取角色详情
// @Summary 获取角色详情
// @Tags Role
// @Security Bearer
// @Param id path int true "角色ID"
// @Success 200 {object} dto.Response
// @Router /api/roles/{id} [get]
func (api *RoleApi) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	role, err := api.roleService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{Code: 404, Message: "角色不存在"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: role})
}

// Update 更新角色
// @Summary 更新角色
// @Tags Role
// @Security Bearer
// @Accept json
// @Param request body dto.RoleUpdateReq true "角色信息"
// @Success 200 {object} dto.Response
// @Router /api/roles [put]
func (api *RoleApi) Update(c *gin.Context) {
	var req dto.RoleUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	role, err := api.roleService.Update(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "更新成功", Data: role})
}

// Delete 删除角色
// @Summary 删除角色
// @Tags Role
// @Security Bearer
// @Param id path int true "角色ID"
// @Success 200 {object} dto.Response
// @Router /api/roles/{id} [delete]
func (api *RoleApi) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := api.roleService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "删除成功"})
}

// BatchDelete 批量删除角色
// @Summary 批量删除角色
// @Tags Role
// @Security Bearer
// @Accept json
// @Param request body dto.BatchDeleteReq true "角色ID列表"
// @Success 200 {object} dto.Response
// @Router /api/roles/batch-delete [post]
func (api *RoleApi) BatchDelete(c *gin.Context) {
	var req dto.BatchDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "请选择要删除的角色"})
		return
	}

	if err := api.roleService.BatchDelete(req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: fmt.Sprintf("成功删除 %d 个角色", len(req.IDs))})
}
