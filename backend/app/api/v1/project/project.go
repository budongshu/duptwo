package api

import (
	"net/http"
	"strconv"

	"datauptwo/app/dto"
	"datauptwo/app/service"

	"github.com/gin-gonic/gin"
)

type ProjectApi struct {
	projectService *service.ProjectService
}

func NewProjectApi() *ProjectApi {
	return &ProjectApi{
		projectService: service.NewProjectService(),
	}
}

// Create 创建项目
// @Summary 创建项目
// @Tags Project
// @Accept json
// @Param request body dto.ProjectCreateReq true "项目信息"
// @Success 200 {object} dto.Response
// @Router /api/projects [post]
func (api *ProjectApi) Create(c *gin.Context) {
	var req dto.ProjectCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	project, err := api.projectService.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "创建成功", Data: project})
}

// List 项目列表
// @Summary 项目列表
// @Tags Project
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param keyword query string false "关键词"
// @Param status query string false "状态"
// @Success 200 {object} dto.Response
// @Router /api/projects [get]
func (api *ProjectApi) List(c *gin.Context) {
	var req dto.ProjectListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PageSize = 20
	}

	result, err := api.projectService.List(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: result})
}

// GetByID 获取项目详情
// @Summary 获取项目详情
// @Tags Project
// @Param id path int true "项目ID"
// @Success 200 {object} dto.Response
// @Router /api/projects/{id} [get]
func (api *ProjectApi) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	project, err := api.projectService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{Code: 404, Message: "项目不存在"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: project})
}

// Update 更新项目
// @Summary 更新项目
// @Tags Project
// @Accept json
// @Param request body dto.ProjectUpdateReq true "项目信息"
// @Success 200 {object} dto.Response
// @Router /api/projects [put]
func (api *ProjectApi) Update(c *gin.Context) {
	var req dto.ProjectUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	if err := api.projectService.Update(req); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "更新成功"})
}

// Delete 删除项目
// @Summary 删除项目
// @Tags Project
// @Param id path int true "项目ID"
// @Success 200 {object} dto.Response
// @Router /api/projects/{id} [delete]
func (api *ProjectApi) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := api.projectService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "删除成功"})
}

// BatchDelete 批量删除项目
// @Summary 批量删除项目
// @Tags Project
// @Accept json
// @Param request body dto.BatchDeleteReq true "项目ID列表"
// @Success 200 {object} dto.Response
// @Router /api/projects/batch-delete [post]
func (api *ProjectApi) BatchDelete(c *gin.Context) {
	var req dto.BatchDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "请选择要删除的项目"})
		return
	}

	if err := api.projectService.BatchDelete(req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "成功删除"})
}

// ListSimple 获取所有项目（用于下拉选择）
// @Summary 获取所有项目
// @Tags Project
// @Success 200 {object} dto.Response
// @Router /api/projects/simple [get]
func (api *ProjectApi) ListSimple(c *gin.Context) {
	projects, err := api.projectService.ListAllSimple()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: projects})
}

// ListKanban 获取所有项目（看板视图）
func (api *ProjectApi) ListKanban(c *gin.Context) {
	projects, err := api.projectService.ListAllForKanban()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: projects})
}
