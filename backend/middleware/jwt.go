package middleware

import (
	"crypto/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"datauptwo/global"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret caches the resolved JWT secret to avoid regenerating it on every call.
var (
	jwtSecret     []byte
	jwtSecretOnce  sync.Once
)

// getJWTSecret returns the JWT signing key, caching it on first call.
// Falls back to a random 32-byte secret only when jwt.secret is truly unset.
func getJWTSecret() []byte {
	jwtSecretOnce.Do(func() {
		secret := global.CONF.JWT.Secret
		if secret != "" {
			jwtSecret = []byte(secret)
			return
		}
		// Generate a cryptographically secure random secret when not configured
		b := make([]byte, 32)
		if _, err := rand.Read(b); err == nil {
			jwtSecret = b
		} else {
			jwtSecret = []byte("insecure-fallback-secret-do-not-use-in-production")
		}
	})
	return jwtSecret
}

// Claims JWT claims
type Claims struct {
	UserID      uint     `json:"userId"`
	Username    string   `json:"username"`
	IsAdmin     bool     `json:"isAdmin,omitempty"`   // 系统内置管理员标识
	Permissions []string  `json:"permissions,omitempty"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(userID uint, username string, isAdmin bool, permissions ...string) (string, error) {
	claims := Claims{
		UserID:      userID,
		Username:    username,
		IsAdmin:     isAdmin,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)), // 7天过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "datauptwo",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

// GenerateTmpToken 生成临时token（用于MFA验证）
func GenerateTmpToken(userID uint, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)), // 5分钟过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "datauptwo-tmp",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// GetToken 从请求中获取token
func GetToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		// 尝试从 query parameter 获取
		token := c.Query("token")
		if token != "" {
			return token
		}
		return ""
	}

	// Bearer token
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := GetToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未登录或登录已过期",
			})
			c.Abort()
			return
		}

		claims, err := ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "登录已过期，请重新登录",
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("isAdmin", claims.IsAdmin)
		c.Set("permissions", claims.Permissions)

		c.Next()
	}
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) uint {
	userID, exists := c.Get("userId")
	if !exists {
		return 0
	}
	if id, ok := userID.(uint); ok {
		return id
	}
	return 0
}

// GetUsername 从上下文获取用户名
func GetUsername(c *gin.Context) string {
	username, exists := c.Get("username")
	if !exists {
		return ""
	}
	if name, ok := username.(string); ok {
		return name
	}
	return ""
}

// GetPermissions 从上下文获取权限列表
func GetPermissions(c *gin.Context) []string {
	perms, exists := c.Get("permissions")
	if !exists {
		return nil
	}
	if p, ok := perms.([]string); ok {
		return p
	}
	return nil
}

// HasPermission 检查是否拥有指定权限
func HasPermission(c *gin.Context, perm string) bool {
	// 系统内置管理员拥有所有权限
	isAdmin, _ := c.Get("isAdmin")
	if isAdmin.(bool) {
		return true
	}
	perms, exists := c.Get("permissions")
	if !exists {
		return false
	}
	for _, p := range perms.([]string) {
		if p == "admin:all" || p == perm {
			return true
		}
	}
	return false
}

// RequirePermission 权限校验中间件
func RequirePermission(perm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !HasPermission(c, perm) {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "无此权限：" + perm,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
