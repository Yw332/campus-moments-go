package handlers

import (
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

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    moment,
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
