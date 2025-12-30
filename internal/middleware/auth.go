package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Yw332/campus-moments-go/pkg/database"
	"github.com/Yw332/campus-moments-go/pkg/jwt"
	"github.com/Yw332/campus-moments-go/pkg/token_blacklist"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "missing authorization header",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(auth, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "invalid authorization header format",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 提取token
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "token cannot be empty",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 检查Token是否在黑名单中（退出登录的Token）
		blacklist := token_blacklist.GetInstance()
		if blacklist.IsBlacklisted(token) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "token has been revoked",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 验证token
		claims, err := jwt.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "invalid or expired token",
				"data":    nil,
			})
			c.Abort()
			return
		}

	// 将用户信息设置到上下文（转换为字符串ID）
		c.Set("userID", fmt.Sprintf("%010d", claims.UserID))
		c.Set("username", claims.Username)
		c.Set("claims", claims)

		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件（用户信息可选）
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth != "" && strings.HasPrefix(auth, "Bearer ") {
			token := strings.TrimPrefix(auth, "Bearer ")
			if claims, err := jwt.ParseToken(token); err == nil {
				c.Set("userID", claims.UserID)
				c.Set("username", claims.Username)
				c.Set("claims", claims)
			}
		}
		c.Next()
	}
}

// AdminMiddleware 管理员认证中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户信息
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未认证",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 检查是否为管理员（负数ID表示来自 admins 表）
		jwtClaims, ok := claims.(*jwt.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的凭证",
				"data":    nil,
			})
			c.Abort()
			return
		}

		if jwtClaims.UserID < 0 {
			// 管理员表登录，直接放行
			c.Next()
			return
		}

		// 普通用户，检查 users 表的 role 字段
		db := database.GetDB()
		userID := fmt.Sprintf("%010d", jwtClaims.UserID)

		var user struct {
			Role int `json:"role"`
		}
		if err := db.Table("users").Select("role").Where("id = ?", userID).First(&user).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "需要管理员权限",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 检查是否为管理员
		if user.Role != 1 {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "需要管理员权限",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
