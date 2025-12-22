package middleware

import (
	"net/http"

	"github.com/Yw332/campus-moments-go/internal/service"
	"github.com/gin-gonic/gin"
)

// AdminMiddleware 确保用户为管理员
func AdminMiddleware() gin.HandlerFunc {
	adminService := service.NewAdminService()
	return func(c *gin.Context) {
		username, ok := c.Get("username")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未认证", "data": nil})
			c.Abort()
			return
		}

		isAdmin, _ := adminService.IsAdminByUsername(username.(string))
		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "需要管理员权限", "data": nil})
			c.Abort()
			return
		}

		c.Next()
	}
}
