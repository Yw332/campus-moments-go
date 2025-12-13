package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 是一个占位的认证中间件。
// 目前仅检查 Authorization 头为 `Bearer <token>` 格式，未来可替换为 JWT 验证。
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}
		if !strings.HasPrefix(auth, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		_ = token // TODO: 验证 token 并设置用户信息到上下文

		// 示例：将 userID 放到上下文，方便 handlers 使用（后续从 token 解出真实 userID）
		c.Set("userID", int64(1))
		c.Next()
	}
}
