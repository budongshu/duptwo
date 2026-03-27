package admin

import (
	"datauptwo/app/ldap"
	"datauptwo/global"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type AdminApi struct{}

func NewAdminApi() *AdminApi {
	return &AdminApi{}
}

// GetConfig 获取当前配置
func (a *AdminApi) GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"serve_web": global.CONF.Base.ServeWeb,
			"port":      global.CONF.Base.Port,
			"mode":      global.CONF.Base.Mode,
			"version":   global.CONF.Base.Version,
		},
	})
}

// GetADConfig 获取AD域配置
func (a *AdminApi) GetADConfig(c *gin.Context) {
	cfg := global.CONF.AD
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"enabled":        cfg.Enabled,
			"server":         cfg.Server,
			"port":           cfg.Port,
			"use_ssl":        cfg.UseSSL,
			"base_dn":        cfg.BaseDN,
			"bind_dn":        cfg.BindDN,
			"bind_password":  "", // 不返回密码
			"user_filter":    cfg.UserFilter,
			"auto_register":  cfg.AutoRegister,
			"default_role_id": cfg.DefaultRoleID,
		},
	})
}

// UpdateADConfig 更新AD域配置
func (a *AdminApi) UpdateADConfig(c *gin.Context) {
	var req struct {
		Enabled        bool   `json:"enabled"`
		Server         string `json:"server"`
		Port           int    `json:"port"`
		UseSSL         bool   `json:"use_ssl"`
		BaseDN         string `json:"base_dn"`
		BindDN         string `json:"bind_dn"`
		BindPassword   string `json:"bind_password"`
		UserFilter     string `json:"user_filter"`
		AutoRegister   bool   `json:"auto_register"`
		DefaultRoleID  uint   `json:"default_role_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 如果密码为空，保留原配置
	if req.BindPassword == "" {
		req.BindPassword = global.CONF.AD.BindPassword
	}

	// 更新全局配置
	global.CONF.AD.Enabled = req.Enabled
	global.CONF.AD.Server = req.Server
	global.CONF.AD.Port = req.Port
	global.CONF.AD.UseSSL = req.UseSSL
	global.CONF.AD.BaseDN = req.BaseDN
	global.CONF.AD.BindDN = req.BindDN
	global.CONF.AD.BindPassword = req.BindPassword
	global.CONF.AD.UserFilter = req.UserFilter
	global.CONF.AD.AutoRegister = req.AutoRegister
	global.CONF.AD.DefaultRoleID = req.DefaultRoleID

	// 写回配置文件
	viper.Set("ad.enabled", req.Enabled)
	viper.Set("ad.server", req.Server)
	viper.Set("ad.port", req.Port)
	viper.Set("ad.use_ssl", req.UseSSL)
	viper.Set("ad.base_dn", req.BaseDN)
	viper.Set("ad.bind_dn", req.BindDN)
	viper.Set("ad.bind_password", req.BindPassword)
	viper.Set("ad.user_filter", req.UserFilter)
	viper.Set("ad.auto_register", req.AutoRegister)
	viper.Set("ad.default_role_id", req.DefaultRoleID)

	if err := viper.WriteConfig(); err != nil {
		// 尝试写入配置路径
		configPath := filepath.Join(global.CONF.Base.InstallDir, "conf", "app.yaml")
		if err2 := viper.WriteConfigAs(configPath); err2 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "配置已更新但保存失败: " + err2.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "AD配置已更新"})
}

// TestADConnection 测试AD连接
func (a *AdminApi) TestADConnection(c *gin.Context) {
	var req struct {
		Server       string `json:"server"`
		Port         int    `json:"port"`
		UseSSL       bool   `json:"use_ssl"`
		BaseDN       string `json:"base_dn"`
		BindDN       string `json:"bind_dn"`
		BindPassword string `json:"bind_password"`
		UserFilter   string `json:"user_filter"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	port := req.Port
	if port == 0 {
		port = 389
	}
	userFilter := req.UserFilter
	if userFilter == "" {
		userFilter = "(sAMAccountName=%s)"
	}

	// 用一个测试账号验证连接（使用bind账号本身测试Bind是否成功）
	_, err := ldap.Authenticate(
		req.Server, port, req.UseSSL,
		req.BaseDN, req.BindDN, req.BindPassword,
		userFilter,
		"test-user-do-not-exist", // 用不存在的用户测试，仅验证Bind连接
		"dummy-password",
	)

	// Bind 失败不一定是坏事——可能是用户不存在，而不是连接问题
	// 所以我们只检查连接本身是否建立
	// ldap library 的 Authenticate 会先 Bind 再 Search，Search 失败会返回特定错误
	// 我们认为 "用户不存在" 的错误码 = 连接正常
	if err != nil && err.Error() == "用户在AD中不存在" {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "AD服务器连接成功！用户名不存在说明连接正常（仅验证了 BindDN 账号的连接性）",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "AD服务器连接失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "AD连接测试通过"})
}

// ReloadConfig 重新加载配置文件
func (a *AdminApi) ReloadConfig(c *gin.Context) {
	if err := viper.ReadInConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "重新加载配置失败: " + err.Error(),
		})
		return
	}
	if err := viper.Unmarshal(&global.CONF); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "重新解析配置失败: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "配置已重新加载",
		"data": gin.H{
			"serve_web": global.CONF.Base.ServeWeb,
			"port":      global.CONF.Base.Port,
			"mode":      global.CONF.Base.Mode,
			"version":   global.CONF.Base.Version,
		},
	})
}
