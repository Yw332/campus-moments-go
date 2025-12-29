package handlers

import (
	"github.com/Yw332/campus-moments-go/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CreateTag 创建标签
func CreateTag(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required,min=1,max=20"`
		Color       string `json:"color" binding:"len=7"`
		Icon        string `json:"icon" binding:"max=50"`
		Description string `json:"description" binding:"max=200"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}
	
	tag, err := service.CreateTag(req.Name, req.Color, req.Icon, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"code": http.StatusCreated,
		"message":  "创建成功",
		"data": tag,
	})
}

// GetTagList 获取标签列表
func GetTagList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))
	status := c.DefaultQuery("status", "0") // 默认只获取正常标签
	
	tags, total, err := service.GetTagList(status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "获取成功",
		"data": gin.H{
			"tags": tags,
			"total": total,
			"page":  page,
			"pageSize": pageSize,
		},
	})
}

// GetTagDetail 获取标签详情
func GetTagDetail(c *gin.Context) {
	var tagIDStr string
	if tagIDStr = c.Param("id"); tagIDStr == "" {
		tagIDStr = c.GetString("tagId")
	}
	
	tagID, err := strconv.ParseInt(tagIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的标签ID"})
		return
	}
	
	tag, err := service.GetTagDetail(tagID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "标签不存在"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "获取成功",
		"data": tag,
	})
}

// UpdateTag 更新标签
func UpdateTag(c *gin.Context) {
	tagID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的标签ID"})
		return
	}
	
	var req struct {
		Name        string `json:"name" binding:"min=1,max=20"`
		Color       string `json:"color" binding:"len=7"`
		Icon        string `json:"icon" binding:"max=50"`
		Description string `json:"description" binding:"max=200"`
		Status      int    `json:"status" binding:"oneof=0 1"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}
	
	tag, err := service.UpdateTag(tagID, req.Name, req.Color, req.Icon, req.Description, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "更新成功",
		"data": tag,
	})
}

// DeleteTag 删除标签
func DeleteTag(c *gin.Context) {
	tagID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的标签ID"})
		return
	}
	
	err = service.DeleteTag(tagID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "删除成功",
	})
}

// GetHotTags 获取热门标签
func GetHotTags(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	
	tags, err := service.GetHotTags(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "获取成功",
		"data": tags,
	})
}

// SearchTags 搜索标签
func SearchTags(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}
	
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	
	tags, err := service.SearchTags(keyword, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "搜索成功",
		"data": tags,
	})
}

// GetTagPosts 获取标签相关帖子
func GetTagPosts(c *gin.Context) {
	tagName := c.Param("name")
	if tagName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标签名不能为空"})
		return
	}
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	userID := c.GetString("userID")
	
	posts, total, err := service.GetTagPosts(tagName, userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "获取成功",
		"data": gin.H{
			"posts": posts,
			"total": total,
			"page":  page,
			"pageSize": pageSize,
		},
	})
}