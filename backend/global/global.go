package global

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB        *gorm.DB
	VIPER     interface{}
	Scheduler interface{} // 全局同步调度器
)

// DBType 返回当前数据库类型：sqlite / mysql / pgsql
func DBType() string {
	if CONF == nil || CONF.Database.Type == "" {
		return "sqlite"
	}
	return CONF.Database.Type
}

// CenterProbeState Center 主动探测状态（内存存储，程序生命周期内有效）
type CenterProbeState struct {
	mu        sync.RWMutex
	Connected map[uint]bool // stationID -> 是否在线
}

var centerProbe = &CenterProbeState{Connected: make(map[uint]bool)}

// GetProbeConnected 获取某站点的探测在线状态
func GetProbeConnected(stationID uint) bool {
	centerProbe.mu.RLock()
	defer centerProbe.mu.RUnlock()
	return centerProbe.Connected[stationID]
}

// SetProbeConnected 设置某站点的探测在线状态
func SetProbeConnected(stationID uint, connected bool) {
	centerProbe.mu.Lock()
	defer centerProbe.mu.Unlock()
	centerProbe.Connected[stationID] = connected
}

// SetProbeConnectedBatch 批量设置探测状态
func SetProbeConnectedBatch(results map[uint]bool) {
	centerProbe.mu.Lock()
	defer centerProbe.mu.Unlock()
	for id, connected := range results {
		centerProbe.Connected[id] = connected
	}
}

// GetSyncScheduler 获取同步调度器
func GetSyncScheduler() interface{} {
	return Scheduler
}

// SetSyncScheduler 设置同步调度器
func SetSyncScheduler(s interface{}) {
	Scheduler = s
}

type Config struct {
	Base      BaseConfig      `mapstructure:"base"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Log       LogConfig       `mapstructure:"log"`
	Session   SessionConfig   `mapstructure:"session"`
	CORS      CORSConfig      `mapstructure:"cors"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	AD        ADConfig        `mapstructure:"ad"`
	Sync      SyncConfig      `mapstructure:"sync"`
}

// SyncConfig 同步配置（mode 互斥：center / agent / ""）
type SyncConfig struct {
	// Mode 同步模式：center（中心站点） / agent（代理节点） / ""（不启用）
	Mode   string       `mapstructure:"mode"`
	Center CenterConfig `mapstructure:"center"`
	Agent  AgentConfig  `mapstructure:"agent"`
	// API Key（Agent 运行时自动从 Center 获取，持久化在此）
	APIKey string `mapstructure:"api_key"`
}

// CenterConfig Center 站点配置
type CenterConfig struct {
	HeartbeatTimeout int `mapstructure:"heartbeat_timeout"` // Agent 心跳超时时间（秒），超过此时间无心跳视为离线，默认 60
}

// AgentConfig Agent 站点配置
type AgentConfig struct {
	CenterURL       string           `mapstructure:"center_url"`       // Center 地址
	StationID       string           `mapstructure:"station_id"`      // 站点唯一标识
	StationName     string           `mapstructure:"station_name"`    // 站点名称
	URL             string           `mapstructure:"url"`             // 站点公开地址（注册到Center时使用）
	Interval        string           `mapstructure:"interval"`        // 同步间隔
	BatchSize       int              `mapstructure:"batch_size"`      // 每批同步数量
	RetryCount      int              `mapstructure:"retry_count"`    // 重试次数
	RetryInterval   string           `mapstructure:"retry_interval"`  // 重试间隔
	HeartbeatInterval string         `mapstructure:"heartbeat_interval"` // 心跳间隔（如 "60s"）
	Proxy           ProxyConfig      `mapstructure:"proxy"`           // 代理配置
	Filter          SyncFilterConfig `mapstructure:"filter"`          // 同步过滤器
}

// SyncFilterConfig 同步过滤器配置
type SyncFilterConfig struct {
	ProjectNames []string `mapstructure:"project_names"` // 只同步这些项目
	StartTime    string   `mapstructure:"start_time"`    // 开始时间
	EndTime      string   `mapstructure:"end_time"`      // 结束时间
}

// ProxyConfig 代理配置
type ProxyConfig struct {
	Enabled  bool   `mapstructure:"enabled"`  // 是否启用代理
	URL      string `mapstructure:"url"`      // 代理URL
	Username string `mapstructure:"username"` // 代理用户名
	Password string `mapstructure:"password"` // 代理密码
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
	LastSyncAt     *time.Time `mapstructure:"-"`            // 上次同步时间（内存，非配置）
	LastSyncCount  int       `mapstructure:"-"`            // 上次同步用户数（内存，非配置）
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

type BaseConfig struct {
	Mode        string `mapstructure:"mode"`
	Port        string `mapstructure:"port"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	Version     string `mapstructure:"version"`
	Language    string `mapstructure:"language"`
	InstallDir  string `mapstructure:"install_dir"`
	ServeWeb    bool   `mapstructure:"serve_web"`
	WebRoot     string `mapstructure:"web_root"` // 前端静态文件目录（可以是绝对路径或相对于配置文件的路径）
	ConfigFile  string `mapstructure:"-"`        // 配置文件路径（非配置项，由启动参数设置）
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
	mu      sync.Mutex
	level   string
}

var AppLogger *Logger

// InitLogger 初始化日志：仅输出到终端
func InitLogger() {
	// 从配置读取（InitGlobal 已在之前调用，CONF 已有值）
	level := "info"
	if CONF != nil {
		level = CONF.Log.Level
	}

	AppLogger = &Logger{
		level: level,
	}
}

// write 向终端输出一条日志
func (l *Logger) write(level, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	ts := time.Now().Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("[%s] [%s] %s", ts, level, msg)
	os.Stdout.WriteString(line + "\n")
}

func (l *Logger) Write(p []byte) (int, error) {
	l.write("INFO", string(p))
	return len(p), nil
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.level != "debug" {
		l.write("INFO", format, v...)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.write("ERROR", format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level == "debug" {
		l.write("DEBUG", format, v...)
	}
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.write("WARN", format, v...)
}

// NewGormLogger 创建干净的 GORM Logger
// 忽略 record not found，只输出真实错误和慢查询（>200ms）
func NewGormLogger(slowThreshold time.Duration) logger.Interface {
	return logger.New(
		log.New(os.Stdout, "[SQL] ", 0),
		logger.Config{
			SlowThreshold:             slowThreshold,
			LogLevel:                  logger.Error,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
}
