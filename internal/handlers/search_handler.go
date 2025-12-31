package handlers

import (
	"net/http"
	"strconv"

	"github.com/Yw332/campus-moments-go/internal/service"
	"github.com/gin-gonic/gin"
)

var searchService = service.NewSearchService()

// SearchContent 搜索内容
func SearchContent(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "搜索关键词不能为空",
			"data":    nil,
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 执行搜索
	results, err := searchService.SearchContent(keyword, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "搜索失败",
			"data":    nil,
		})
		return
	}

	// 转换 posts 为响应格式（id -> postId, user -> author）
	convertedPosts := make([]map[string]interface{}, 0, len(results.Posts))
	for _, post := range results.Posts {
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
		"code":    http.StatusOK,
		"message": "success",
		"data": map[string]interface{}{
			"posts":      convertedPosts,
			"users":      results.Users,
			"pagination": results.Pagination,
		},
	})
}

// GetHotWords 获取热词
func GetHotWords(c *gin.Context) {
	hotWords, err := searchService.GetHotWords()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "获取热词失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    hotWords,
	})
}

// GetSearchHistory 获取搜索历史
func GetSearchHistory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "未认证",
			"data":    nil,
		})
		return
	}

	uid := userID.(string)
	history, err := searchService.GetSearchHistory(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "获取搜索历史失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    history,
	})
}

// SaveSearchHistory 保存搜索历史
func SaveSearchHistory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "未认证",
			"data":    nil,
		})
		return
	}

	var req struct {
		Keyword string `json:"keyword" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	uid := userID.(string)
	err := searchService.SaveSearchHistory(uid, req.Keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "保存搜索历史失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "保存成功",
		"data":    nil,
	})
}

// GetFilteredContent 获取筛选内容
func GetFilteredContent(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 获取筛选参数
	visibility := c.Query("visibility")
	tags := c.QueryArray("tags")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	filter := service.FilterRequest{
		Page:       page,
		PageSize:   pageSize,
		Visibility: visibility,
		Tags:       tags,
		StartDate:  startDate,
		EndDate:    endDate,
	}

	results, err := searchService.GetFilteredContent(&filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "获取筛选内容失败",
			"data":    nil,
		})
		return
	}

	// 转换 posts 为响应格式（id -> postId, user -> author）
	convertedPosts := make([]map[string]interface{}, 0, len(results.Posts))
	for _, post := range results.Posts {
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
		"code":    http.StatusOK,
		"message": "success",
		"data": map[string]interface{}{
			"posts":      convertedPosts,
			"users":      results.Users,
			"pagination": results.Pagination,
		},
	})
}

// GetSearchSuggestions 获取搜索建议
func GetSearchSuggestions(c *gin.Context) {
	keyword := c.Query("keyword")
	
	suggestions, err := searchService.GetSearchSuggestions(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "获取搜索建议失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    suggestions,
	})
}