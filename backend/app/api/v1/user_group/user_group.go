package user_group

import (
	"datauptwo/app/dto"
	"datauptwo/app/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserGroupApi struct {
	groupService *service.UserGroupService
}

func NewUserGroupApi() *UserGroupApi {
	return &UserGroupApi{
		groupService: service.NewUserGroupService(),
	}
}

// Create 创建用户组
// @Summary 创建用户组
// @Tags UserGroup
// @Security Bearer
// @Accept json
// @Param request body dto.UserGroupCreateReq true "用户组信息"
// @Success 200 {object} dto.Response
// @Router /api/user-groups [post]
func (api *UserGroupApi) Create(c *gin.Context) {
	var req dto.UserGroupCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	group, err := api.groupService.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "创建成功", Data: group})
}

// List 用户组列表
// @Summary 用户组列表
// @Tags UserGroup
// @Security Bearer
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} dto.Response
// @Router /api/user-groups [get]
func (api *UserGroupApi) List(c *gin.Context) {
	var req dto.UserGroupListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PageSize = 20
	}

	result, err := api.groupService.List(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: result})
}

// GetAll 获取所有用户组
// @Summary 获取所有用户组
// @Tags UserGroup
// @Security Bearer
// @Success 200 {object} dto.Response
// @Router /api/user-groups/all [get]
func (api *UserGroupApi) GetAll(c *gin.Context) {
	groups, err := api.groupService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: groups})
}

// GetByID 获取用户组详情
// @Summary 获取用户组详情
// @Tags UserGroup
// @Security Bearer
// @Param id path int true "用户组ID"
// @Success 200 {object} dto.Response
// @Router /api/user-groups/{id} [get]
func (api *UserGroupApi) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	group, err := api.groupService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{Code: 404, Message: "用户组不存在"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: group})
}

// Update 更新用户组
// @Summary 更新用户组
// @Tags UserGroup
// @Security Bearer
// @Accept json
// @Param request body dto.UserGroupUpdateReq true "用户组信息"
// @Success 200 {object} dto.Response
// @Router /api/user-groups [put]
func (api *UserGroupApi) Update(c *gin.Context) {
	var req dto.UserGroupUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	group, err := api.groupService.Update(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "更新成功", Data: group})
}

// Delete 删除用户组
// @Summary 删除用户组
// @Tags UserGroup
// @Security Bearer
// @Param id path int true "用户组ID"
// @Success 200 {object} dto.Response
// @Router /api/user-groups/{id} [delete]
func (api *UserGroupApi) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := api.groupService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "删除成功"})
}

// BatchDelete 批量删除用户组
// @Summary 批量删除用户组
// @Tags UserGroup
// @Security Bearer
// @Accept json
// @Param request body dto.BatchDeleteReq true "用户组ID列表"
// @Success 200 {object} dto.Response
// @Router /api/user-groups/batch-delete [post]
func (api *UserGroupApi) BatchDelete(c *gin.Context) {
	var req dto.BatchDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "请选择要删除的用户组"})
		return
	}

	if err := api.groupService.BatchDelete(req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: fmt.Sprintf("成功删除 %d 个用户组", len(req.IDs))})
}
