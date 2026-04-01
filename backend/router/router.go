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
	"datauptwo/app/api/v1/user"
	"datauptwo/app/api/v1/user_group"
	"datauptwo/global"
	"datauptwo/middleware"

	_ "datauptwo/docs"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"

	"github.com/gin-contrib/static"
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

	// Swagger 重定向 /swagger -> /swagger/
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(301, "/swagger/")
	})

	// Swagger API 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 自定义中间件：运行时检查 serve_web 配置，实现热切换
	r.Use(func(c *gin.Context) {
		if !global.CONF.Base.ServeWeb {
			// 禁用 Web 时，对于前端路由返回 404
			// 保留 API 接口正常访问
			path := c.Request.URL.Path
			if path == "/" || (len(path) > 0 && path[0] != '/') || (len(path) > 1 && path[:2] != "/a" && path[:2] != "/p" && path != "/health" && path[:8] != "/swagger") {
				c.AbortWithStatusJSON(404, gin.H{"code": 404, "message": "Web frontend is disabled"})
				return
			}
		}
		c.Next()
	})

	// 静态文件服务（前端）- 根据配置动态决定是否启用
	r.Use(static.Serve("/", static.LocalFile("./cmd/server/web", false)))

	// 公开上传记录接口（无需认证）
	r.POST("/public/upload-records", RouterGroupApp.PublicUploadRecordApi.Create)
	r.GET("/public/upload-records/:serialNo", RouterGroupApp.PublicUploadRecordApi.GetBySerialNo)

	// 公开认证接口
	r.GET("/api/auth/captcha", RouterGroupApp.AuthApi.GetCaptcha)
	r.POST("/api/auth/login", RouterGroupApp.AuthApi.Login)
	r.POST("/api/auth/ad-login", RouterGroupApp.AuthApi.ADLogin)
	r.POST("/api/auth/register", RouterGroupApp.AuthApi.Register)
	r.POST("/api/auth/mfa/verify", RouterGroupApp.AuthApi.MFAVerify)

	// 需要认证的接口
	authGroup := r.Group("/api")
	authGroup.Use(middleware.JWTAuth())
	{
		// 认证相关（无需额外权限，用户自己操作）
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
		authGroup.GET("/upload-records/recent", middleware.RequirePermission("upload:read"), RouterGroupApp.UploadRecordApi.GetRecent)
		authGroup.GET("/upload-records/template", middleware.RequirePermission("upload:read"), RouterGroupApp.UploadRecordApi.GetTemplate)
		authGroup.POST("/upload-records/preview", middleware.RequirePermission("upload:read"), RouterGroupApp.UploadRecordApi.Preview)
		authGroup.GET("/upload-records/uploaders", middleware.RequirePermission("upload:read"), RouterGroupApp.UploadRecordApi.GetUploaderList)
		authGroup.GET("/upload-records/:id", middleware.RequirePermission("upload:read"), RouterGroupApp.UploadRecordApi.GetByID)
		authGroup.PUT("/upload-records", middleware.RequirePermission("upload:update"), RouterGroupApp.UploadRecordApi.Update)
		authGroup.DELETE("/upload-records/:id", middleware.RequirePermission("upload:delete"), RouterGroupApp.UploadRecordApi.Delete)
		authGroup.POST("/upload-records/batch-delete", middleware.RequirePermission("upload:delete"), RouterGroupApp.UploadRecordApi.BatchDelete)
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
		authGroup.GET("/personnels/export", middleware.RequirePermission("personnel:export"), RouterGroupApp.PersonnelApi.Export)
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
		authGroup.GET("/audit/login-logs/export", middleware.RequirePermission("audit:login:read"), RouterGroupApp.AuditApi.ExportLoginLogs)

		// 系统配置接口
		authGroup.GET("/admin/config", middleware.RequirePermission("config:read"), RouterGroupApp.AdminApi.GetConfig)
		authGroup.POST("/admin/config/reload", middleware.RequirePermission("config:update"), RouterGroupApp.AdminApi.ReloadConfig)
		authGroup.GET("/admin/ad-config", middleware.RequirePermission("config:read"), RouterGroupApp.AdminApi.GetADConfig)
		authGroup.POST("/admin/ad-config", middleware.RequirePermission("config:update"), RouterGroupApp.AdminApi.UpdateADConfig)
		authGroup.POST("/admin/ad-config/test", middleware.RequirePermission("config:update"), RouterGroupApp.AdminApi.TestADConnection)

		// 安全设置接口
		authGroup.GET("/admin/security-settings", middleware.RequirePermission("config:read"), RouterGroupApp.SecurityApi.GetSettings)
		authGroup.PUT("/admin/security-settings", middleware.RequirePermission("config:update"), RouterGroupApp.SecurityApi.UpdateSettings)
		authGroup.GET("/admin/security/overview", middleware.RequirePermission("config:read"), RouterGroupApp.SecurityApi.GetSecurityOverview)
		authGroup.GET("/admin/security/locked-users", middleware.RequirePermission("config:read"), RouterGroupApp.SecurityApi.GetLockedUsers)
		authGroup.POST("/admin/security/unlock-user", middleware.RequirePermission("config:update"), RouterGroupApp.SecurityApi.UnlockUser)
		authGroup.GET("/admin/security/locked-ips", middleware.RequirePermission("config:read"), RouterGroupApp.SecurityApi.GetLockedIPs)
		authGroup.POST("/admin/security/unlock-ip", middleware.RequirePermission("config:update"), RouterGroupApp.SecurityApi.UnlockIP)
	}

	// SPA 路由兜底
	r.NoRoute(func(c *gin.Context) {
		if global.CONF.Base.ServeWeb {
			c.File("./cmd/server/web/index.html")
		} else {
			c.AbortWithStatusJSON(404, gin.H{"code": 404, "message": "Web frontend is disabled"})
		}
	})

	return r
}
