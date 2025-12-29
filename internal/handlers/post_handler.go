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
		Title      string            `json:"title" binding:"max=100"`
		Content    string            `json:"content" binding:"required,min=1,max=10000"`
		Images     []string          `json:"images"`
		Video      string            `json:"video"`
		Visibility int               `json:"visibility" binding:"oneof=0 1 2"`
		Tags       []string          `json:"tags"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}
	
	userID := c.GetString("userID")
	post, err := service.CreatePost(userID, req.Title, req.Content, req.Images, req.Video, req.Visibility, req.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"code": http.StatusCreated,
		"message":  "创建成功",
		"data": post,
	})
}

// GetPostList 获取帖子列表
func GetPostList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	visibility := c.DefaultQuery("visibility", "0") // 默认只看公开帖子
	
	userID := c.GetString("userID")
	
	posts, total, err := service.GetEnhancedPostList(userID, visibility, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "获取成功",
		"data": gin.H{
			"posts": posts,
			"total": total,
			"page":  page,
			"pageSize": pageSize,
		},
	})
}

// GetPostDetail 获取帖子详情
func GetPostDetail(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的帖子ID"})
		return
	}
	
	userID := c.GetString("userID")
	post, err := service.GetEnhancedPostDetail(postID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
		return
	}
	
	// 增加浏览量
	service.IncrementViewCount(postID)
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "获取成功",
		"data": post,
	})
}

// UpdatePost 更新帖子
func UpdatePost(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的帖子ID"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}
	
	userID := c.GetString("userID")
	post, err := service.UpdatePost(postID, userID, req.Title, req.Content, req.Images, req.Video, req.Visibility, req.Tags)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "更新失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "更新成功",
		"data": post,
	})
}

// DeletePost 删除帖子
func DeletePost(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的帖子ID"})
		return
	}
	
	userID := c.GetString("userID")
	err = service.DeletePost(postID, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "删除失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "删除成功",
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "获取成功",
		"data": gin.H{
			"posts": posts,
			"total": total,
			"page":  page,
			"pageSize": pageSize,
		},
	})
}