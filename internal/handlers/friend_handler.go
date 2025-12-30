package handlers

import (
	"github.com/Yw332/campus-moments-go/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// SendFriendRequest 发送好友请求
func SendFriendRequest(c *gin.Context) {
	var req struct {
		ToUserID string `json:"toUserId" binding:"required"`
		Message  string `json:"message" binding:"max=200"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}
	
	fromUserID := c.GetString("userID")
	request, err := service.SendFriendRequest(fromUserID, req.ToUserID, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"code": http.StatusOK,
		"message":  "发送成功",
		"data": request,
	})
}

// GetFriendRequests 获取好友请求列表
func GetFriendRequests(c *gin.Context) {
	userID := c.GetString("userID")
	requestType := c.DefaultQuery("type", "received") // sent/received
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	
	requests, total, err := service.GetFriendRequests(userID, requestType, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "获取成功",
		"data": gin.H{
			"requests": requests,
			"total": total,
			"page":  page,
			"pageSize": pageSize,
		},
	})
}

// HandleFriendRequest 处理好友请求
func HandleFriendRequest(c *gin.Context) {
	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求ID"})
		return
	}
	
	var req struct {
		Action string `json:"action" binding:"required,oneof=accept reject"` // accept/reject
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}
	
	userID := c.GetString("userID")
	err = service.HandleFriendRequest(requestID, userID, req.Action)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "处理失败: " + err.Error()})
		return
	}

	var msg string
	if req.Action == "accept" {
		msg = "已同意好友请求"
	} else {
		msg = "已拒绝好友请求"
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": msg,
		"data":    nil,
	})
}

// GetFriendList 获取好友列表
func GetFriendList(c *gin.Context) {
	userID := c.GetString("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))
	keyword := c.Query("keyword")
	
	friends, total, err := service.GetFriendList(userID, keyword, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "获取成功",
		"data": gin.H{
			"friends": friends,
			"total": total,
			"page":  page,
			"pageSize": pageSize,
		},
	})
}

// DeleteFriend 删除好友
func DeleteFriend(c *gin.Context) {
	friendID := c.Param("friendId")
	if friendID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "好友ID不能为空"})
		return
	}

	userID := c.GetString("userID")
	err := service.DeleteFriend(userID, friendID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "删除成功",
		"data":    nil,
	})
}

// UpdateFriendRemark 更新好友备注
func UpdateFriendRemark(c *gin.Context) {
	friendID := c.Param("friendId")
	if friendID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "好友ID不能为空"})
		return
	}

	var req struct {
		RemarkName string `json:"remarkName" binding:"max=50"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	userID := c.GetString("userID")
	err := service.UpdateFriendRemark(userID, friendID, req.RemarkName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "更新成功",
		"data":    nil,
	})
}

// SearchFriends 搜索好友
func SearchFriends(c *gin.Context) {
	userID := c.GetString("userID")
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	
	friends, total, err := service.SearchFriends(userID, keyword, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "搜索成功",
		"data": gin.H{
			"friends": friends,
			"total": total,
			"page":  page,
			"pageSize": pageSize,
		},
	})
}