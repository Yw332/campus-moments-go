package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
func Register(c *gin.Context) {
	type RegisterRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 保存用户到数据库
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "注册成功",
		"data": gin.H{
			"userId":   1,
			"username": req.Username,
		},
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	type LoginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 验证用户名密码，生成JWT
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"data": gin.H{
			"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
			"userInfo": gin.H{
				"id":       1,
				"username": req.Username,
				"avatar":   "https://example.com/avatar.jpg",
			},
		},
	})
}

// Logout 退出登录
func Logout(c *gin.Context) {
	// TODO: 使token失效
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "退出成功",
	})
}
