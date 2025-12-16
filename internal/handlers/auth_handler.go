package handlers

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/Yw332/campus-moments-go/internal/service"
	"github.com/Yw332/campus-moments-go/pkg/jwt"
	"github.com/Yw332/campus-moments-go/pkg/token_blacklist"
	"github.com/gin-gonic/gin"
)

var authService = service.NewAuthService()

// Register 用户注册
func Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	user, err := authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
		"data": gin.H{
			"userId":   user.ID,
			"username": user.Username,
			"phone":    user.Phone,
		},
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	response, err := authService.Login(&req)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "账户不存在" {
			statusCode = http.StatusNotFound
		}
		if err.Error() == "密码错误" {
			statusCode = http.StatusUnauthorized
		}
		if err.Error() == "账户已被锁定，请联系管理员" {
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, gin.H{
			"code":    statusCode,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data":    response,
	})
}

// Logout 退出登录
func Logout(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未认证",
			"data":    nil,
		})
		return
	}

	// 获取Token并添加到黑名单
	auth := c.GetHeader("Authorization")
	if auth != "" && strings.HasPrefix(auth, "Bearer ") {
		token := strings.TrimPrefix(auth, "Bearer ")
		
		// 解析Token获取过期时间
		if claims, err := jwt.ParseToken(token); err == nil {
			blacklist := token_blacklist.GetInstance()
			// 将Token添加到黑名单，直到其原定的过期时间
			blacklist.AddToken(token, claims.ExpiresAt.Time)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "退出成功",
		"data": gin.H{
			"userId":   userID,
			"logoutAt": time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}

// GetProfile 获取用户资料
func GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未认证",
			"data":    nil,
		})
		return
	}

	uid, _ := userID.(int64)
	user, err := authService.GetUserByID(uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "用户不存在",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"userId":   user.ID,
			"username": user.Username,
			"phone":    user.Phone,
		},
	})
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未认证",
			"data":    nil,
		})
		return
	}

	type ChangePasswordRequest struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	uid, _ := userID.(int64)
	if err := authService.UpdatePassword(uid, req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "密码修改成功",
		"data":    nil,
	})
}

// UpdateUserProfile 更新用户资料
func UpdateUserProfile(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未认证",
			"data":    nil,
		})
		return
	}

	uid, _ := userID.(int64)
	
	// 临时实现，实际应该从请求体获取数据
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data": gin.H{
			"userID": uid,
			"updated": true,
		},
	})
}

// SendVerificationCode 发送验证码
func SendVerificationCode(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "手机号不能为空",
			"data":    nil,
		})
		return
	}

	// 验证手机号格式
	if len(req.Phone) != 11 || req.Phone[0:1] != "1" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "手机号格式不正确",
			"data":    nil,
		})
		return
	}

	// 生成6位验证码
	verificationCode := fmt.Sprintf("%06d", rand.Intn(1000000))
	
	// 存储验证码（这里使用内存存储，生产环境建议使用Redis）
	// verificationCodes[req.Phone] = VerificationCode{
	// 	Code:      verificationCode,
	// 	ExpiresAt: time.Now().Add(5 * time.Minute),
	// }

	// 模拟发送短信
	log.Printf("验证码已发送至 %s: %s", req.Phone, verificationCode)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "验证码发送成功",
		"data": gin.H{
			"phone":       req.Phone,
			"expiresIn":   300, // 5分钟
			"resendAfter": 60,  // 1分钟后可重发
		},
	})
}

// VerifyAndResetPassword 验证验证码并重置密码
func VerifyAndResetPassword(c *gin.Context) {
	var req struct {
		Phone           string `json:"phone" binding:"required"`
		VerificationCode string `json:"verificationCode" binding:"required"`
		NewPassword     string `json:"newPassword" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 临时实现：演示用，实际应该验证存储的验证码
	// 生产环境需要验证：
	// 1. 验证码是否存在
	// 2. 验证码是否正确
	// 3. 验证码是否过期
	// 4. 手机号对应的用户是否存在

	// 模拟验证成功
	// 如果验证码正确，重置密码
	// if err := authService.ResetPasswordByPhone(req.Phone, req.NewPassword); err != nil {
	//     c.JSON(http.StatusBadRequest, gin.H{
	//         "code":    400,
	//         "message": "重置密码失败",
	//         "data":    nil,
	//     })
	//     return
	// }

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "密码重置成功",
		"data":    nil,
	})
}
