package router

import (
	"datauptwo/app/api/v1/admin"
	"datauptwo/app/api/v1/auth"
	"datauptwo/app/api/v1/audit"
	fieldconfig "datauptwo/app/api/v1/field_config"
	personnelapi "datauptwo/app/api/v1/personnel"
	projectapi "datauptwo/app/api/v1/project"
	"datauptwo/app/api/v1/role"
	uploadRecord "datauptwo/app/api/v1/upload_record"
	"datauptwo/app/api/v1/sync"
	"datauptwo/app/api/v1/user"
	"datauptwo/app/api/v1/user_group"
	"datauptwo/global"
	"datauptwo/middleware"
	"os"
	"path/filepath"
	"strings"

	_ "datauptwo/docs"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"

	"github.com/gin-gonic/gin"
)

type RouterGroup struct {
	UploadRecordApi       *uploadRecord.UploadRecordApi
	PublicUploadRecordApi *uploadRecord.PublicUploadRecordApi
	FieldConfigApi        *fieldconfig.FieldConfigApi
	ProjectApi            *projectapi.ProjectApi
	PersonnelApi          *personnelapi.PersonnelApi
	AuthApi               *auth.AuthApi
	UserApi               *user.UserApi
	RoleApi               *role.RoleApi
	UserGroupApi          *user_group.UserGroupApi
	AuditApi              *audit.AuditApi
	AdminApi              *admin.AdminApi
	SecurityApi           *admin.SecurityApi
	SyncApi               *sync.SyncApi
}

var RouterGroupApp = &RouterGroup{
	UploadRecordApi:       uploadRecord.NewUploadRecordApi(),
	PublicUploadRecordApi: uploadRecord.NewPublicUploadRecordApi(),
	FieldConfigApi:        fieldconfig.NewFieldConfigApi(),
	ProjectApi:            projectapi.NewProjectApi(),
	PersonnelApi:          personnelapi.NewPersonnelApi(),
	AuthApi:               auth.NewAuthApi(),
	UserApi:               user.NewUserApi(),
	RoleApi:               role.NewRoleApi(),
	UserGroupApi:          user_group.NewUserGroupApi(),
	AuditApi:              audit.NewAuditApi(),
	AdminApi:              admin.NewAdminApi(),
	SecurityApi:           admin.NewSecurityApi(),
	SyncApi:               sync.NewSyncApi(),
}

// resolveWebRoot resolves web root path.
// Priority: config value > ./cmd/server/web (dev) > ./web (release).
// Relative paths are resolved relative to the project root (parent of conf/).
func resolveWebRoot() string {
	webRoot := global.CONF.Base.WebRoot
	if webRoot == "" {
		for _, candidate := range []string{"./cmd/server/web", "./web"} {
			if _, err := os.Stat(candidate); err == nil {
				webRoot = candidate
				break
			}
		}
	}
	if webRoot == "" {
		return ""
	}
	if filepath.IsAbs(webRoot) {
		return webRoot
	}
	if global.CONF.Base.ConfigFile == "" {
		return webRoot
	}
	confDir := filepath.Dir(global.CONF.Base.ConfigFile)
	projectRoot := filepath.Dir(confDir)
	return filepath.Join(projectRoot, webRoot)
}

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger：统一由 ginSwagger 处理裸路径（gin-swagger v1.6.1 对 /swagger/ 等裸路径有匹配 bug，
	// 直接在 handler 内部重定向到 /swagger/index.html，避免 Gin 路由冲突）
	swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(301, "/swagger/index.html")
	})
	r.GET("/swagger/*any", func(c *gin.Context) {
		path := c.Param("any")
		if path == "" || path == "/" {
			c.Redirect(301, "/swagger/index.html")
			return
		}
		swaggerHandler(c)
	})

	// 公开上传记录接口（无需认证）
	r.POST("/public/upload-records", RouterGroupApp.PublicUploadRecordApi.Create)
	r.GET("/public/upload-records/:serialNo", RouterGroupApp.PublicUploadRecordApi.GetBySerialNo)
	r.PUT("/public/upload-records/:serialNo", RouterGroupApp.PublicUploadRecordApi.UpdateBySerialNo)
	r.POST("/public/upload-records/:serialNo/update", RouterGroupApp.PublicUploadRecordApi.UpdateBySerialNoPost)

	// 公开认证接口
	r.GET("/api/auth/captcha", RouterGroupApp.AuthApi.GetCaptcha)
	r.POST("/api/auth/login", RouterGroupApp.AuthApi.Login)
	r.POST("/api/auth/ad-login", RouterGroupApp.AuthApi.ADLogin)
	r.POST("/api/auth/register", RouterGroupApp.AuthApi.Register)
	r.POST("/api/auth/mfa/verify", RouterGroupApp.AuthApi.MFAVerify)
	r.GET("/api/auth/registration-status", RouterGroupApp.AuthApi.GetRegistrationStatus)

	// 公开同步接口（使用API Key认证）
	r.POST("/api/sync/register", RouterGroupApp.SyncApi.Register)
	r.POST("/api/sync/upload-records", middleware.ApiKeyAuth(), RouterGroupApp.SyncApi.UploadRecords)
	r.POST("/api/sync/heartbeat", middleware.ApiKeyAuth(), RouterGroupApp.SyncApi.Heartbeat)

	// 需要认证的接口
	authGroup := r.Group("/api")
	authGroup.Use(middleware.JWTAuth())
	{
		// 认证相关
		authGroup.GET("/auth/current", RouterGroupApp.AuthApi.GetCurrentUser)
		authGroup.PUT("/auth/profile", RouterGroupApp.AuthApi.UpdateProfile)
		authGroup.POST("/auth/change-password", RouterGroupApp.AuthApi.ChangePassword)
		authGroup.GET("/auth/mfa/status", RouterGroupApp.AuthApi.GetMFAStatus)
		authGroup.POST("/auth/mfa/enable", middleware.RequirePermission("admin:all"), RouterGroupApp.AuthApi.EnableMFA)
		authGroup.POST("/auth/mfa/disable", middleware.RequirePermission("admin:all"), RouterGroupApp.AuthApi.DisableMFA)

		// 上传记录管理接口
		authGroup.GET("/upload-records", middleware.RequirePermission("upload:read"), RouterGroupApp.UploadRecordApi.List)
		authGroup.POST("/upload-records", middleware.RequirePermission("upload:create"), RouterGroupApp.UploadRecordApi.Create)
		authGroup.GET("/upload-records/export", middleware.RequirePermission("upload:export"), RouterGroupApp.UploadRecordApi.Export)
		authGroup.GET("/upload-records/statistics", middleware.RequirePermission("upload:read"), RouterGroupApp.UploadRecordApi.GetStatistics)
		authGroup.GET("/upload-records/disk-labels", middleware.RequirePermission("upload:read"), RouterGroupApp.UploadRecordApi.GetDiskLabelStatuses)
		authGroup.GET("/upload-records/recent", middleware.RequirePermission("upload:read"), RouterGroupApp.UploadRecordApi.GetRecent)
		authGroup.GET("/upload-records/template", middleware.RequirePermission("upload:read"), RouterGroupApp.UploadRecordApi.GetTemplate)
		authGroup.POST("/upload-records/preview", middleware.RequirePermission("upload:read"), RouterGroupApp.UploadRecordApi.Preview)
		authGroup.GET("/upload-records/uploaders", middleware.RequirePermission("upload:read"), RouterGroupApp.UploadRecordApi.GetUploaderList)
		authGroup.GET("/upload-records/:id", middleware.RequirePermission("upload:read"), RouterGroupApp.UploadRecordApi.GetByID)
		authGroup.PUT("/upload-records", middleware.RequirePermission("upload:update"), RouterGroupApp.UploadRecordApi.Update)
		authGroup.DELETE("/upload-records/:id", middleware.RequirePermission("upload:delete"), RouterGroupApp.UploadRecordApi.Delete)
		authGroup.POST("/upload-records/batch-delete", middleware.RequirePermission("upload:delete"), RouterGroupApp.UploadRecordApi.BatchDelete)
		authGroup.POST("/upload-records/batch-update-status", middleware.RequirePermission("upload:update"), RouterGroupApp.UploadRecordApi.BatchUpdateStatus)
		authGroup.POST("/upload-records/import", middleware.RequirePermission("upload:create"), RouterGroupApp.UploadRecordApi.Import)

		// 字段配置接口
		authGroup.GET("/field-configs", middleware.RequirePermission("field-config:read"), RouterGroupApp.FieldConfigApi.List)
		authGroup.POST("/field-configs", middleware.RequirePermission("field-config:create"), RouterGroupApp.FieldConfigApi.Create)
		authGroup.GET("/field-configs/all", middleware.RequirePermission("field-config:read"), RouterGroupApp.FieldConfigApi.GetAllEnabled)
		authGroup.GET("/field-configs/:id", middleware.RequirePermission("field-config:read"), RouterGroupApp.FieldConfigApi.GetByID)
		authGroup.PUT("/field-configs", middleware.RequirePermission("field-config:update"), RouterGroupApp.FieldConfigApi.Update)
		authGroup.DELETE("/field-configs/:id", middleware.RequirePermission("field-config:delete"), RouterGroupApp.FieldConfigApi.Delete)
		authGroup.POST("/field-configs/batch-delete", middleware.RequirePermission("field-config:delete"), RouterGroupApp.FieldConfigApi.BatchDelete)

		// 项目管理接口
		authGroup.GET("/projects", middleware.RequirePermission("project:read"), RouterGroupApp.ProjectApi.List)
		authGroup.POST("/projects", middleware.RequirePermission("project:create"), RouterGroupApp.ProjectApi.Create)
		authGroup.GET("/projects/simple", middleware.RequirePermission("project:read"), RouterGroupApp.ProjectApi.ListSimple)
		authGroup.GET("/projects/kanban", middleware.RequirePermission("project:read"), RouterGroupApp.ProjectApi.ListKanban)
		authGroup.GET("/projects/:id", middleware.RequirePermission("project:read"), RouterGroupApp.ProjectApi.GetByID)
		authGroup.PUT("/projects", middleware.RequirePermission("project:update"), RouterGroupApp.ProjectApi.Update)
		authGroup.DELETE("/projects/:id", middleware.RequirePermission("project:delete"), RouterGroupApp.ProjectApi.Delete)
		authGroup.POST("/projects/batch-delete", middleware.RequirePermission("project:delete"), RouterGroupApp.ProjectApi.BatchDelete)

		// 人员管理接口
		authGroup.GET("/personnels", middleware.RequirePermission("personnel:read"), RouterGroupApp.PersonnelApi.List)
		authGroup.POST("/personnels", middleware.RequirePermission("personnel:create"), RouterGroupApp.PersonnelApi.Create)
		authGroup.GET("/personnels/all", middleware.RequirePermission("personnel:read"), RouterGroupApp.PersonnelApi.ListAll)
		authGroup.GET("/personnels/statistics", middleware.RequirePermission("personnel:read"), RouterGroupApp.PersonnelApi.Statistics)
		authGroup.GET("/personnels/export", middleware.RequirePermission("personnel:export"), RouterGroupApp.PersonnelApi.Export)
		authGroup.GET("/personnels/template", middleware.RequirePermission("personnel:read"), RouterGroupApp.PersonnelApi.GetTemplate)
		authGroup.POST("/personnels/preview", middleware.RequirePermission("personnel:read"), RouterGroupApp.PersonnelApi.Preview)
		authGroup.POST("/personnels/import", middleware.RequirePermission("personnel:create"), RouterGroupApp.PersonnelApi.Import)
		authGroup.GET("/personnels/:id", middleware.RequirePermission("personnel:read"), RouterGroupApp.PersonnelApi.GetByID)
		authGroup.PUT("/personnels", middleware.RequirePermission("personnel:update"), RouterGroupApp.PersonnelApi.Update)
		authGroup.DELETE("/personnels/:id", middleware.RequirePermission("personnel:delete"), RouterGroupApp.PersonnelApi.Delete)
		authGroup.POST("/personnels/batch-delete", middleware.RequirePermission("personnel:delete"), RouterGroupApp.PersonnelApi.BatchDelete)

		// 用户管理接口
		authGroup.GET("/users", middleware.RequirePermission("user:read"), RouterGroupApp.UserApi.List)
		authGroup.POST("/users", middleware.RequirePermission("user:create"), RouterGroupApp.UserApi.Create)
		authGroup.GET("/users/all", middleware.RequirePermission("user:read"), RouterGroupApp.UserApi.GetAll)
		authGroup.GET("/users/:id", middleware.RequirePermission("user:read"), RouterGroupApp.UserApi.GetByID)
		authGroup.PUT("/users", middleware.RequirePermission("user:update"), RouterGroupApp.UserApi.Update)
		authGroup.DELETE("/users/:id", middleware.RequirePermission("user:delete"), RouterGroupApp.UserApi.Delete)
		authGroup.POST("/users/batch-delete", middleware.RequirePermission("user:delete"), RouterGroupApp.UserApi.BatchDelete)
		authGroup.POST("/users/batch-update-role", middleware.RequirePermission("user:update"), RouterGroupApp.UserApi.BatchUpdateRole)
		authGroup.GET("/users/export", middleware.RequirePermission("user:read"), RouterGroupApp.UserApi.Export)
		authGroup.GET("/users/template", middleware.RequirePermission("user:create"), RouterGroupApp.UserApi.GetTemplate)
		authGroup.POST("/users/preview", middleware.RequirePermission("user:create"), RouterGroupApp.UserApi.Preview)
		authGroup.POST("/users/import", middleware.RequirePermission("user:create"), RouterGroupApp.UserApi.Import)
		authGroup.POST("/users/reset-password", middleware.RequirePermission("admin:all"), RouterGroupApp.UserApi.ResetPassword)
		authGroup.POST("/users/reset-mfa", middleware.RequirePermission("admin:all"), RouterGroupApp.UserApi.ResetMFA)
		authGroup.GET("/users/generate-mfa-secret", middleware.RequirePermission("admin:all"), RouterGroupApp.UserApi.GenerateMFASecret)
		authGroup.POST("/users/admin-enable-mfa", middleware.RequirePermission("admin:all"), RouterGroupApp.UserApi.AdminEnableMFA)

		// 角色管理接口
		authGroup.GET("/roles", middleware.RequirePermission("role:read"), RouterGroupApp.RoleApi.List)
		authGroup.POST("/roles", middleware.RequirePermission("role:create"), RouterGroupApp.RoleApi.Create)
		authGroup.GET("/roles/all", middleware.RequirePermission("role:read"), RouterGroupApp.RoleApi.GetAll)
		authGroup.GET("/roles/:id", middleware.RequirePermission("role:read"), RouterGroupApp.RoleApi.GetByID)
		authGroup.PUT("/roles", middleware.RequirePermission("role:update"), RouterGroupApp.RoleApi.Update)
		authGroup.DELETE("/roles/:id", middleware.RequirePermission("role:delete"), RouterGroupApp.RoleApi.Delete)
		authGroup.POST("/roles/batch-delete", middleware.RequirePermission("role:delete"), RouterGroupApp.RoleApi.BatchDelete)

		// 用户组管理接口
		authGroup.GET("/user-groups", middleware.RequirePermission("role:read"), RouterGroupApp.UserGroupApi.List)
		authGroup.POST("/user-groups", middleware.RequirePermission("role:create"), RouterGroupApp.UserGroupApi.Create)
		authGroup.GET("/user-groups/all", middleware.RequirePermission("role:read"), RouterGroupApp.UserGroupApi.GetAll)
		authGroup.GET("/user-groups/:id", middleware.RequirePermission("role:read"), RouterGroupApp.UserGroupApi.GetByID)
		authGroup.PUT("/user-groups", middleware.RequirePermission("role:update"), RouterGroupApp.UserGroupApi.Update)
		authGroup.DELETE("/user-groups/:id", middleware.RequirePermission("role:delete"), RouterGroupApp.UserGroupApi.Delete)
		authGroup.POST("/user-groups/batch-delete", middleware.RequirePermission("role:delete"), RouterGroupApp.UserGroupApi.BatchDelete)

		// 日志管理接口
		authGroup.GET("/audit/operation-logs", middleware.RequirePermission("audit:operation:read"), RouterGroupApp.AuditApi.ListOperationLogs)
		authGroup.GET("/audit/operation-logs/export", middleware.RequirePermission("audit:operation:read"), RouterGroupApp.AuditApi.ExportOperationLogs)
		authGroup.GET("/audit/login-logs", middleware.RequirePermission("audit:login:read"), RouterGroupApp.AuditApi.ListLoginLogs)
		authGroup.GET("/audit/login-logs/export", middleware.RequirePermission("audit:login:read"), RouterGroupApp.AuditApi.ExportOperationLogs)

		// 系统配置接口
		authGroup.GET("/admin/config", middleware.RequirePermission("config:read"), RouterGroupApp.AdminApi.GetConfig)
		authGroup.POST("/admin/config/reload", middleware.RequirePermission("config:update"), RouterGroupApp.AdminApi.ReloadConfig)
		authGroup.GET("/admin/ad-config", middleware.RequirePermission("config:read"), RouterGroupApp.AdminApi.GetADConfig)
		authGroup.POST("/admin/ad-config", middleware.RequirePermission("config:update"), RouterGroupApp.AdminApi.UpdateADConfig)
		authGroup.POST("/admin/ad-config/test", middleware.RequirePermission("config:update"), RouterGroupApp.AdminApi.TestADConnection)
		authGroup.POST("/admin/ad-users/sync", middleware.RequirePermission("config:update"), RouterGroupApp.AdminApi.SyncADUsers)
		authGroup.POST("/admin/ad-users/reset-all", middleware.RequirePermission("config:update"), RouterGroupApp.AdminApi.ResetAllADUsers)
		authGroup.GET("/admin/ad-users", middleware.RequirePermission("config:read"), RouterGroupApp.AdminApi.GetADUsers)

		// 安全设置接口
		authGroup.GET("/admin/security-settings", middleware.RequirePermission("config:read"), RouterGroupApp.SecurityApi.GetSettings)
		authGroup.PUT("/admin/security-settings", middleware.RequirePermission("config:update"), RouterGroupApp.SecurityApi.UpdateSettings)
		authGroup.GET("/admin/security/overview", middleware.RequirePermission("config:read"), RouterGroupApp.SecurityApi.GetSecurityOverview)
		authGroup.GET("/admin/security/locked-users", middleware.RequirePermission("config:read"), RouterGroupApp.SecurityApi.GetLockedUsers)
		authGroup.POST("/admin/security/unlock-user", middleware.RequirePermission("config:update"), RouterGroupApp.SecurityApi.UnlockUser)
		authGroup.GET("/admin/security/locked-ips", middleware.RequirePermission("config:read"), RouterGroupApp.SecurityApi.GetLockedIPs)
		authGroup.POST("/admin/security/unlock-ip", middleware.RequirePermission("config:update"), RouterGroupApp.SecurityApi.UnlockIP)

		// 同步管理接口
		authGroup.GET("/sync/stations", middleware.RequirePermission("config:read"), RouterGroupApp.SyncApi.ListStations)
		authGroup.POST("/sync/stations", middleware.RequirePermission("config:update"), RouterGroupApp.SyncApi.CreateStation)
		authGroup.GET("/sync/stations/all", middleware.RequirePermission("config:read"), RouterGroupApp.SyncApi.GetAllStations)
		authGroup.GET("/sync/stations/:id", middleware.RequirePermission("config:read"), RouterGroupApp.SyncApi.GetStation)
		authGroup.PUT("/sync/stations", middleware.RequirePermission("config:update"), RouterGroupApp.SyncApi.UpdateStation)
		authGroup.DELETE("/sync/stations/:id", middleware.RequirePermission("config:update"), RouterGroupApp.SyncApi.DeleteStation)
		authGroup.POST("/sync/stations/:id/reset-key", middleware.RequirePermission("config:update"), RouterGroupApp.SyncApi.ResetApiKey)
		authGroup.GET("/sync/history", middleware.RequirePermission("config:read"), RouterGroupApp.SyncApi.GetHistory)
		authGroup.GET("/sync/history/:id", middleware.RequirePermission("config:read"), RouterGroupApp.SyncApi.GetHistoryDetails)
		authGroup.GET("/sync/station-summaries", middleware.RequirePermission("config:read"), RouterGroupApp.SyncApi.GetStationSummaries)
		authGroup.GET("/sync/status", middleware.RequirePermission("config:read"), RouterGroupApp.SyncApi.GetStatus)
	}

	// 静态文件服务 + SPA 兜底
	// 注意：不能使用 r.Static("/", root)，因为它会注册 /*filepath 通配路由，
	// 与已有的 /api/* 等前缀路由冲突。故改用 NoRoute 统一兜底。
	if global.CONF.Base.ServeWeb {
		webRoot := resolveWebRoot()
		if webRoot != "" {
			r.NoRoute(func(c *gin.Context) {
				path := c.Request.URL.Path
				// index.html 禁用缓存（避免浏览器缓存旧前端）
				if path == "/" {
					c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
					c.Header("Pragma", "no-cache")
					c.Header("Expires", "0")
				}
				// 文件存在则直接服务，否则返回 index.html（SPA 兜底）
				fullPath := filepath.Join(webRoot, filepath.Clean(strings.TrimPrefix(path, "/")))
				if _, err := os.Stat(fullPath); err == nil {
					c.File(fullPath)
					return
				}
				c.File(filepath.Join(webRoot, "index.html"))
			})
		}
	} else {
		r.NoRoute(func(c *gin.Context) {
			c.AbortWithStatusJSON(404, gin.H{"code": 404, "message": "Web frontend is disabled"})
		})
	}

	return r
}
