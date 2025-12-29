package handlers

import (
	"github.com/Yw332/campus-moments-go/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetHomePage 获取主页内容（包含公开帖子和好友帖子）
func GetHomePage(c *gin.Context) {
	userID := c.GetString("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	posts, total, err := service.GetHomePagePosts(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"message":  "获取失败: " + err.Error(),
		})
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

// GetPostListEnhanced 增强版帖子列表（包含用户信息）
func GetPostListEnhanced(c *gin.Context) {
	userID := c.GetString("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	postType := c.DefaultQuery("type", "public") // public, friends, all

	posts, total, err := service.GetEnhancedPostList(userID, postType, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"message":  "获取失败: " + err.Error(),
		})
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