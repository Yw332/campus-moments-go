package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Yw332/campus-moments-go/internal/models"
	"github.com/Yw332/campus-moments-go/internal/service"
	"github.com/Yw332/campus-moments-go/pkg/database"
	"github.com/Yw332/campus-moments-go/pkg/jwt"
	"github.com/gin-gonic/gin"
)

var adminService = service.NewAdminService()

// AdminLogin 管理员登录
func AdminLogin(c *gin.Context) {
	var req struct {
		Account  string `json:"account" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error(), "data": nil})
		return
	}

	admin, err := adminService.Authenticate(req.Account, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": err.Error(), "data": nil})
		return
	}

	// 生成 token（admin 不依赖 userId，使用 0）
	token, _ := jwt.GenerateToken(0, admin.Username)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "登录成功", "data": gin.H{"token": token}})
}

// AdminMenu 返回管理员导航菜单
func AdminMenu(c *gin.Context) {
	menu := []gin.H{
		{"name": "用户管理", "path": "/admin/users"},
		{"name": "内容管理", "path": "/admin/contents"},
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "ok", "data": menu})
}

// ListUsers 管理员列出用户
func ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	db := database.GetDB()
	var users []models.User
	if err := db.Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "ok", "data": users})
}

// GetUserDetail 获取指定用户
func GetUserDetail(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	var user models.User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "ok", "data": gin.H{"userId": user.ID, "username": user.Username, "phone": user.Phone}})
}

// ResetUserPassword 管理重置用户密码
func ResetUserPassword(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		NewPassword string `json:"newPassword" binding:"required"`
		Confirm     bool   `json:"confirm" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error(), "data": nil})
		return
	}
	if !req.Confirm {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请确认操作", "data": nil})
		return
	}

	a := service.NewAuthService()
	// 使用字符串ID版本的更新
	if err := a.UpdatePasswordStr(id, "", req.NewPassword); err != nil {
		// UpdatePasswordStr 現在需要旧密码，为管理重置应直接更新
		db := database.GetDB()
		var user models.User
		if err := db.Where("id = ?", id).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在", "data": nil})
			return
		}
		hashed := []byte(req.NewPassword)
		user.Password = string(hashed)
		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "重置密码失败", "data": nil})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "重置成功", "data": nil})
}

// DeleteUser 删除指定用户
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	confirm := c.Query("confirm")
	if confirm != "true" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请确认删除（?confirm=true）", "data": nil})
		return
	}
	db := database.GetDB()
	if err := db.Where("id = ?", id).Delete(&models.User{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功", "data": nil})
}

// ListContents 管理列出内容（使用 posts 表）
func ListContents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	db := database.GetDB()
	// posts 表映射到 models.Moment 可能不存在；我们使用 posts 表原始查询
	type Post struct {
		ID        int64     `json:"id"`
		UserID    string    `json:"userId"`
		Title     string    `json:"title"`
		Content   string    `json:"content"`
		Status    int       `json:"status"`
		CreatedAt time.Time `json:"createdAt"`
	}
	var posts []Post
	if err := db.Raw("SELECT id, user_id, title, content, status, created_at FROM posts ORDER BY created_at DESC LIMIT ? OFFSET ?", pageSize, offset).Scan(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "ok", "data": posts})
}

// GetContentDetail 内容详情
func GetContentDetail(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	var post struct {
		ID        int64     `json:"id"`
		UserID    string    `json:"userId"`
		Title     string    `json:"title"`
		Content   string    `json:"content"`
		Status    int       `json:"status"`
		CreatedAt time.Time `json:"createdAt"`
	}
	if err := db.Raw("SELECT id, user_id, title, content, status, created_at FROM posts WHERE id = ?", id).Scan(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "内容不存在", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "ok", "data": post})
}

// DeleteContent 管理删除内容
func DeleteContent(c *gin.Context) {
	id := c.Param("id")
	confirm := c.Query("confirm")
	if confirm != "true" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请确认删除（?confirm=true）", "data": nil})
		return
	}
	db := database.GetDB()
	// 软删除或直接删除，这里直接删除
	if err := db.Exec("DELETE FROM posts WHERE id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功", "data": nil})
}
