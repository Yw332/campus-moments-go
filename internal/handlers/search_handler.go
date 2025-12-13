package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SearchContent 搜索内容
func SearchContent(c *gin.Context) {
	keyword := c.Query("keyword")

	// TODO: 根据关键词搜索
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"moments": []gin.H{
				{"id": 1, "content": "包含 " + keyword + " 的动态"},
			},
			"users": []gin.H{
				{"id": 1, "username": "用户" + keyword},
			},
		},
	})
}

// GetHotWords 获取热词
func GetHotWords(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    []string{"校园", "活动", "学习", "美食", "运动"},
	})
}
