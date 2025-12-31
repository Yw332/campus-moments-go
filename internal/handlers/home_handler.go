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
			"likeCount": post.LikeCount,
			"commentCount": post.CommentCount,
			"viewCount": post.ViewCount,
			"visibility": post.Visibility,
		}

		// 添加作者信息
		if post.User != nil {
			postData["author"] = map[string]interface{}{
				"userId":   post.User.ID,
				"username": post.User.Username,
				"avatarUrl":   post.User.Avatar,
			}
		}

		convertedPosts = append(convertedPosts, postData)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "获取成功",
		"data": gin.H{
			"posts": convertedPosts,
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

	// 转换为响应格式（id -> postId, username/avatar -> author）
	convertedPosts := make([]map[string]interface{}, 0, len(posts))
	for _, post := range posts {
		postData := map[string]interface{}{
			"postId":       post.ID,
			"title":        post.Title,
			"content":      post.Content,
			"images":       post.Images,
			"video":        post.Video,
			"createdAt":    post.CreatedAt,
			"likeCount":    post.LikeCount,
			"commentCount": post.CommentCount,
			"viewCount":    post.ViewCount,
		}

		// 添加作者信息
		postData["author"] = map[string]interface{}{
			"userId":   post.AuthorID,
			"username": post.Username,
			"avatarUrl":   post.Avatar,
		}

		convertedPosts = append(convertedPosts, postData)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "获取成功",
		"data": gin.H{
			"posts": convertedPosts,
			"total": total,
			"page":  page,
			"pageSize": pageSize,
		},
	})
}