package handlers

import (
	"github.com/Yw332/campus-moments-go/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// LikePost 点赞/取消点赞帖子
func LikePost(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("postId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无效的帖子ID",
			"data":    nil,
		})
		return
	}

	userID := c.GetString("userID")
	liked, err := service.ToggleLikePost(postID, userID)
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

// GetPostLikes 获取帖子点赞列表
func GetPostLikes(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("postId"), 10, 64)
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

	likes, total, err := service.GetPostLikes(postID, page, pageSize)
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
			"likes": likes,
			"total": total,
			"page":  page,
			"pageSize": pageSize,
		},
	})
}

// GetCommentLikes 获取评论点赞列表
func GetCommentLikes(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("commentId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无效的评论ID",
			"data":    nil,
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	likes, total, err := service.GetCommentLikes(commentID, page, pageSize)
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
			"likes": likes,
			"total": total,
			"page":  page,
			"pageSize": pageSize,
		},
	})
}

// GetUserLikes 获取用户点赞列表
func GetUserLikes(c *gin.Context) {
	targetUserID := c.Param("userId")
	if targetUserID == "" {
		targetUserID = c.GetString("userID")
	}

	targetType := c.DefaultQuery("type", "1") // 默认获取帖子点赞
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	likes, total, err := service.GetUserLikes(targetUserID, targetType, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "获取失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 转换为响应格式,修正avatar字段名
	convertedLikes := make([]map[string]interface{}, 0, len(likes))
	for _, like := range likes {
		likeData := map[string]interface{}{
			"id":         like.ID,
			"targetId":   like.TargetID,
			"targetType": like.TargetType,
			"createdAt":  like.CreatedAt,
		}

		// 添加用户信息,统一使用avatarUrl
		if like.User.ID != "" {
			likeData["user"] = map[string]interface{}{
				"id":              like.User.ID,
				"username":        like.User.Username,
				"avatarUrl":       like.User.AvatarURL,
				"avatarType":      like.User.AvatarType,
				"avatarUpdatedAt": like.User.AvatarUpdatedAt,
				"signature":       like.User.Signature,
			}
		}

		// 如果是帖子点赞,添加帖子信息
		if like.TargetType == 1 && like.Post.ID != 0 {
			likeData["post"] = map[string]interface{}{
				"id":          like.Post.ID,
				"title":       like.Post.Title,
				"content":     like.Post.Content,
				"images":      like.Post.Images,
				"video":       like.Post.Video,
				"createdAt":   like.Post.CreatedAt,
				"likeCount":   like.Post.LikeCount,
				"commentCount": like.Post.CommentCount,
				"viewCount":   like.Post.ViewCount,
				"visibility":  like.Post.Visibility,
			}
		}

		convertedLikes = append(convertedLikes, likeData)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "获取成功",
		"data": gin.H{
			"likes":    convertedLikes,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}