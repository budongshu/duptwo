package main

import (
	"datauptwo/app/model"
	"datauptwo/app/repo"
	"datauptwo/app/service"
	"datauptwo/global"
	"datauptwo/router"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	_ "datauptwo/docs"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// @title datauptwo API
// @version 1.0
// @description datauptwo 数据管理平台 API 文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:18421
// @BasePath /

var (
	configFile string
	serveWeb   bool
	rootCmd    = &cobra.Command{
		Use:   "datauptwo",
		Short: "Data Registry - 数据登记管理平台",
		Long:  `一个灵活、扩展、易用的数据登记管理解决方案`,
		Run:   runServer,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "conf/app.yaml", "配置文件路径")
	rootCmd.PersistentFlags().BoolVarP(&serveWeb, "serve-web", "w", false, "是否启用 Web 前端（覆盖配置文件）")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runServer(cmd *cobra.Command, args []string) {
	// 初始化配置
	initConfig()

	// 初始化日志
	global.InitLogger()
	global.AppLogger.Info("Starting Data Registry server...")

	// 命令行参数覆盖配置文件
	if cmd.Flags().Changed("serve-web") {
		global.CONF.Base.ServeWeb = serveWeb
		global.AppLogger.Info("Command line --serve-web=%v overrides config", serveWeb)
	}

	// 初始化数据库
	initDatabase()

	// 初始化数据
	initData()

	// 启动服务器
	startServer()
}

func initConfig() {
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	if err := viper.Unmarshal(&global.CONF); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	os.Setenv("TZ", global.CONF.Log.TimeZone)
}

func initDatabase() {
	var db *gorm.DB
	var err error

	dbType := global.CONF.Database.Type

	switch dbType {
	case "sqlite":
		dbPath := global.CONF.Database.Path
		dir := path.Dir(dbPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create database directory: %v", err)
		}
		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			global.CONF.Database.User,
			global.CONF.Database.Pass,
			global.CONF.Database.Host,
			global.CONF.Database.Port,
			global.CONF.Database.Name,
		)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			global.CONF.Database.Host,
			global.CONF.Database.User,
			global.CONF.Database.Pass,
			global.CONF.Database.Name,
			global.CONF.Database.Port,
		)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

	default:
		log.Fatalf("Unsupported database type: %s, supported: sqlite, mysql, postgres", dbType)
	}

	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)

	global.DB = db
	global.AppLogger.Info("Database(%s) initialized successfully", dbType)
}

func initData() {
	// 自动迁移表结构
	if err := migrateModels(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化默认角色（先于用户创建，以便关联）
	initDefaultRoles()
	// 初始化默认管理员账号
	initDefaultUser()
	// 初始化默认安全设置
	initDefaultSecuritySettings()

	global.AppLogger.Info("Data initialization completed")
}

// initDefaultUser 初始化默认管理员账号
func initDefaultUser() {
	userRepo := repo.NewUserRepo()
	roleRepo := repo.NewRoleRepo()

	// 查找 admin 角色
	var adminRoleID uint
	if role, err := roleRepo.GetByCode("admin"); err == nil {
		adminRoleID = role.ID
	}

	// 查找 admin 用户
	var admin model.User
	if err := global.DB.Where("username = ?", "admin").First(&admin).Error; err != nil {
		// admin 用户不存在，创建默认管理员
		if adminRoleID == 0 {
			global.AppLogger.Error("Admin role not found, cannot create admin user")
			return
		}
		adminPassword, _ := service.HashPassword("admin123")
		admin := &model.User{
			Username: "admin",
			Password: adminPassword,
			Nickname: "管理员",
			Email:    "admin@example.com",
			Status:   "active",
			RoleID:   adminRoleID,
			GroupID:  0,
		}
		if err := userRepo.Create(admin); err != nil {
			global.AppLogger.Error("Failed to create default admin user: %v", err)
		} else {
			global.AppLogger.Info("Default admin user created: admin / admin123 (roleID=%d)", adminRoleID)
		}
		return
	}

	// admin 用户已存在，确保关联到 admin 角色
	if adminRoleID > 0 && admin.RoleID != adminRoleID {
		global.DB.Model(&admin).Update("role_id", adminRoleID)
		global.AppLogger.Info("Admin user linked to admin role (roleID=%d)", adminRoleID)
	}
}

// initDefaultSecuritySettings 初始化默认安全设置
func initDefaultSecuritySettings() {
	var count int64
	global.DB.Model(&model.SecuritySettings{}).Count(&count)
	if count > 0 {
		global.AppLogger.Info("Security settings already exist, skipping initialization")
		return
	}
	settings := &model.SecuritySettings{
		CaptchaEnabled:          true,
		CaptchaMinLen:           3,
		InactiveAutoDisable:      false,
		InactiveDaysThreshold:   90,
		UserLoginMaxAttempts:    5,
		UserLoginLockMinutes:   30,
		IPLoginMaxAttempts:      20,
		IPLoginLockMinutes:     60,
		IPWhitelist:            "",
		IPBlacklist:            "",
		PasswordExpiryDays:     0,
		PasswordMinLength:      6,
		PasswordRequireUppercase: false,
		PasswordRequireLowercase: false,
		PasswordRequireDigit:   false,
		PasswordRequireSpecial: false,
		SessionTimeoutHours:    24,
	}
	if err := global.DB.Create(settings).Error; err != nil {
		global.AppLogger.Error("Failed to create default security settings: %v", err)
	} else {
		global.AppLogger.Info("Default security settings created")
	}
}

// initDefaultRoles 初始化默认角色
func initDefaultRoles() {
	var count int64
	global.DB.Model(&model.Role{}).Count(&count)
	if count > 0 {
		global.AppLogger.Info("Roles already exist, skipping default role creation")
		return
	}

	roleRepo := repo.NewRoleRepo()

	permissionsJSON := func(perms []string) string {
		b, _ := json.Marshal(perms)
		return string(b)
	}

	roles := []model.Role{
		{
			Name:        "管理员",
			Code:        "admin",
			Description: "拥有系统全部权限",
			Permissions: permissionsJSON([]string{
				"upload:create", "upload:read", "upload:update", "upload:delete", "upload:export",
				"field-config:create", "field-config:read", "field-config:update", "field-config:delete",
				"project:create", "project:read", "project:update", "project:delete",
				"personnel:create", "personnel:read", "personnel:update", "personnel:delete", "personnel:export",
				"user:create", "user:read", "user:update", "user:delete",
				"role:create", "role:read", "role:update", "role:delete",
				"audit:operation:read", "audit:login:read",
				"config:read", "config:update",
				"admin:all",
			}),
			Sort: 1,
		},
		{
			Name:        "数据操作员",
			Code:        "data_operator",
			Description: "负责数据上传、字段配置，可查看项目人员信息",
			Permissions: permissionsJSON([]string{
				"upload:create", "upload:read", "upload:update", "upload:delete", "upload:export",
				"field-config:create", "field-config:read", "field-config:update", "field-config:delete",
				"project:read",
				"personnel:read",
			}),
			Sort: 2,
		},
		{
			Name:        "项目管理员",
			Code:        "project_manager",
			Description: "管理项目、人员，可查看上传记录",
			Permissions: permissionsJSON([]string{
				"project:create", "project:read", "project:update", "project:delete",
				"personnel:create", "personnel:read", "personnel:update", "personnel:delete", "personnel:export",
				"upload:read",
			}),
			Sort: 3,
		},
		{
			Name:        "审计账号",
			Code:        "auditor",
			Description: "仅可查看操作日志和登录日志",
			Permissions: permissionsJSON([]string{
				"upload:read",
				"field-config:read",
				"role:read",
				"audit:operation:read", "audit:login:read",
				"config:read",
			}),
			Sort: 4,
		},
		{
			Name:        "普通查看者",
			Code:        "readonly",
			Description: "仅可查看上传记录和字段配置",
			Permissions: permissionsJSON([]string{
				"upload:read",
				"field-config:read",
			}),
			Sort: 5,
		},
	}

	for _, role := range roles {
		if err := roleRepo.Create(&role); err != nil {
			global.AppLogger.Error("Failed to create default role %s: %v", role.Code, err)
		} else {
			global.AppLogger.Info("Default role created: %s", role.Name)
		}
	}
}

func migrateModels() error {
	// 数据库类型决定索引处理方式
	dbType := global.CONF.Database.Type
	if dbType == "sqlite" {
		// SQLite: DROP INDEX 支持 IF EXISTS
		global.DB.Exec("DROP INDEX IF EXISTS idx_users_username")
	} else {
		// MySQL/PostgreSQL: 先删除再创建
		global.DB.Exec("DROP INDEX IF EXISTS idx_users_username")
	}
	global.DB.Exec("CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)")

	return global.DB.AutoMigrate(
		&model.FieldConfig{},
		&model.UploadRecord{},
		&model.Project{},
		&model.Personnel{},
		&model.User{},
		&model.Role{},
		&model.UserGroup{},
		&model.OperationLog{},
		&model.LoginLog{},
		&model.SecuritySettings{},
		&model.LoginLockout{},
	)
}

func startServer() {
	port := global.CONF.Base.Port
	if port == "" {
		port = "8080"
	}

	r := router.InitRouter()

	global.AppLogger.Info("Server listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
