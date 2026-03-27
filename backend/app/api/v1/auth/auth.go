package auth

import (
	"datauptwo/app/dto"
	"datauptwo/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthApi struct {
	authService *service.AuthService
}

func NewAuthApi() *AuthApi {
	return &AuthApi{
		authService: service.NewAuthService(),
	}
}

// GetCaptcha 获取登录验证码
// @Summary 获取登录验证码
// @Tags Auth
// @Success 200 {object} dto.Response
// @Router /api/auth/captcha [get]
func (api *AuthApi) GetCaptcha(c *gin.Context) {
	captcha, err := api.authService.GenerateCaptcha()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "生成验证码失败"})
		return
	}
	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: captcha})
}

// Login 用户登录
// @Summary 用户登录
// @Tags Auth
// @Accept json
// @Param request body dto.LoginReq true "登录信息"
// @Success 200 {object} dto.Response
// @Router /api/auth/login [post]
func (api *AuthApi) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	result, err := api.authService.Login(req, ip, userAgent)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Response{Code: 401, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "登录成功", Data: result})
}

// Register 用户注册
// @Summary 用户注册
// @Tags Auth
// @Accept json
// @Param request body dto.RegisterReq true "注册信息"
// @Success 200 {object} dto.Response
// @Router /api/auth/register [post]
func (api *AuthApi) Register(c *gin.Context) {
	var req dto.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	result, err := api.authService.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "注册成功", Data: result})
}

// ADLogin AD用户登录
// @Summary AD用户登录
// @Tags Auth
// @Accept json
// @Param request body dto.ADLoginReq true "AD登录信息"
// @Success 200 {object} dto.Response
// @Router /api/auth/ad-login [post]
func (api *AuthApi) ADLogin(c *gin.Context) {
	var req dto.ADLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	result, err := api.authService.ADLogin(req, ip, userAgent)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Response{Code: 401, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "登录成功", Data: result})
}

// GetCurrentUser 获取当前用户信息
// @Summary 获取当前用户信息
// @Tags Auth
// @Security Bearer
// @Success 200 {object} dto.Response
// @Router /api/auth/current [get]
func (api *AuthApi) GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("userId")

	user, err := api.authService.GetCurrentUser(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{Code: 404, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: user})
}

// UpdateProfile 更新个人资料
// @Summary 更新个人资料
// @Tags Auth
// @Security Bearer
// @Accept json
// @Param request body dto.ProfileUpdateReq true "个人资料"
// @Success 200 {object} dto.Response
// @Router /api/auth/profile [put]
func (api *AuthApi) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req dto.ProfileUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	user, err := api.authService.UpdateProfile(userID.(uint), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "更新成功", Data: user})
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Tags Auth
// @Security Bearer
// @Accept json
// @Param request body dto.ChangePasswordReq true "密码信息"
// @Success 200 {object} dto.Response
// @Router /api/auth/change-password [post]
func (api *AuthApi) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req dto.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	if err := api.authService.ChangePassword(userID.(uint), req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "密码修改成功"})
}

// MFAVerify MFA验证
// @Summary MFA验证
// @Tags Auth
// @Accept json
// @Param request body dto.MFAVerifyReq true "MFA验证信息"
// @Success 200 {object} dto.Response
// @Router /api/auth/mfa/verify [post]
func (api *AuthApi) MFAVerify(c *gin.Context) {
	var req dto.MFAVerifyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	result, err := api.authService.MFAVerify(req, ip, userAgent)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Response{Code: 401, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "验证成功", Data: result})
}

// GetMFAStatus 获取MFA状态
// @Summary 获取MFA状态
// @Tags Auth
// @Security Bearer
// @Success 200 {object} dto.Response
// @Router /api/auth/mfa/status [get]
func (api *AuthApi) GetMFAStatus(c *gin.Context) {
	userID, _ := c.Get("userId")

	result, err := api.authService.GetMFAStatus(userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: result})
}

// EnableMFA 启用MFA
// @Summary 启用MFA
// @Tags Auth
// @Security Bearer
// @Accept json
// @Param request body dto.MFAEnableReq true "MFA启用信息"
// @Success 200 {object} dto.Response
// @Router /api/auth/mfa/enable [post]
func (api *AuthApi) EnableMFA(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req dto.MFAEnableReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	result, err := api.authService.EnableMFA(userID.(uint), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "MFA启用成功", Data: result})
}

// DisableMFA 禁用MFA
// @Summary 禁用MFA
// @Tags Auth
// @Security Bearer
// @Success 200 {object} dto.Response
// @Router /api/auth/mfa/disable [post]
func (api *AuthApi) DisableMFA(c *gin.Context) {
	userID, _ := c.Get("userId")

	if err := api.authService.DisableMFA(userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "MFA已禁用"})
}
