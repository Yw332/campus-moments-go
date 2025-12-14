package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdateUserProfile 更新用户资料（占位实现）
func UpdateUserProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "UpdateUserProfile 未实现（占位）",
		"data":    nil,
	})
}