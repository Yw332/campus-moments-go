package handlers

<<<<<<< HEAD
import (
	"net/http"

	"github.com/Yw332/campus-moments-go/internal/service"
	"github.com/gin-gonic/gin"
)

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

	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	uid := userID.(string)
	// Log the update request
	println("Updating profile for user:", uid)

	user, err := authService.UpdateProfile(uid, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data": gin.H{
			"userId":   user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"gender":   user.Gender,
			"bio":      user.Bio,
		},
	})
}
=======
// 用户相关函数在 auth_handler.go 中已实现，避免重复定义
>>>>>>> be7109d45b16980427c35fc3f6c3874bbda68e13
