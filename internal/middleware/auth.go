package middleware

import (
	"fmt"
	"net/http"
	"strings"

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
