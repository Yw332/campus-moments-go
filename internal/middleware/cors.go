package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware 跨域资源共享中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求来源
		origin := c.GetHeader("Origin")
		
		// 允许的域名列表 - 生产环境应该限制具体的域名
		allowedOrigins := []string{
			"http://localhost:3000",        // React开发服务器
			"http://localhost:5173",        // Vite开发服务器
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
			"http://localhost:8080",       // 如果前端也在8080
			"http://localhost",            // 本地host
			"http://127.0.0.1",
			"https://106.52.165.122:3000", // 服务器上的前端域名
			"https://106.52.165.122:5173",
			"https://106.52.165.122:8080",
			"https://106.52.165.122",      // 服务器主域名
			"*",                          // 允许所有域名（开发时使用，生产环境建议移除）
		}
		
		// 检查是否在允许的域名列表中
		isAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				isAllowed = true
				break
			}
		}
		
		// 如果是允许的域名或者请求头中没有Origin，设置CORS头
		if origin == "" || isAllowed {
			if origin != "" {
				c.Header("Access-Control-Allow-Origin", origin)
			} else {
				// 如果没有Origin头，允许所有（通常是服务器端请求）
				c.Header("Access-Control-Allow-Origin", "*")
			}
		} else {
			// 如果不在允许列表中，仍然返回允许的头，但是只返回预检请求
			c.Header("Access-Control-Allow-Origin", "null")
		}
		
		// 设置允许的方法
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
		
		// 设置允许的头
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type, Authorization, X-Requested-With, Accept, Cache-Control, X-File-Name")
		
		// 设置暴露的头（让前端可以访问这些头）
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type, Authorization, X-Total-Count")
		
		// 设置是否允许携带凭证
		c.Header("Access-Control-Allow-Credentials", "true")
		
		// 设置预检请求的缓存时间（秒）
		c.Header("Access-Control-Max-Age", "86400")
		
		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

// CORSMiddlewareForProduction 生产环境CORS中间件（更严格的配置）
func CORSMiddlewareForProduction() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		
		// 生产环境只允许特定的域名
		allowedOrigins := []string{
			"https://106.52.165.122",        // 主域名
			"https://106.52.165.122:3000",   // 前端端口3000
			"https://106.52.165.122:5173",   // Vite端口
			"https://yourdomain.com",        // 你的正式域名
		}
		
		// 检查域名是否允许
		isAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			if strings.HasPrefix(origin, allowedOrigin) {
				isAllowed = true
				break
			}
		}
		
		if isAllowed && origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		} else if origin == "" {
			// 服务器端请求
			c.Header("Access-Control-Allow-Origin", "*")
		}
		
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type, Authorization, X-Requested-With, Accept, Cache-Control, X-File-Name")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type, Authorization, X-Total-Count")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}