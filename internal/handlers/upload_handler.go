package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// UploadFile 上传文件
func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择文件"})
		return
	}

	// 生成文件名
	filename := filepath.Base(file.Filename)
	// TODO: 保存文件到服务器或云存储

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "上传成功",
		"data": gin.H{
			"url": "https://example.com/uploads/" + filename,
		},
	})
}

// UploadAvatar 上传头像
func UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择头像文件"})
		return
	}

	// 生成文件名（示例）并返回 URL
	filename := filepath.Base(file.Filename)
	// TODO: 验证文件类型、大小，保存文件到存储
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "头像上传成功",
		"data": gin.H{
			"avatarUrl": "https://example.com/avatars/" + filename,
		},
	})
}
