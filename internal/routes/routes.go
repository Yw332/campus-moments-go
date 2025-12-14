package routes

import (
	"github.com/Yw332/campus-moments-go/internal/handlers"
	"github.com/Yw332/campus-moments-go/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// ========== 公共路由（无需认证）==========
	router.GET("/", handlers.Home)
	router.GET("/health", handlers.HealthCheck)

	// 认证相关
	auth := router.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// 公开的动态列表（不需要登录也能看）
	router.GET("/moments", handlers.GetMoments)

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
		}

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
			search.GET("", handlers.SearchContent) // GET /search?keyword=xxx
		}
	}
}
