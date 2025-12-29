package handlers

import (
	"github.com/Yw332/campus-moments-go/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// SendMessage 发送消息
func SendMessage(c *gin.Context) {
	var req struct {
		ReceiverID     string `json:"receiverId" binding:"required"`
		MsgType        int    `json:"msgType" binding:"required,oneof=1 2 3 4"` // 1-文本 2-图片 3-视频 4-文件
		ContentPreview string `json:"contentPreview"`
		FileURL        string `json:"fileUrl"`
		FileSize       int    `json:"fileSize"`
		IsEncrypted    bool   `json:"isEncrypted"`
		DeviceID       string `json:"deviceId"`
		ServerMsgID    string `json:"serverMsgId"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}
	
	senderID := c.GetString("userID")
	message, err := service.SendMessage(senderID, req.ReceiverID, req.MsgType, 
		req.ContentPreview, req.FileURL, req.FileSize, req.IsEncrypted, 
		req.DeviceID, req.ServerMsgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"code": http.StatusOK,
		"message":  "发送成功",
		"data": message,
	})
}

// GetMessageList 获取消息列表
func GetMessageList(c *gin.Context) {
	userID := c.GetString("userID")
	peerID := c.Param("peerId")
	if peerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "会话对象ID不能为空"})
		return
	}
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))
	beforeMsgID := c.Query("beforeMsgId") // 分页使用
	
	messages, total, err := service.GetMessageList(userID, peerID, beforeMsgID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "获取成功",
		"data": gin.H{
			"messages": messages,
			"total": total,
			"page":  page,
			"pageSize": pageSize,
		},
	})
}

// GetConversationList 获取会话列表
func GetConversationList(c *gin.Context) {
	userID := c.GetString("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	
	conversations, total, err := service.GetConversationList(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "获取成功",
		"data": gin.H{
			"conversations": conversations,
			"total": total,
			"page":  page,
			"pageSize": pageSize,
		},
	})
}

// MarkMessagesAsRead 标记消息为已读
func MarkMessagesAsRead(c *gin.Context) {
	peerID := c.Param("peerId")
	if peerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "会话对象ID不能为空"})
		return
	}
	
	userID := c.GetString("userID")
	err := service.MarkMessagesAsRead(userID, peerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "标记失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "标记成功",
	})
}

// PinConversation 置顶会话
func PinConversation(c *gin.Context) {
	peerID := c.Param("peerId")
	if peerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "会话对象ID不能为空"})
		return
	}
	
	userID := c.GetString("userID")
	err := service.PinConversation(userID, peerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "置顶成功",
	})
}

// UnpinConversation 取消置顶会话
func UnpinConversation(c *gin.Context) {
	peerID := c.Param("peerId")
	if peerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "会话对象ID不能为空"})
		return
	}
	
	userID := c.GetString("userID")
	err := service.UnpinConversation(userID, peerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "取消置顶成功",
	})
}

// MuteConversation 静音会话
func MuteConversation(c *gin.Context) {
	peerID := c.Param("peerId")
	if peerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "会话对象ID不能为空"})
		return
	}
	
	userID := c.GetString("userID")
	err := service.MuteConversation(userID, peerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "静音成功",
	})
}

// UnmuteConversation 取消静音会话
func UnmuteConversation(c *gin.Context) {
	peerID := c.Param("peerId")
	if peerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "会话对象ID不能为空"})
		return
	}
	
	userID := c.GetString("userID")
	err := service.UnmuteConversation(userID, peerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "取消静音成功",
	})
}

// DeleteConversation 删除会话
func DeleteConversation(c *gin.Context) {
	peerID := c.Param("peerId")
	if peerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "会话对象ID不能为空"})
		return
	}
	
	userID := c.GetString("userID")
	err := service.DeleteConversation(userID, peerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "删除成功",
	})
}

// GetUnreadCount 获取未读消息数
func GetUnreadCount(c *gin.Context) {
	userID := c.GetString("userID")
	count, err := service.GetUnreadCount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message":  "获取成功",
		"data": gin.H{
			"unreadCount": count,
		},
	})
}