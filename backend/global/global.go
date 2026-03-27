package global

import (
	"log"
	"os"
	"time"

	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	VIPER interface{}
)

type Config struct {
	Base      BaseConfig      `mapstructure:"base"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Log       LogConfig       `mapstructure:"log"`
	Session   SessionConfig   `mapstructure:"session"`
	CORS      CORSConfig      `mapstructure:"cors"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	AD        ADConfig        `mapstructure:"ad"`
}

type ADConfig struct {
	Enabled        bool   `mapstructure:"enabled"`         // 是否启用AD认证
	Server         string `mapstructure:"server"`          // AD服务器地址，如 ldap://192.168.1.100:389
	Port           int    `mapstructure:"port"`            // 端口，默认 389
	UseSSL         bool   `mapstructure:"use_ssl"`         // 是否使用LDAPS
	BaseDN         string `mapstructure:"base_dn"`         // 基准DN，如 OU=Users,DC=company,DC=com
	BindDN         string `mapstructure:"bind_dn"`         // 绑定DN，用于搜索用户，如 CN=admin,DC=company,DC=com
	BindPassword   string `mapstructure:"bind_password"`   // 绑定密码
	UserFilter     string `mapstructure:"user_filter"`     // 用户搜索过滤器，如 (sAMAccountName=%s)
	AutoRegister   bool   `mapstructure:"auto_register"`   // 当用户不存在时，是否自动注册AD用户
	DefaultRoleID  uint   `mapstructure:"default_role_id"` // 自动注册AD用户的默认角色ID
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

type BaseConfig struct {
	Mode       string `mapstructure:"mode"`
	Port       string `mapstructure:"port"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	Version    string `mapstructure:"version"`
	Language   string `mapstructure:"language"`
	InstallDir string `mapstructure:"install_dir"`
	ServeWeb   bool   `mapstructure:"serve_web"`
}

type DatabaseConfig struct {
	Type string `mapstructure:"type"`
	Path string `mapstructure:"path"`
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
	Name string `mapstructure:"name"`
}

type LogConfig struct {
	Level     string `mapstructure:"level"`
	TimeZone  string `mapstructure:"time_zone"`
	LogName   string `mapstructure:"log_name"`
	MaxBackup int    `mapstructure:"max_backup"`
	MaxSize   int    `mapstructure:"max_size"`
}

type SessionConfig struct {
	Timeout int    `mapstructure:"timeout"`
	Secret  string `mapstructure:"secret"`
}

type CORSConfig struct {
	AllowOrigins []string `mapstructure:"allow_origins"`
	AllowMethods []string `mapstructure:"allow_methods"`
	AllowHeaders []string `mapstructure:"allow_headers"`
	MaxAge       int      `mapstructure:"max_age"`
}

var CONF *Config

func GetUnix() int64 {
	return time.Now().Unix()
}

func GetTimeNow() time.Time {
	return time.Now()
}

func InitGlobal() {
	CONF = &Config{}
}

type Logger struct {
	*log.Logger
}

var AppLogger *Logger

func InitLogger() {
	AppLogger = &Logger{
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.Printf("[INFO] "+format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.Printf("[ERROR] "+format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if CONF != nil && CONF.Log.Level == "debug" {
		l.Printf("[DEBUG] "+format, v...)
	}
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.Printf("[WARN] "+format, v...)
}
