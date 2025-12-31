package handlers

import (
	"github.com/Yw332/campus-moments-go/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// getPostID 从请求中获取帖子ID的辅助函数
func getPostID(c *gin.Context) (int64, error) {
	var postIDStr string
	
	// 尝试从不同的参数名获取帖子ID
	if postIDStr = c.Param("postId"); postIDStr == "" {
		if postIDStr = c.GetString("postId"); postIDStr == "" {
			postIDStr = c.Param("id")
		}
	}
	
	return strconv.ParseInt(postIDStr, 10, 64)
}

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	postID, err := getPostID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无效的帖子ID",
			"data":    nil,
		})
		return
	}
	
	var req struct {
		Content string                 `json:"content" binding:"required,min=1,max=1000"`
		Replies []map[string]interface{} `json:"replies"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	userID := c.GetString("userID")
	comment, err := service.CreateComment(postID, userID, req.Content, req.Replies)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "创建失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"code": http.StatusOK,
		"message":  "评论成功",
		"data": comment,
	})
}

// GetCommentList 获取评论列表
func GetCommentList(c *gin.Context) {
	postID, err := getPostID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无效的帖子ID",
			"data":    nil,
		})
		return
	}
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	
	comments, total, err := service.GetCommentList(postID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "获取失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "获取成功",
		"data": gin.H{
			"comments": comments,
			"total": total,
			"page":  page,
			"pageSize": pageSize,
		},
	})
}

// UpdateComment 更新评论
func UpdateComment(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无效的评论ID",
			"data":    nil,
		})
		return
	}
	
	var req struct {
		Content string `json:"content" binding:"required,min=1,max=1000"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	userID := c.GetString("userID")
	comment, err := service.UpdateComment(commentID, userID, req.Content)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    http.StatusForbidden,
			"message": "更新失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "更新成功",
		"data": comment,
	})
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无效的评论ID",
			"data":    nil,
		})
		return
	}
	
	userID := c.GetString("userID")
	err = service.DeleteComment(commentID, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    http.StatusForbidden,
			"message": "删除失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "删除成功",
	})
}

// LikeComment 点赞评论
func LikeComment(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无效的评论ID",
			"data":    nil,
		})
		return
	}
	
	userID := c.GetString("userID")
	liked, err := service.ToggleLikeComment(commentID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "操作失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	var msg string
	if liked {
		msg = "点赞成功"
	} else {
		msg = "取消点赞成功"
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  msg,
		"data": gin.H{
			"liked": liked,
		},
	})
}

// ReplyComment 回复评论
func ReplyComment(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无效的评论ID",
			"data":    nil,
		})
		return
	}
	
	var req struct {
		Content string `json:"content" binding:"required,min=1,max=1000"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	userID := c.GetString("userID")
	comment, err := service.ReplyComment(commentID, userID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "回复失败: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "回复成功",
		"data": comment,
	})
}

// AdminDeleteComment 管理员删除评论
func AdminDeleteComment(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无效的评论ID",
			"data":    nil,
		})
		return
	}

	comment, err := service.AdminDeleteComment(commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "删除失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "删除成功",
		"data": gin.H{
			"commentId": comment.ID,
		},
	})
}