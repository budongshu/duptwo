package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"datauptwo/app/dto"
	"datauptwo/app/model"
	"datauptwo/app/repo"
	"datauptwo/global"

	"github.com/gin-gonic/gin"
)

// stationCache API Key → Station 的内存缓存，TTL 10 分钟自动刷新
// 采用 Cache-Aside 模式：缓存未命中时直接查 DB，保证数据一致性
type stationCache struct {
	mu      sync.RWMutex
	byHash  map[string]model.SyncStation // key: sha256(apiKey)
	loadedAt time.Time
	ttl     time.Duration
}

var apiKeyCache = stationCache{ttl: 10 * time.Minute}

// loadCache 刷新缓存（从 DB 加载所有活跃站点）
func (c *stationCache) loadCache() {
	stationRepo := repo.NewSyncStationRepo()
	stations, err := stationRepo.GetAll()
	if err != nil {
		global.AppLogger.Error("[ApiKeyCache] 刷新站点缓存失败: %v", err)
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	// DB 存的 api_key 是明文的 hash，直接用 hash 作为 map 的 key
	c.byHash = make(map[string]model.SyncStation, len(stations))
	for _, s := range stations {
		// s.APIKey 存储的是明文的 SHA256 hash，直接作为缓存 key
		c.byHash[s.APIKey] = s
	}
	c.loadedAt = time.Now()
	global.AppLogger.Info("[ApiKeyCache] 缓存已刷新，共 %d 个活跃站点", len(stations))
}

// get 返回缓存中对应的站点
// Cache-Aside 模式：缓存命中直接返回；缓存未命中则查 DB，保证一致性
func (c *stationCache) get(apiKey string) (model.SyncStation, bool) {
	hash := sha256.Sum256([]byte(apiKey))
	hashStr := hex.EncodeToString(hash[:])

	// 1. 先查缓存（TTL 内直接返回）
	c.mu.RLock()
	s, ok := c.byHash[hashStr]
	c.mu.RUnlock()

	if ok {
		global.AppLogger.Debug("[ApiKeyCache] 缓存命中，站点=%s [%s]", s.Name, s.Code)
		return s, true
	}

	// 2. 缓存未命中 → 直接查 DB（hash 对 hash 比对）
	global.AppLogger.Info("[ApiKeyCache] 缓存未命中 key=%s，查询 DB...", hashStr[:16]+"...")
	stationRepo := repo.NewSyncStationRepo()
	station, err := stationRepo.GetByAPIKey(hashStr) // DB 存的也是 hash，直接比对
	if err != nil {
		global.AppLogger.Warn("[ApiKeyCache] DB 查询失败或 Key 不存在: %v", err)
		return model.SyncStation{}, false
	}

	if station.Status != "active" {
		global.AppLogger.Warn("[ApiKeyCache] 站点已禁用: %s [%s]", station.Name, station.Code)
		return model.SyncStation{}, false
	}

	// 3. 异步写回缓存，下次请求直接命中
	go func() {
		c.mu.Lock()
		c.byHash[hashStr] = *station
		c.mu.Unlock()
		global.AppLogger.Info("[ApiKeyCache] Key 已缓存: %s [%s]", station.Name, station.Code)
	}()

	global.AppLogger.Info("[ApiKeyCache] DB 验证成功，站点=%s [%s]（缓存未命中，已回写）", station.Name, station.Code)
	return *station, true
}

// ApiKeyAuth API Key认证中间件
// 验证 X-API-Key 头，用于站点间API调用认证
func ApiKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			global.AppLogger.Warn("[ApiKeyAuth] 缺少 API Key，请求路径: %s", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, dto.Response{
				Code:    401,
				Message: "缺少 API Key",
			})
			c.Abort()
			return
		}

		// 验证 API Key（缓存优先，未命中则查 DB）
		station, ok := apiKeyCache.get(apiKey)
		if !ok {
			global.AppLogger.Warn("[ApiKeyAuth] API Key 验证失败: 无效或站点已禁用")
			c.JSON(http.StatusUnauthorized, dto.Response{
				Code:    401,
				Message: "无效的 API Key",
			})
			c.Abort()
			return
		}

		// 将站点信息存储到上下文
		c.Set("stationId", station.ID)
		c.Set("stationCode", station.Code)
		c.Set("stationName", station.Name)

		global.AppLogger.Debug("[ApiKeyAuth] 站点 %s [%s] 认证成功", station.Name, station.Code)
		c.Next()
	}
}

// GetStationID 获取当前站点ID
func GetStationID(c *gin.Context) uint {
	if id, exists := c.Get("stationId"); exists {
		return id.(uint)
	}
	return 0
}

// GetStationCode 获取当前站点代码
func GetStationCode(c *gin.Context) string {
	if code, exists := c.Get("stationCode"); exists {
		return code.(string)
	}
	return ""
}

// GetStationName 获取当前站点名称
func GetStationName(c *gin.Context) string {
	if name, exists := c.Get("stationName"); exists {
		return name.(string)
	}
	return ""
}

// InitAPIKeyCache 启动时预热缓存（Center 模式调用一次即可）
func InitAPIKeyCache() {
	apiKeyCache.loadCache()
}

// RefreshCache 主动刷新缓存（预热用，实际不再依赖）
func RefreshCache() {
	apiKeyCache.loadCache()
}
