package admin

import (
	"datauptwo/app/ldap"
	"datauptwo/app/model"
	"datauptwo/app/repo"
	"datauptwo/global"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminApi struct {
	securitySettingsRepo *repo.SecuritySettingsRepo
}

func NewAdminApi() *AdminApi {
	return &AdminApi{
		securitySettingsRepo: repo.NewSecuritySettingsRepo(),
	}
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
	settings, err := a.securitySettingsRepo.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取AD配置失败: " + err.Error()})
		return
	}
	if settings == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": gin.H{
				"enabled":        false,
				"server":         "",
				"port":           389,
				"use_ssl":        false,
				"base_dn":        "",
				"bind_dn":        "",
				"bind_password":  "",
				"user_filter":    "(sAMAccountName=%s)",
				"auto_register":  false,
				"default_role_id": 0,
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"enabled":        settings.ADEnabled,
			"server":         settings.ADServer,
			"port":           settings.ADPort,
			"use_ssl":        settings.ADUseSSL,
			"base_dn":        settings.ADBaseDN,
			"bind_dn":        settings.ADBindDN,
			"bind_password":  "", // 不返回密码
			"user_filter":    settings.ADUserFilter,
			"auto_register":  settings.ADAutoRegister,
			"default_role_id": settings.ADDefaultRoleID,
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

	// 获取当前配置
	settings, err := a.securitySettingsRepo.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取配置失败: " + err.Error()})
		return
	}
	if settings == nil {
		settings = &model.SecuritySettings{}
	}

	// 如果密码为空，保留原配置（从 global.CONF.AD 获取）
	if req.BindPassword == "" {
		req.BindPassword = global.CONF.AD.BindPassword
	}

	// 更新 AD 配置
	settings.ADEnabled = req.Enabled
	settings.ADServer = req.Server
	settings.ADPort = req.Port
	settings.ADUseSSL = req.UseSSL
	settings.ADBaseDN = req.BaseDN
	settings.ADBindDN = req.BindDN
	settings.ADBindPassword = req.BindPassword
	settings.ADUserFilter = req.UserFilter
	settings.ADAutoRegister = req.AutoRegister
	settings.ADDefaultRoleID = req.DefaultRoleID

	// 保存到数据库
	if err := a.securitySettingsRepo.SaveADConfig(settings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存AD配置失败: " + err.Error()})
		return
	}

	// 同步更新内存中的全局配置
	global.CONF.AD.Enabled = settings.ADEnabled
	global.CONF.AD.Server = settings.ADServer
	global.CONF.AD.Port = settings.ADPort
	global.CONF.AD.UseSSL = settings.ADUseSSL
	global.CONF.AD.BaseDN = settings.ADBaseDN
	global.CONF.AD.BindDN = settings.ADBindDN
	global.CONF.AD.BindPassword = settings.ADBindPassword
	global.CONF.AD.UserFilter = settings.ADUserFilter
	global.CONF.AD.AutoRegister = settings.ADAutoRegister
	global.CONF.AD.DefaultRoleID = settings.ADDefaultRoleID

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
	// 如果测试密码为空，从内存配置获取
	bindPassword := req.BindPassword
	if bindPassword == "" {
		bindPassword = global.CONF.AD.BindPassword
	}

	// 用一个测试账号验证连接（使用bind账号本身测试Bind是否成功）
	_, err := ldap.Authenticate(
		req.Server, port, req.UseSSL,
		req.BaseDN, req.BindDN, bindPassword,
		userFilter,
		"test-user-do-not-exist", // 用不存在的用户测试，仅验证Bind连接
		"dummy-password",
	)

	// Bind 失败不一定是坏事——可能是用户不存在，而不是连接问题
	// 所以我们只检查连接本身是否建立
	// ldap library 的 Authenticate 会先 Bind 再 Search，Search 失败会返回特定错误
	// 我们认为 "用户不存在" 的错误码 = 连接正常
	if err != nil && err.Error() == "用户在目录中不存在" {
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
	// 从数据库重新加载 AD 配置
	settings, err := a.securitySettingsRepo.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "重新加载配置失败: " + err.Error(),
		})
		return
	}
	if settings != nil {
		global.CONF.AD.Enabled = settings.ADEnabled
		global.CONF.AD.Server = settings.ADServer
		global.CONF.AD.Port = settings.ADPort
		global.CONF.AD.UseSSL = settings.ADUseSSL
		global.CONF.AD.BaseDN = settings.ADBaseDN
		global.CONF.AD.BindDN = settings.ADBindDN
		global.CONF.AD.BindPassword = settings.ADBindPassword
		global.CONF.AD.UserFilter = settings.ADUserFilter
		global.CONF.AD.AutoRegister = settings.ADAutoRegister
		global.CONF.AD.DefaultRoleID = settings.ADDefaultRoleID
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

// SyncADUsers 同步AD用户到本地数据库
func (a *AdminApi) SyncADUsers(c *gin.Context) {
	var req struct {
		Mode string `json:"mode"` // "full" 全量 / "incremental" 增量（默认）
	}
	c.ShouldBindJSON(&req)
	syncMode := req.Mode
	if syncMode == "" {
		syncMode = "incremental"
	}

	cfg := global.CONF.AD
	if !cfg.Enabled {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "AD认证未启用"})
		return
	}

	// 从AD搜索所有用户
	adUsers, err := ldap.SearchAllUsers(
		cfg.Server, cfg.Port, cfg.UseSSL,
		cfg.BaseDN, cfg.BindDN, cfg.BindPassword,
		cfg.UserFilter,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "同步失败: " + err.Error()})
		return
	}

	userRepo := repo.NewUserRepo()
	var existingUsers []model.User
	global.DB.Unscoped().Where("source = ?", "AD").Find(&existingUsers)
	existingMap := make(map[string]*model.User)
	for i := range existingUsers {
		existingMap[existingUsers[i].Username] = &existingUsers[i]
	}

	var created, updated, total, skipped int
	now := time.Now()
	adUsernameSet := make(map[string]bool)

	for _, adUser := range adUsers {
		total++
		adUsernameSet[adUser.Username] = true
		existing, exists := existingMap[adUser.Username]

		if !exists {
			// 新增用户
			defaultRoleID := cfg.DefaultRoleID
			if defaultRoleID == 0 {
				defaultRoleID = 3
			}
			user := &model.User{
				Username: adUser.Username,
				Nickname: adUser.Nickname,
				Email:    adUser.Email,
				Source:   "AD",
				ADDN:     adUser.DN,
				Status:   "active",
				RoleID:   defaultRoleID,
			}
			if err := userRepo.Create(user); err != nil {
				global.AppLogger.Error("Failed to create AD user %s: %v", adUser.Username, err)
				continue
			}
			created++
		} else {
			changed := false
			// 如果用户之前被软删除，重新同步时恢复
			if existing.IsDeleted {
				global.DB.Unscoped().Model(&model.User{}).Where("id = ?", existing.ID).Update("is_deleted", false)
				existing.IsDeleted = false
				existing.Status = "active"
				changed = true
			}
			if adUser.Email != "" && existing.Email != adUser.Email {
				existing.Email = adUser.Email
				changed = true
			}
			if adUser.Nickname != "" && existing.Nickname != adUser.Nickname {
				existing.Nickname = adUser.Nickname
				changed = true
			}
			if adUser.DN != "" && existing.ADDN != adUser.DN {
				existing.ADDN = adUser.DN
				changed = true
			}
			if existing.Status == "inactive" {
				existing.Status = "active"
				changed = true
			}
			if changed {
				if err := userRepo.Update(existing); err != nil {
					global.AppLogger.Error("Failed to update AD user %s: %v", adUser.Username, err)
				} else {
					updated++
				}
			} else {
				skipped++
			}
		}
	}

	// 全量同步：标记本地存在但AD中已删除的用户为禁用
	var disabled, deletedCount int
	if syncMode == "full" {
		for _, existing := range existingUsers {
			if !adUsernameSet[existing.Username] {
				if existing.Status == "active" {
					global.DB.Model(&existing).Update("status", "inactive")
					disabled++
				}
				// 完全从AD删除的用户，执行软删除
				if !existing.IsDeleted {
					global.DB.Model(&existing).Update("is_deleted", true)
					deletedCount++
				}
			}
		}
	}

	// 更新同步记录到 security_settings 表
	settings, err := a.securitySettingsRepo.Get()
	if err != nil {
		global.AppLogger.Error("Failed to get security settings: %v", err)
	} else {
		if settings == nil {
			settings = &model.SecuritySettings{}
		}
		settings.ADLastSyncAt = &now
		settings.ADLastSyncCount = total
		if saveErr := a.securitySettingsRepo.SaveADConfig(settings); saveErr != nil {
			global.AppLogger.Error("Failed to save sync info: %v", saveErr)
		}
	}

	// 同步到内存配置
	global.CONF.AD.LastSyncAt = &now
	global.CONF.AD.LastSyncCount = total

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "同步完成",
		"data": gin.H{
			"total":    total,
			"created":  created,
			"updated":  updated,
			"skipped":  skipped,
			"disabled": disabled,
			"deleted":  deletedCount,
			"mode":     syncMode,
			"lastSyncAt": now.Format("2006-01-02 15:04:05"),
		},
	})
}

// ResetAllADUsers 重置所有AD同步用户（删除所有Source=AD的用户）
func (a *AdminApi) ResetAllADUsers(c *gin.Context) {
	cfg := global.CONF.AD
	if !cfg.Enabled {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "AD认证未启用"})
		return
	}

	// 删除所有 Source = "AD" 的用户
	result := global.DB.Where("source = ?", "AD").Delete(&model.User{})
	count := result.RowsAffected

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "重置完成",
		"data": gin.H{
			"count": count,
		},
	})
}

// GetADUsers 获取AD用户列表（预览，不写入本地）
func (a *AdminApi) GetADUsers(c *gin.Context) {
	cfg := global.CONF.AD
	if !cfg.Enabled {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "AD认证未启用"})
		return
	}

	adUsers, err := ldap.SearchAllUsers(
		cfg.Server, cfg.Port, cfg.UseSSL,
		cfg.BaseDN, cfg.BindDN, cfg.BindPassword,
		cfg.UserFilter,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户列表失败: " + err.Error()})
		return
	}

	// 获取本地已存在的AD用户
	var localUsers []model.User
	global.DB.Where("source = ?", "AD").Find(&localUsers)
	localMap := make(map[string]bool)
	for _, u := range localUsers {
		localMap[u.Username] = true
	}

	// 标记是否已同步
	type ADUserPreview struct {
		Username    string `json:"username"`
		Nickname    string `json:"nickname"`
		Email       string `json:"email"`
		DN          string `json:"dn"`
		Synced      bool   `json:"synced"`
		LocalStatus string `json:"localStatus"`
	}
	// 预先构建本地用户状态映射
	localStatusMap := make(map[string]string)
	for _, lu := range localUsers {
		localStatusMap[lu.Username] = lu.Status
	}
	previews := make([]ADUserPreview, 0, len(adUsers))
	for _, u := range adUsers {
		preview := ADUserPreview{
			Username:    u.Username,
			Nickname:    u.Nickname,
			Email:       u.Email,
			DN:          u.DN,
			Synced:      localMap[u.Username],
			LocalStatus: localStatusMap[u.Username],
		}
		previews = append(previews, preview)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total":   len(adUsers),
			"synced":  len(localUsers),
			"users":   previews,
		},
	})
}