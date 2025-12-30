package handlers

import (
	"github.com/Yw332/campus-moments-go/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CreatePost 创建帖子
func CreatePost(c *gin.Context) {
	var req struct {
		Title      string   `json:"title" binding:"max=100"`
		Content    string   `json:"content" binding:"required,min=1,max=10000"`
		Images     []string `json:"images"`
		Video      string   `json:"video"`
		Visibility int      `json:"visibility" binding:"oneof=0 1 2"`
		Tags       []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	userID := c.GetString("userID")
	post, err := service.CreatePost(userID, req.Title, req.Content, req.Images, req.Video, req.Visibility, req.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "创建成功",
		"data":    post,
	})
}

// GetPostList 获取帖子列表
func GetPostList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	visibility := c.DefaultQuery("visibility", "0")

	userID := c.GetString("userID")

	posts, total, err := service.GetPostList(userID, visibility, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 转换为响应格式（id -> postId, user -> author）
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
				"avatar":   post.User.Avatar,
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

// GetPostDetail 获取帖子详情
func GetPostDetail(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的帖子ID",
			"data":    nil,
		})
		return
	}

	userID := c.GetString("userID")
	post, err := service.GetPostDetail(postID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "帖子不存在",
			"data":    nil,
		})
		return
	}

	// 增加浏览量
	service.IncrementViewCount(postID)

	// 转换为响应格式
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
			"avatar":   post.User.Avatar,
		}
	}

	// 添加统计信息
	postData["likeCount"] = post.LikeCount
	postData["commentCount"] = post.CommentCount
	postData["viewCount"] = post.ViewCount
	postData["visibility"] = post.Visibility

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    postData,
	})
}

// UpdatePost 更新帖子
func UpdatePost(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的帖子ID",
			"data":    nil,
		})
		return
	}

	var req struct {
		Title      string   `json:"title" binding:"max=100"`
		Content    string   `json:"content" binding:"required,min=1,max=10000"`
		Images     []string `json:"images"`
		Video      string   `json:"video"`
		Visibility int      `json:"visibility" binding:"oneof=0 1 2"`
		Tags       []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	userID := c.GetString("userID")
	post, err := service.UpdatePost(postID, userID, req.Title, req.Content, req.Images, req.Video, req.Visibility, req.Tags)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "更新失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    post,
	})
}

// DeletePost 删除帖子
func DeletePost(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的帖子ID",
			"data":    nil,
		})
		return
	}

	userID := c.GetString("userID")
	err = service.DeletePost(postID, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "删除失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
		"data":    nil,
	})
}

// GetUserPosts 获取用户帖子列表
func GetUserPosts(c *gin.Context) {
	targetUserID := c.Param("userId")
	if targetUserID == "" {
		targetUserID = c.GetString("userID")
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	userID := c.GetString("userID")
	posts, total, err := service.GetUserPosts(userID, targetUserID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取失败: " + err.Error(),
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
				"avatar":   post.User.Avatar,
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