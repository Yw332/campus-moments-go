package handlers

import (
	"net/http"
	"strconv"

	"github.com/Yw332/campus-moments-go/internal/service"
	"github.com/gin-gonic/gin"
)

var userService = service.NewUserService()

// GetUserByID 根据ID获取用户信息（公开信息）
func GetUserByID(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户ID不能为空",
			"data":    nil,
		})
		return
	}

	user, err := userService.GetPublicUserInfo(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "用户不存在",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    user,
	})
}

// UpdateAvatar 更新头像
func UpdateAvatar(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未认证",
			"data":    nil,
		})
		return
	}

	uid := userID.(string)

	var req struct {
		AvatarURL  string `json:"avatarUrl" binding:"required"`
		AvatarType int    `json:"avatarType"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	user, err := userService.UpdateAvatar(uid, req.AvatarURL, req.AvatarType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "头像更新成功",
		"data":    user,
	})
}

// UpdateSignature 更新个性签名
func UpdateSignature(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未认证",
			"data":    nil,
		})
		return
	}

	uid := userID.(string)

	var req struct {
		Signature string `json:"signature" binding:"max=200"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	user, err := userService.UpdateSignature(uid, req.Signature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "签名更新成功",
		"data":    user,
	})
}

// UpdateLastActive 更新最后活跃时间
func UpdateLastActive(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未认证",
			"data":    nil,
		})
		return
	}

	uid := userID.(string)

	if err := userService.UpdateLastActive(uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    nil,
	})
}

// SearchUsers 搜索用户
func SearchUsers(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "搜索关键词不能为空",
			"data":    nil,
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	users, total, err := userService.SearchUsers(keyword, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "搜索失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "搜索成功",
		"data": gin.H{
			"users":    users,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// AdminResetUserPassword 管理员重置用户密码
func AdminResetUserPassword(c *gin.Context) {
	targetUserID := c.Param("userId")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户ID不能为空",
			"data":    nil,
		})
		return
	}

	var req struct {
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	if err := userService.ResetUserPassword(targetUserID, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "密码重置成功",
		"data": gin.H{
			"userId": targetUserID,
		},
	})
}

// AdminGetAllUsers 管理员获取所有用户列表
func AdminGetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")

	users, total, err := userService.AdminGetAllUsers(page, pageSize, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取用户列表失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"users":    users,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// AdminBanUser 管理员封禁用户
func AdminBanUser(c *gin.Context) {
	targetUserID := c.Param("userId")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户ID不能为空",
			"data":    nil,
		})
		return
	}

	if err := userService.AdminBanUser(targetUserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "封禁成功",
		"data": gin.H{
			"userId": targetUserID,
		},
	})
}

// AdminUnbanUser 管理员解封用户
func AdminUnbanUser(c *gin.Context) {
	targetUserID := c.Param("userId")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户ID不能为空",
			"data":    nil,
		})
		return
	}

	if err := userService.AdminUnbanUser(targetUserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "解封成功",
		"data": gin.H{
			"userId": targetUserID,
		},
	})
}

// AdminDeleteUser 管理员删除用户
func AdminDeleteUser(c *gin.Context) {
	targetUserID := c.Param("userId")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户ID不能为空",
			"data":    nil,
		})
		return
	}

	if err := userService.AdminDeleteUser(targetUserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
		"data": nil,
	})
}

// AdminGetUserPosts 管理员查看用户动态
func AdminGetUserPosts(c *gin.Context) {
	targetUserID := c.Param("userId")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户ID不能为空",
			"data":    nil,
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// 使用 moment service 获取用户的动态
	// 这里复用现有的 moment service，因为管理员需要能看到所有动态
	var momentService *service.MomentService = service.NewMomentService()
	posts, total, err := momentService.GetUserMoments(targetUserID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取用户动态失败",
			"data":    nil,
		})
		return
	}

	// 转换为响应格式
	convertedPosts := make([]map[string]interface{}, 0, len(posts))
	for _, post := range posts {
		postData := map[string]interface{}{
			"postId":    post.ID,
			"title":     post.Title,
			"content":   post.Content,
			"images":    post.Images,
			"video":     post.Video,
			"tags":      post.Tags,
			"createdAt": post.CreatedAt,
		}

		// 添加作者信息
		if post.User != nil {
			postData["author"] = map[string]interface{}{
				"userId":   post.User.ID,
				"username": post.User.Username,
				"avatarUrl":   post.User.Avatar,
			}
		}

		// 添加统计信息
		postData["likeCount"] = post.LikeCount
		postData["commentCount"] = post.CommentCount
		postData["viewCount"] = post.ViewCount
		postData["visibility"] = post.Visibility

		convertedPosts = append(convertedPosts, postData)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"posts":    convertedPosts,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// AdminGetUserFriends 管理员查看用户好友列表
func AdminGetUserFriends(c *gin.Context) {
	targetUserID := c.Param("userId")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户ID不能为空",
			"data":    nil,
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// 使用 friend service 获取好友列表
	friends, total, err := service.GetFriendList(targetUserID, "", page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取好友列表失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"friends":  friends,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}