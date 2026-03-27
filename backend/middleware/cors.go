package middleware

import (
	"fmt"
	"net/http"

	"datauptwo/global"

	"github.com/gin-gonic/gin"
)

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowed := global.CONF.CORS.AllowOrigins

		// 检查请求的 Origin 是否在白名单中
		validOrigin := ""
		if origin != "" {
			for _, allowedOrigin := range allowed {
				if allowedOrigin == "*" {
					// 配置文件明确配 * 时退化为反射模式（仅对同源请求有效）
					validOrigin = origin
					break
				}
				if allowedOrigin == origin {
					validOrigin = origin
					break
				}
			}
		}

		if validOrigin != "" {
			c.Header("Access-Control-Allow-Origin", validOrigin)
			c.Header("Vary", "Origin")
		}

		methods := "GET, POST, PUT, DELETE, OPTIONS"
		if len(global.CONF.CORS.AllowMethods) > 0 {
			methods = joinStrings(global.CONF.CORS.AllowMethods, ", ")
		}
		c.Header("Access-Control-Allow-Methods", methods)

		headers := "Origin, Content-Type, Authorization, Accept-Language, X-Submit-Source, X-API-Key, X-Token"
		if len(global.CONF.CORS.AllowHeaders) > 0 {
			headers = joinStrings(global.CONF.CORS.AllowHeaders, ", ")
		}
		c.Header("Access-Control-Allow-Headers", headers)

		if global.CONF.CORS.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", fmt.Sprintf("%d", global.CONF.CORS.MaxAge))
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
