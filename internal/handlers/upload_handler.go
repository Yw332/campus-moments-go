package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadFile 上传文件
func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请选择文件",
			"data":    nil,
		})
		return
	}

	// 1. 验证文件大小 (限制10MB)
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文件大小不能超过10MB",
			"data":    nil,
		})
		return
	}

	// 2. 验证文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".mp4":  true,
		".mov":  true,
		".avi":  true,
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".txt":  true,
		".zip":  true,
		".rar":  true,
	}
	
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "不支持的文件类型",
			"data":    nil,
		})
		return
	}

	// 3. 创建上传目录
	uploadDir := "./uploads/files"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建上传目录失败",
			"data":    nil,
		})
		return
	}

	// 4. 生成唯一文件名
	uuid := uuid.New().String()
	timestamp := time.Now().Format("20060102150405")
	newFilename := fmt.Sprintf("%s_%s%s", timestamp, uuid[:8], ext)
	savePath := filepath.Join(uploadDir, newFilename)

	// 5. 保存文件
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "文件保存失败",
			"data":    nil,
		})
		return
	}

	// 6. 返回访问URL (使用相对路径,前端会自动拼接当前域名)
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	host := c.Request.Host
	if host == "" {
		host = "localhost:8080" // 默认值
	}
	fileUrl := fmt.Sprintf("%s://%s/static/files/%s", scheme, host, newFilename)
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文件上传成功",
		"data": gin.H{
			"fileId":   uuid[:16],
			"filename": newFilename,
			"originalName": file.Filename,
			"fileSize": file.Size,
			"fileType": ext,
			"fileUrl":  fileUrl,
		},
	})
}

// UploadAvatar 上传头像
func UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请选择头像文件",
			"data":    nil,
		})
		return
	}

	// 1. 验证文件大小 (限制5MB)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "头像文件大小不能超过5MB",
			"data":    nil,
		})
		return
	}

	// 2. 验证文件类型 (只允许图片)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "只支持 JPG、PNG、GIF、WebP 格式的图片",
			"data":    nil,
		})
		return
	}

	// 3. 创建上传目录
	uploadDir := "./uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建上传目录失败",
			"data":    nil,
		})
		return
	}

	// 4. 生成唯一文件名
	uuid := uuid.New().String()
	timestamp := time.Now().Format("20060102150405")
	newFilename := fmt.Sprintf("%s_%s%s", timestamp, uuid[:8], ext)
	savePath := filepath.Join(uploadDir, newFilename)

	// 5. 保存文件
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "文件保存失败",
			"data":    nil,
		})
		return
	}

	// 6. 返回真实的访问URL (使用相对路径,前端会自动拼接当前域名)
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	host := c.Request.Host
	if host == "" {
		host = "localhost:8080" // 默认值
	}
	avatarUrl := fmt.Sprintf("%s://%s/static/avatars/%s", scheme, host, newFilename)
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "头像上传成功",
		"data": gin.H{
			"avatarUrl": avatarUrl,
			"filename":  newFilename,
			"size":      file.Size,
		},
	})
}
