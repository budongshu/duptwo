package admin

import (
	"datauptwo/app/dto"
	"datauptwo/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SecurityApi struct {
	securityService *service.SecurityService
}

func NewSecurityApi() *SecurityApi {
	return &SecurityApi{
		securityService: service.NewSecurityService(),
	}
}

// GetSettings 获取安全设置
// @Summary 获取安全设置
// @Tags Security
// @Security Bearer
// @Success 200 {object} dto.Response
// @Router /api/admin/security-settings [get]
func (api *SecurityApi) GetSettings(c *gin.Context) {
	settings, err := api.securityService.GetSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "获取安全设置失败"})
		return
	}
	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: settings})
}

// UpdateSettings 更新安全设置
// @Summary 更新安全设置
// @Tags Security
// @Security Bearer
// @Accept json
// @Param request body dto.SecuritySettingsUpdateReq true "安全设置"
// @Success 200 {object} dto.Response
// @Router /api/admin/security-settings [put]
func (api *SecurityApi) UpdateSettings(c *gin.Context) {
	var req dto.SecuritySettingsUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	if err := api.securityService.UpdateSettings(req); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "更新安全设置失败"})
		return
	}
	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "安全设置已更新"})
}

// GetLockedUsers 获取被锁定的用户列表
// @Summary 获取被锁定的用户列表
// @Tags Security
// @Security Bearer
// @Success 200 {object} dto.Response
// @Router /api/admin/security/locked-users [get]
func (api *SecurityApi) GetLockedUsers(c *gin.Context) {
	users, err := api.securityService.GetLockedUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "获取锁定用户失败"})
		return
	}
	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: users})
}

// UnlockUser 解锁用户
// @Summary 解锁用户
// @Tags Security
// @Security Bearer
// @Param username query string true "用户名"
// @Success 200 {object} dto.Response
// @Router /api/admin/security/unlock-user [post]
func (api *SecurityApi) UnlockUser(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "用户名不能为空"})
		return
	}
	if err := api.securityService.UnlockUser(username); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "解锁用户失败"})
		return
	}
	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "用户已解锁"})
}

// GetLockedIPs 获取被锁定的IP列表
// @Summary 获取被锁定的IP列表
// @Tags Security
// @Security Bearer
// @Success 200 {object} dto.Response
// @Router /api/admin/security/locked-ips [get]
func (api *SecurityApi) GetLockedIPs(c *gin.Context) {
	ips, err := api.securityService.GetLockedIPs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "获取锁定IP失败"})
		return
	}
	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: ips})
}

// UnlockIP 解锁IP
// @Summary 解锁IP
// @Tags Security
// @Security Bearer
// @Param ip query string true "IP地址"
// @Success 200 {object} dto.Response
// @Router /api/admin/security/unlock-ip [post]
func (api *SecurityApi) UnlockIP(c *gin.Context) {
	ip := c.Query("ip")
	if ip == "" {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "IP地址不能为空"})
		return
	}
	if err := api.securityService.UnlockIP(ip); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "解锁IP失败"})
		return
	}
	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "IP已解锁"})
}

// GetSecurityOverview 获取安全总览
// @Summary 获取安全总览
// @Tags Security
// @Security Bearer
// @Success 200 {object} dto.Response
// @Router /api/admin/security/overview [get]
func (api *SecurityApi) GetSecurityOverview(c *gin.Context) {
	overview, err := api.securityService.GetSecurityOverview()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "获取安全总览失败"})
		return
	}
	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: overview})
}
