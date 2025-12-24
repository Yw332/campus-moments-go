package routes

import (
	"github.com/Yw332/campus-moments-go/internal/handlers"
	"github.com/Yw332/campus-moments-go/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// 添加CORS中间件到所有路由（在路由设置之前）
	router.Use(middleware.CORSMiddleware())
	
	// 全局OPTIONS处理（解决CORS预检问题）
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.AbortWithStatus(204)
	})
	
	// ========== 公共路由（无需认证）==========
	router.GET("/", handlers.Home)
	router.GET("/health", handlers.HealthCheck)

	// 认证相关
	auth := router.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.POST("/send-verification", handlers.SendVerificationCode)
		auth.POST("/verify-and-reset", handlers.VerifyAndResetPassword)
	}

	// 公开的动态列表（不需要登录也能看）
	router.GET("/moments", handlers.GetMoments)
	
	// 公开的搜索功能（不需要登录）
	router.GET("/search/hot-words", handlers.GetHotWords)

	// 管理端登录
	router.POST("/admin/login", handlers.AdminLogin)

	// 管理端路由（需要登录 + 管理员权限）
	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		admin.GET("/menu", handlers.AdminMenu)
		admin.GET("/users", handlers.ListUsers)
		admin.GET("/users/:id", handlers.GetUserDetail)
		admin.POST("/users/:id/password", handlers.ResetUserPassword)
		admin.DELETE("/users/:id", handlers.DeleteUser)
		admin.GET("/contents", handlers.ListContents)
		admin.GET("/contents/:id", handlers.GetContentDetail)
		admin.DELETE("/contents/:id", handlers.DeleteContent)
	}

	// ========== 需要认证的路由 ==========
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// 认证相关
		api.POST("/auth/logout", handlers.Logout)

		// 上传相关
		upload := api.Group("/upload")
		{
			upload.POST("/file", handlers.UploadFile)
			upload.POST("/avatar", handlers.UploadAvatar)
		}

		// 动态相关
		moments := api.Group("/moments")
		{
			moments.POST("", handlers.CreateMoment)
			moments.GET("/:id", handlers.GetMomentDetail)
			moments.PATCH("/:id", handlers.UpdateMoment)
			moments.DELETE("/:id", handlers.DeleteMoment)
			moments.GET("/my", handlers.GetUserMoments) // 获取当前用户的动态列表
			moments.GET("/:id/comments", handlers.GetMomentComments) // 获取动态评论
		}

		// 互动相关
		api.POST("/comments", handlers.CreateComment)
		api.DELETE("/comments/:id", handlers.DeleteComment)
		api.POST("/likes", handlers.ToggleLike)

		// 用户相关
		users := api.Group("/users")
		{
			users.GET("/profile", handlers.GetProfile)
			users.PUT("/profile", handlers.UpdateUserProfile)
			users.PUT("/password", handlers.ChangePassword)
		}

		// 搜索相关
		search := api.Group("/search")
		{
			search.GET("/hot-words", handlers.GetHotWords)
			search.GET("/history", handlers.GetSearchHistory)
			search.GET("/filter", handlers.GetFilteredContent)
			search.POST("/history", handlers.SaveSearchHistory)
			search.GET("/suggestions", handlers.GetSearchSuggestions) // 搜索建议
			search.GET("", handlers.SearchContent)                    // GET /search?keyword=xxx
		}
	}
}
