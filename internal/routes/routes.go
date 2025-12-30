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
		auth.POST("/send-verification", handlers.SendVerificationCode)
		auth.POST("/verify-and-reset", handlers.VerifyAndResetPassword)
	}

	// 公开的帖子列表（不需要登录也能看）
	router.GET("/public/posts", handlers.GetPostList)
	router.GET("/public/posts/:id", handlers.GetPostDetail)
	
	// 主页帖子列表（支持公开和好友帖子）
	router.GET("/home", handlers.GetHomePage)
	
	// 获取公开评论列表
	router.GET("/public/posts/:id/comments", func(c *gin.Context) {
		postID := c.Param("id")
		c.Set("postId", postID)
		handlers.GetCommentList(c)
	})
	
	// 获取标签列表和热门标签
	router.GET("/public/tags", handlers.GetTagList)
	router.GET("/public/tags/hot", handlers.GetHotTags)
	router.GET("/public/tags/search", handlers.SearchTags)
	router.GET("/public/tags/by-name/:name/posts", handlers.GetTagPosts)
	router.GET("/public/tags/by-id/:id", func(c *gin.Context) {
		// 将参数名从 id 转换为标准格式
		c.Set("tagId", c.Param("id"))
		handlers.GetTagDetail(c)
	})

	// 公开搜索功能（无需认证）
	router.GET("/search", handlers.SearchContent)
	router.GET("/search/hot-words", handlers.GetHotWords)
	router.GET("/search/suggestions", handlers.GetSearchSuggestions)

		// ========== 需要认证的路由 ==========
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// ========== 管理员专用路由 ==========
		admin := api.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			// 管理员用户管理
			admin.GET("/users", handlers.AdminGetAllUsers)
			admin.GET("/users/:userId/posts", handlers.AdminGetUserPosts)
			admin.GET("/users/:userId/friends", handlers.AdminGetUserFriends)
			admin.PUT("/users/:userId/password", handlers.AdminResetUserPassword)
			admin.PUT("/users/:userId/ban", handlers.AdminBanUser)
			admin.PUT("/users/:userId/unban", handlers.AdminUnbanUser)
			admin.DELETE("/users/:userId", handlers.AdminDeleteUser)
			// 管理员删除帖子
			admin.DELETE("/posts/:id", handlers.AdminDeleteMoment)
		}

		// 认证相关
		api.POST("/auth/logout", handlers.Logout)

		// 上传相关
		upload := api.Group("/upload")
		{
			upload.POST("/file", handlers.UploadFile)
			upload.POST("/avatar", handlers.UploadAvatar)
		}

		// ========== 帖子相关 ==========
		posts := api.Group("/posts")
		{
			posts.POST("", handlers.CreatePost)
			posts.PUT("/:id", handlers.UpdatePost)
			posts.DELETE("/:id", handlers.DeletePost)
			posts.GET("/my", handlers.GetUserPosts)
			posts.GET("/user/:userId", handlers.GetUserPosts)
			posts.GET("/:id/likes", handlers.GetPostLikes)
		}

		// ========== 评论相关 ==========
		comments := api.Group("/comments")
		{
			comments.POST("/post/:postId", handlers.CreateComment)
			comments.PUT("/:id", handlers.UpdateComment)
			comments.DELETE("/:id", handlers.DeleteComment)
			comments.POST("/:id/like", handlers.LikeComment)
			comments.POST("/:id/reply", handlers.ReplyComment)
			comments.GET("/:id/likes", handlers.GetCommentLikes)
		}

		// ========== 点赞相关 ==========
		likes := api.Group("/likes")
		{
			likes.POST("/post/:postId", handlers.LikePost)
			likes.GET("/posts/:postId", handlers.GetPostLikes)
			likes.GET("/comments/:commentId", handlers.GetCommentLikes)
			likes.GET("/users/:userId", handlers.GetUserLikes)
		}

		// ========== 好友相关 ==========
		friends := api.Group("/friends")
		{
			friends.POST("/request", handlers.SendFriendRequest)
			friends.GET("/requests", handlers.GetFriendRequests)
			friends.PUT("/requests/:id", handlers.HandleFriendRequest)
			friends.GET("", handlers.GetFriendList)
			friends.DELETE("/:friendId", handlers.DeleteFriend)
			friends.PUT("/:friendId/remark", handlers.UpdateFriendRemark)
			friends.GET("/search", handlers.SearchFriends)
		}

		// ========== 消息相关 ==========
		messages := api.Group("/messages")
		{
			messages.POST("", handlers.SendMessage)
			messages.GET("/:peerId", handlers.GetMessageList)
			messages.PUT("/:peerId/read", handlers.MarkMessagesAsRead)
		}

		// ========== 会话相关 ==========
		conversations := api.Group("/conversations")
		{
			conversations.GET("", handlers.GetConversationList)
			conversations.PUT("/:peerId/pin", handlers.PinConversation)
			conversations.DELETE("/:peerId/pin", handlers.UnpinConversation)
			conversations.PUT("/:peerId/mute", handlers.MuteConversation)
			conversations.DELETE("/:peerId/mute", handlers.UnmuteConversation)
			conversations.DELETE("/:peerId", handlers.DeleteConversation)
			conversations.GET("/unread", handlers.GetUnreadCount)
		}

		// ========== 标签相关（管理功能） ==========
		tags := api.Group("/tags")
		{
			tags.POST("", handlers.CreateTag)
			tags.PUT("/:id", handlers.UpdateTag)
			tags.DELETE("/:id", handlers.DeleteTag)
		}

		// ========== 用户相关 ==========
		users := api.Group("/users")
		{
			users.GET("/profile", handlers.GetProfile)
			users.PUT("/profile", handlers.UpdateUserProfile)
			users.PUT("/password", handlers.ChangePassword)
			users.PUT("/avatar", handlers.UpdateAvatar)
			users.PUT("/signature", handlers.UpdateSignature)
			users.POST("/active", handlers.UpdateLastActive)
			users.GET("/:userId", handlers.GetUserByID)
			users.GET("/search", handlers.SearchUsers)
		}

		// ========== 搜索相关 ==========
		search := api.Group("/search")
		{
			search.GET("/hot-words", handlers.GetHotWords)
			search.GET("/history", handlers.GetSearchHistory)
			search.GET("/filter", handlers.GetFilteredContent)
			search.POST("/history", handlers.SaveSearchHistory)
			search.GET("/suggestions", handlers.GetSearchSuggestions)
			search.GET("", handlers.SearchContent)
		}

		// ========== 兼容旧的动态接口 ==========
		moments := api.Group("/moments")
		{
			moments.POST("", handlers.CreateMoment)
			moments.GET("/:id", handlers.GetMomentDetail)
			moments.PATCH("/:id", handlers.UpdateMoment)
			moments.DELETE("/:id", handlers.DeleteMoment)
			moments.GET("/my", handlers.GetUserMoments)
		}
	}
}
