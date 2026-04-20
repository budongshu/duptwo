package sync

import (
	"datauptwo/app/dto"
	"datauptwo/app/service"
	"datauptwo/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SyncApi 同步API
type SyncApi struct {
	syncService  *service.SyncService
	auditService *service.AuditService
}

func NewSyncApi() *SyncApi {
	return &SyncApi{
		syncService:  service.NewSyncService(),
		auditService: service.NewAuditService(),
	}
}

// ============ 站点管理 ============

// CreateStation 创建站点
// @Summary 创建站点
// @Tags Sync
// @Security Bearer
// @Accept json
// @Param request body dto.SyncStationCreateReq true "站点信息"
// @Success 200 {object} dto.Response
// @Router /api/sync/stations [post]
func (api *SyncApi) CreateStation(c *gin.Context) {
	var req dto.SyncStationCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	station, err := api.syncService.CreateStation(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	// 记录操作日志
	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"同步站点管理", "create", "SyncStation", station.ID,
		station.Name, c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "创建成功", Data: station})
}

// ListStations 获取站点列表
// @Summary 获取站点列表
// @Tags Sync
// @Security Bearer
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param keyword query string false "搜索关键词"
// @Param status query string false "状态"
// @Success 200 {object} dto.Response
// @Router /api/sync/stations [get]
func (api *SyncApi) ListStations(c *gin.Context) {
	var req dto.SyncStationListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PageSize = 20
	}

	result, err := api.syncService.ListStations(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: result})
}

// GetAllStations 获取所有站点
// @Summary 获取所有站点
// @Tags Sync
// @Security Bearer
// @Success 200 {object} dto.Response
// @Router /api/sync/stations/all [get]
func (api *SyncApi) GetAllStations(c *gin.Context) {
	stations, err := api.syncService.GetAllStations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: stations})
}

// GetStation 获取站点详情
// @Summary 获取站点详情
// @Tags Sync
// @Security Bearer
// @Param id path int true "站点ID"
// @Success 200 {object} dto.Response
// @Router /api/sync/stations/{id} [get]
func (api *SyncApi) GetStation(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	station, err := api.syncService.GetStation(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{Code: 404, Message: "站点不存在"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: station})
}

// UpdateStation 更新站点
// @Summary 更新站点
// @Tags Sync
// @Security Bearer
// @Accept json
// @Param request body dto.SyncStationUpdateReq true "站点信息"
// @Success 200 {object} dto.Response
// @Router /api/sync/stations [put]
func (api *SyncApi) UpdateStation(c *gin.Context) {
	var req dto.SyncStationUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	station, err := api.syncService.UpdateStation(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	// 记录操作日志
	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"同步站点管理", "update", "SyncStation", req.ID,
		station.Name, c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "更新成功", Data: station})
}

// DeleteStation 删除站点
// @Summary 删除站点
// @Tags Sync
// @Security Bearer
// @Param id path int true "站点ID"
// @Success 200 {object} dto.Response
// @Router /api/sync/stations/{id} [delete]
func (api *SyncApi) DeleteStation(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// 先获取站点信息用于日志
	station, _ := api.syncService.GetStation(uint(id))
	stationName := ""
	if station != nil {
		stationName = station.Name
	}

	if err := api.syncService.DeleteStation(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	// 记录操作日志
	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"同步站点管理", "delete", "SyncStation", uint(id),
		stationName, c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "删除成功"})
}

// ============ 站点注册 ============

// Register 站点注册
// @Summary 站点注册
// @Tags Sync
// @Accept json
// @Param request body dto.SyncRegisterReq true "注册信息"
// @Success 200 {object} dto.Response
// @Router /api/sync/register [post]
func (api *SyncApi) Register(c *gin.Context) {
	var req dto.SyncRegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	resp, err := api.syncService.RegisterStation(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "注册成功", Data: resp})
}

// ============ 记录同步 ============

// UploadRecords 上传记录同步
// @Summary 上传记录同步
// @Tags Sync
// @Accept json
// @Param X-API-Key header string true "API Key"
// @Param request body dto.SyncUploadReq true "同步请求"
// @Success 200 {object} dto.Response
// @Router /api/sync/upload-records [post]
func (api *SyncApi) UploadRecords(c *gin.Context) {
	var req dto.SyncUploadReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	if len(req.Records) == 0 {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "记录列表不能为空"})
		return
	}

	// 获取站点ID
	stationID := middleware.GetStationID(c)

	resp, err := api.syncService.UploadRecords(req, stationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: resp})
}

// ============ 同步历史 ============

// GetHistory 获取同步历史
// @Summary 获取同步历史
// @Tags Sync
// @Security Bearer
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param stationId query int false "站点ID"
// @Param direction query string false "方向"
// @Param status query string false "状态"
// @Param startDate query string false "开始日期"
// @Param endDate query string false "结束日期"
// @Success 200 {object} dto.Response
// @Router /api/sync/history [get]
func (api *SyncApi) GetHistory(c *gin.Context) {
	var req dto.SyncHistoryListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PageSize = 20
	}

	result, err := api.syncService.GetHistory(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: result})
}

// GetHistoryDetails 获取同步历史详情
// @Summary 获取同步历史详情
// @Tags Sync
// @Security Bearer
// @Param id path int true "历史ID"
// @Success 200 {object} dto.Response
// @Router /api/sync/history/{id} [get]
func (api *SyncApi) GetHistoryDetails(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	result, err := api.syncService.GetHistoryDetails(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{Code: 404, Message: "同步历史不存在"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: result})
}

// ============ 同步状态 ============

// GetStatus 获取同步状态
// @Summary 获取同步状态
// @Tags Sync
// @Security Bearer
// @Success 200 {object} dto.Response
// @Router /api/sync/status [get]
func (api *SyncApi) GetStatus(c *gin.Context) {
	status := api.syncService.GetSyncStatus()
	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: status})
}

// ============ 辅助方法 ============

// getUserID 获取当前用户ID
func (api *SyncApi) getUserID(c *gin.Context) uint {
	if id, exists := c.Get("userId"); exists {
		return id.(uint)
	}
	return 0
}

// getUsername 获取当前用户名
func (api *SyncApi) getUsername(c *gin.Context) string {
	if username, exists := c.Get("username"); exists {
		return username.(string)
	}
	return ""
}
