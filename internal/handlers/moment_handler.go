package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Yw332/campus-moments-go/internal/service"
	"github.com/gin-gonic/gin"
)

// CreateMoment 发布动态
func CreateMoment(c *gin.Context) {
	type CreateMomentRequest struct {
		Content string   `json:"content" binding:"required"`
		Images  []string `json:"images"`
	}

	var req CreateMomentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 保存到数据库
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "发布成功",
		"data": gin.H{
			"id":      123,
			"content": req.Content,
			"images":  req.Images,
		},
	})
}

// GetMoments 获取动态列表（支持分页）
func GetMoments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	// 调用 service 层获取真实数据
	list, total, err := service.ListMoments(page, pageSize)
	if err != nil {
		// 如果数据库不可用或查询失败，回退到示例数据以保持接口可用
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"list": []gin.H{
					{
						"id":      1,
						"content": "今天天气真好！",
						"images":  []string{"https://example.com/1.jpg"},
						"author": gin.H{
							"id":       1,
							"username": "张三",
							"avatar":   "https://example.com/avatar.jpg",
						},
						"likeCount":    10,
						"commentCount": 5,
						"createdAt":    time.Now().Format("2006-01-02 15:04:05"),
					},
				},
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
		"code":    0,
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
	id := c.Param("id")
	// TODO: 从数据库查询详情
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id":      id,
			"content": "动态详情内容...",
		},
	})
}
