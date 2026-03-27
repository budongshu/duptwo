package user

import (
	"datauptwo/app/dto"
	"datauptwo/app/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserApi struct {
	userService  *service.UserService
	auditService *service.AuditService
}

func NewUserApi() *UserApi {
	return &UserApi{
		userService:  service.NewUserService(),
		auditService: service.NewAuditService(),
	}
}

// Create 创建用户
// @Summary 创建用户
// @Tags User
// @Security Bearer
// @Accept json
// @Param request body dto.UserCreateReq true "用户信息"
// @Success 200 {object} dto.Response
// @Router /api/users [post]
func (api *UserApi) Create(c *gin.Context) {
	var req dto.UserCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	user, err := api.userService.Create(req)
	if err != nil {
		// 根据错误类型返回合适的HTTP状态码
		errMsg := err.Error()
		if errMsg == "用户名已存在" || errMsg == "检查用户名失败" {
			c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: errMsg})
		} else {
			c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: errMsg})
		}
		return
	}

	// 记录操作日志
	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"用户管理", "create", "User", user.ID,
		user.Username, c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "创建成功", Data: user})
}

// List 用户列表
// @Summary 用户列表
// @Tags User
// @Security Bearer
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param keyword query string false "搜索关键词"
// @Param status query string false "状态"
// @Param roleId query int false "角色ID"
// @Param groupId query int false "用户组ID"
// @Success 200 {object} dto.Response
// @Router /api/users [get]
func (api *UserApi) List(c *gin.Context) {
	var req dto.UserListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PageSize = 20
	}

	result, err := api.userService.List(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	// 记录操作日志
	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"用户管理", "view", "User", 0,
		"查看用户列表", c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: result})
}

// GetAll 获取所有用户
// @Summary 获取所有用户
// @Tags User
// @Security Bearer
// @Success 200 {object} dto.Response
// @Router /api/users/all [get]
func (api *UserApi) GetAll(c *gin.Context) {
	users, err := api.userService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: users})
}

// GetByID 获取用户详情
// @Summary 获取用户详情
// @Tags User
// @Security Bearer
// @Param id path int true "用户ID"
// @Success 200 {object} dto.Response
// @Router /api/users/{id} [get]
func (api *UserApi) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	user, err := api.userService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{Code: 404, Message: "用户不存在"})
		return
	}

	// 记录操作日志
	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"用户管理", "view", "User", uint(id),
		user.Username, c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: user})
}

// Update 更新用户
// @Summary 更新用户
// @Tags User
// @Security Bearer
// @Accept json
// @Param request body dto.UserUpdateReq true "用户信息"
// @Success 200 {object} dto.Response
// @Router /api/users [put]
func (api *UserApi) Update(c *gin.Context) {
	var req dto.UserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	user, err := api.userService.Update(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	// 记录操作日志
	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"用户管理", "update", "User", req.ID,
		user.Username, c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "更新成功", Data: user})
}

// Delete 删除用户
// @Summary 删除用户
// @Tags User
// @Security Bearer
// @Param id path int true "用户ID"
// @Success 200 {object} dto.Response
// @Router /api/users/{id} [delete]
func (api *UserApi) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// 先获取用户信息用于日志
	user, _ := api.userService.GetByID(uint(id))
	username := ""
	if user != nil {
		username = user.Username
	}

	if err := api.userService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	// 记录操作日志
	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"用户管理", "delete", "User", uint(id),
		username, c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "删除成功"})
}

// BatchDelete 批量删除用户
// @Summary 批量删除用户
// @Tags User
// @Security Bearer
// @Accept json
// @Param request body dto.BatchDeleteReq true "用户ID列表"
// @Success 200 {object} dto.Response
// @Router /api/users/batch-delete [post]
func (api *UserApi) BatchDelete(c *gin.Context) {
	var req dto.BatchDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "请选择要删除的用户"})
		return
	}

	if err := api.userService.BatchDelete(req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: fmt.Sprintf("成功删除 %d 个用户", len(req.IDs))})
}
func (api *UserApi) getUserID(c *gin.Context) uint {
	if id, exists := c.Get("userId"); exists {
		return id.(uint)
	}
	return 0
}

// getUsername 获取当前用户名
func (api *UserApi) getUsername(c *gin.Context) string {
	if username, exists := c.Get("username"); exists {
		return username.(string)
	}
	return ""
}

// ResetPassword 重置用户密码
// @Summary 重置用户密码
// @Tags User
// @Security Bearer
// @Accept json
// @Param request body dto.ResetPasswordReq true "重置密码信息"
// @Success 200 {object} dto.Response
// @Router /api/users/reset-password [post]
func (api *UserApi) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	err := api.userService.ResetPassword(req.UserID, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	// 记录操作日志
	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"用户管理", "reset_password", "User", req.UserID,
		"", c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "密码重置成功"})
}

// ResetMFA 重置用户MFA
// @Summary 重置用户MFA
// @Tags User
// @Security Bearer
// @Accept json
// @Param request body dto.ResetMFAReq true "重置MFA信息"
// @Success 200 {object} dto.Response
// @Router /api/users/reset-mfa [post]
func (api *UserApi) ResetMFA(c *gin.Context) {
	var req dto.ResetMFAReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	err := api.userService.ResetMFA(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	// 记录操作日志
	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"用户管理", "reset_mfa", "User", req.UserID,
		"", c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "MFA重置成功"})
}

// GenerateMFASecret 生成MFA密钥
// @Summary 生成MFA密钥
// @Tags User
// @Security Bearer
// @Param id query int true "用户ID"
// @Success 200 {object} dto.Response
// @Router /api/users/generate-mfa-secret [get]
func (api *UserApi) GenerateMFASecret(c *gin.Context) {
	userIDStr := c.Query("id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "用户ID不能为空"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "无效的用户ID"})
		return
	}

	secret, qrCode, err := api.userService.GenerateMFASecret(uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: dto.GenerateMFASecretResp{
		Secret: secret,
		QRCode: qrCode,
	}})
}

// AdminEnableMFA 管理员为用户启用MFA
// @Summary 管理员为用户启用MFA
// @Tags User
// @Security Bearer
// @Accept json
// @Param request body dto.AdminEnableMFAReq true "启用MFA信息"
// @Success 200 {object} dto.Response
// @Router /api/users/admin-enable-mfa [post]
func (api *UserApi) AdminEnableMFA(c *gin.Context) {
	var req dto.AdminEnableMFAReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	err := api.userService.AdminEnableMFA(req.UserID, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	// 记录操作日志
	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"用户管理", "enable_mfa", "User", req.UserID,
		"", c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "MFA启用成功"})
}
