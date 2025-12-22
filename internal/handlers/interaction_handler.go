package handlers

import (
	"net/http"
	"strconv"

	"github.com/Yw332/campus-moments-go/internal/models"
	"github.com/Yw332/campus-moments-go/internal/service"
	"github.com/gin-gonic/gin"
)

var interactionService = service.NewInteractionService()

// CreateComment 发表评论
func CreateComment(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	var req service.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	comment, err := interactionService.CreateComment(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "评论成功",
		"data": comment,
	})
}

// GetMomentComments 获取动态评论
func GetMomentComments(c *gin.Context) {
	momentIDStr := c.Param("id")
	momentID, err := strconv.ParseInt(momentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的动态ID", "data": nil})
		return
	}

	comments, err := interactionService.GetMomentComments(momentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "获取成功",
		"data": comments,
	})
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context) {
	userID, _ := c.Get("userID")
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的评论ID", "data": nil})
		return
	}

	if err := interactionService.DeleteComment(userID.(string), commentID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "删除成功",
		"data": nil,
	})
}

// LikeRequest 点赞请求
type LikeRequest struct {
	TargetID   int64 `json:"targetId" binding:"required"`
	TargetType int   `json:"targetType" binding:"required"` // 1:动态 2:评论
}

// ToggleLike 点赞/取消点赞
func ToggleLike(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req LikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	if req.TargetType != models.LikeTargetTypeMoment && req.TargetType != models.LikeTargetTypeComment {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的目标类型", "data": nil})
		return
	}

	isLiked, err := interactionService.ToggleLike(userID.(string), req.TargetID, req.TargetType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error(), "data": nil})
		return
	}

	msg := "取消点赞成功"
	if isLiked {
		msg = "点赞成功"
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": msg,
		"data": gin.H{
			"isLiked": isLiked,
		},
	})
}
