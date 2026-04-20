package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"

	"datauptwo/app/dto"
	"datauptwo/app/model"
	"datauptwo/app/repo"

	"github.com/gin-gonic/gin"
)

// ApiKeyAuth API Key认证中间件
// 验证 X-API-Key 头，用于站点间API调用认证
func ApiKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, dto.Response{
				Code:    401,
				Message: "缺少 API Key",
			})
			c.Abort()
			return
		}

		// 验证 API Key
		station, err := validateAPIKey(apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.Response{
				Code:    401,
				Message: "无效的 API Key",
			})
			c.Abort()
			return
		}

		// 检查站点状态
		if station.Status != "active" {
			c.JSON(http.StatusForbidden, dto.Response{
				Code:    403,
				Message: "站点已被禁用",
			})
			c.Abort()
			return
		}

		// 将站点信息存储到上下文
		c.Set("stationId", station.ID)
		c.Set("stationCode", station.Code)
		c.Set("stationName", station.Name)

		c.Next()
	}
}

// validateAPIKey 验证API Key（内部函数，避免循环导入）
func validateAPIKey(apiKey string) (*model.SyncStation, error) {
	// 查找所有活跃站点，逐一验证
	stationRepo := repo.NewSyncStationRepo()
	stations, err := stationRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// 计算输入的哈希值
	hash := sha256.Sum256([]byte(apiKey))
	inputHash := hex.EncodeToString(hash[:])

	for _, station := range stations {
		if station.APIKey == inputHash {
			return &station, nil
		}
	}

	return nil, errors.New("invalid API key")
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
