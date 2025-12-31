package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Yw332/campus-moments-go/internal/service"
	"github.com/gin-gonic/gin"
)

var momentService = service.NewMomentService()

// CreateMoment 发布动态
func CreateMoment(c *gin.Context) {
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

	var req service.CreateMomentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	uid := userID.(string)
	moment, err := momentService.CreateMoment(uid, &req)
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
		"message": "发布成功",
		"data":    moment,
	})
}

// GetMoments 获取动态列表（支持分页）
func GetMoments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	
	// 支持按用户ID筛选
	var userID *string
	if uidStr := c.Query("userId"); uidStr != "" {
		userID = &uidStr
	}

	list, total, err := momentService.ListMoments(page, pageSize, userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
			"data": gin.H{
				"list": []gin.H{},
				"pagination": gin.H{
					"page":     page,
					"pageSize": pageSize,
					"total":    0,
				},
			},
		})
		return
	}

	// 转换为前端需要的格式
	convertedList := make([]gin.H, 0, len(list))
	for _, moment := range list {
		// 提取第一张图片作为imageUrl
		var imageUrl string
		if moment.Images != nil && len(moment.Images) > 0 {
			var images []string
			if err := json.Unmarshal(moment.Images, &images); err == nil && len(images) > 0 {
				imageUrl = images[0]
			}
		}

		// 获取作者名称
		author := "未知作者"
		if moment.User != nil {
			author = moment.User.Username
		}

		// 格式化创建时间
		createTime := moment.CreatedAt.Format("2006-01-02 15:04")

		item := gin.H{
			"id":        moment.ID,
			"title":     moment.Title,
			"author":    author,
			"imageUrl":  imageUrl,
			"likeCount": moment.LikeCount,
			"createTime": createTime,
		}

		convertedList = append(convertedList, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"list": convertedList,
			"pagination": gin.H{
				"page":     page,
				"pageSize": pageSize,
				"total":    total,
			},
		},
	})
}

// GetMomentDetail 获取动态详情
func GetMomentDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的动态ID",
			"data":    nil,
		})
		return
	}

	moment, err := momentService.GetMomentByID(id)
	if err != nil {
		if err.Error() == "动态不存在" {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "动态不存在",
				"data":    nil,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": err.Error(),
				"data":    nil,
			})
		}
		return
	}

	// 添加 author 字段（兼容）
	data := make(map[string]interface{})
	if moment != nil {
		data["id"] = moment.ID
		data["userId"] = moment.UserID
		data["authorId"] = moment.UserID
		data["title"] = moment.Title
		data["content"] = moment.Content
		data["images"] = moment.Images
		data["video"] = moment.Video
		data["visibility"] = moment.Visibility
		data["status"] = moment.Status
		data["tags"] = moment.Tags
		data["likedUsers"] = moment.LikedUsers
		data["commentsSummary"] = moment.CommentsSummary
		data["likeCount"] = moment.LikeCount
		data["commentCount"] = moment.CommentCount
		data["viewCount"] = moment.ViewCount
		data["createdAt"] = moment.CreatedAt
		data["updatedAt"] = moment.UpdatedAt
		data["media"] = moment.Media
		data["author"] = moment.User
		data["user"] = moment.User
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    data,
	})
}

// UpdateMoment 更新动态（部分更新）
func UpdateMoment(c *gin.Context) {
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

	momentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的动态ID",
			"data":    nil,
		})
		return
	}

	var req service.UpdateMomentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	uid := userID.(string)
	moment, err := momentService.UpdateMoment(uid, momentID, &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := err.Error()
		
		if err.Error() == "动态不存在或无权限修改" {
			statusCode = http.StatusNotFound
		}
		
		c.JSON(statusCode, gin.H{
			"code":    statusCode,
			"message": message,
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    moment,
	})
}

// DeleteMoment 删除动态
func DeleteMoment(c *gin.Context) {
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

	momentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的动态ID",
			"data":    nil,
		})
		return
	}

	uid := userID.(string)
	if err := momentService.DeleteMoment(uid, momentID); err != nil {
		statusCode := http.StatusInternalServerError
		message := err.Error()

		if err.Error() == "动态不存在或无权限删除" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{
			"code":    statusCode,
			"message": message,
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
		"data": gin.H{
			"postId": momentID,
		},
	})
}

// AdminDeleteMoment 管理员删除动态（可删除任意用户的动态）
func AdminDeleteMoment(c *gin.Context) {
	momentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的动态ID",
			"data":    nil,
		})
		return
	}

	if err := momentService.AdminDeleteMoment(momentID); err != nil {
		statusCode := http.StatusInternalServerError
		message := err.Error()

		if err.Error() == "动态不存在" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{
			"code":    statusCode,
			"message": message,
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
		"data": gin.H{
			"postId": momentID,
		},
	})
}

// GetUserMoments 获取当前用户的所有动态
func GetUserMoments(c *gin.Context) {
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

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	uid := userID.(string)
	list, total, err := momentService.GetUserMoments(uid, page, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
			"data": gin.H{
				"list": []gin.H{},
				"pagination": gin.H{
					"page":     page,
					"pageSize": pageSize,
					"total":    0,
				},
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"list": list,
			"pagination": gin.H{
				"page":     page,
				"pageSize": pageSize,
				"total":    total,
			},
		},
	})
}
