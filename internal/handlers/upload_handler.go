package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Yw332/campus-moments-go/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请选择头像文件",
		})
		return
	}

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !allowedExtensions[ext] {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "头像文件格式不支持，支持格式：jpg, jpeg, png, gif, webp",
		})
		return
	}

	// 验证文件大小 (最大5MB)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "头像文件大小不能超过5MB",
		})
		return
	}

	// 生成唯一的文件名
	uniqueID := uuid.New().String()
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s_%s%s", userID, timestamp, uniqueID, ext)

	// 创建上传目录
	uploadDir := "uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建上传目录失败",
		})
		return
	}

	// 保存文件
	filePath := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "保存头像文件失败",
		})
		return
	}

	// 生成头像URL
	avatarURL := fmt.Sprintf("http://106.52.165.122:8080/uploads/avatars/%s", filename)

	// 获取数据库连接
	db := models.GetDB()
	
	// 更新用户头像
	if err := db.Model(&models.User{}).Where("id = ?", userID).Update("avatar", avatarURL).Error; err != nil {
		// 如果数据库更新失败，删除已上传的文件
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新用户头像失败",
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "头像上传成功",
		"data": gin.H{
			"avatarUrl": avatarURL,
			"userId":    userID,
		},
	})
}
