package service

import (
	"github.com/Yw332/campus-moments-go/internal/models"
	"gorm.io/gorm"
	"time"
)

// SendMessage 发送消息
func SendMessage(senderID, receiverID string, msgType int, contentPreview, fileURL string, 
	fileSize int, isEncrypted bool, deviceID, serverMsgID string) (*models.Message, error) {
	
	// 创建消息
	message := &models.Message{
		SenderID:       senderID,
		ReceiverID:     receiverID,
		MsgType:        msgType,
		ContentPreview: contentPreview,
		FileURL:        fileURL,
		FileSize:       fileSize,
		IsEncrypted:    isEncrypted,
		IsRead:         false,
		DeviceID:       deviceID,
		ServerMsgID:    serverMsgID,
		CreatedAt:      time.Now(),
	}
	
	// 开启事务
	tx := getDB().Begin()
	
	// 保存消息
	if err := tx.Create(message).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	
	// 关联用户信息
	tx.Preload("Sender").Preload("Receiver").First(message, message.ID)
	
	// 更新或创建会话记录
	var conversation models.Conversation
	err := tx.Where("user_id = ? AND peer_id = ?", senderID, receiverID).
		First(&conversation).Error
	
	if err == gorm.ErrRecordNotFound {
		// 创建新会话
		conversation = models.Conversation{
			UserID:          senderID,
			PeerID:          receiverID,
			LastMsgID:       message.ID,
			LastMsgPreview:  contentPreview,
			UnreadCount:     0, // 发送者未读数为0
			IsPinned:        false,
			IsMuted:         false,
			UpdatedAt:       time.Now(),
		}
		if err := tx.Create(&conversation).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	} else if err == nil {
		// 更新现有会话
		if err := tx.Model(&conversation).Updates(map[string]interface{}{
			"last_msg_id":      message.ID,
			"last_msg_preview": contentPreview,
			"updated_at":       time.Now(),
		}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	} else {
		tx.Rollback()
		return nil, err
	}
	
	// 更新接收者的会话记录
	var receiverConversation models.Conversation
	err = tx.Where("user_id = ? AND peer_id = ?", receiverID, senderID).
		First(&receiverConversation).Error
	
	if err == gorm.ErrRecordNotFound {
		// 创建接收者的会话
		receiverConversation = models.Conversation{
			UserID:          receiverID,
			PeerID:          senderID,
			LastMsgID:       message.ID,
			LastMsgPreview:  contentPreview,
			UnreadCount:     1, // 接收者未读数+1
			IsPinned:        false,
			IsMuted:         false,
			UpdatedAt:       time.Now(),
		}
		if err := tx.Create(&receiverConversation).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	} else if err == nil {
		// 更新接收者的会话
		if err := tx.Model(&receiverConversation).Updates(map[string]interface{}{
			"last_msg_id":      message.ID,
			"last_msg_preview": contentPreview,
			"unread_count":     gorm.Expr("unread_count + ?", 1),
			"updated_at":       time.Now(),
		}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	} else {
		tx.Rollback()
		return nil, err
	}
	
	tx.Commit()
	
	return message, nil
}

// GetMessageList 获取消息列表
func GetMessageList(userID, peerID, beforeMsgID string, page, pageSize int) ([]models.Message, int64, error) {
	var messages []models.Message
	var total int64
	
	// 构建查询条件
	query := getDB().Model(&models.Message{}).
		Where("((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?))",
			userID, peerID, peerID, userID)
	
	// 如果指定了beforeMsgID，获取更早的消息
	if beforeMsgID != "" {
		query = query.Where("id < ?", beforeMsgID)
	}
	
	// 获取总数（简化处理，实际可以优化）
	query.Count(&total)
	
	// 获取消息列表
	err := query.Preload("Sender").
		Preload("Receiver").
		Order("created_at DESC").
		Limit(pageSize).
		Find(&messages).Error
	
	// 反转消息顺序（最新的在后面）
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	
	return messages, total, err
}

// GetConversationList 获取会话列表
func GetConversationList(userID string, page, pageSize int) ([]models.Conversation, int64, error) {
	var conversations []models.Conversation
	var total int64
	
	offset := (page - 1) * pageSize
	
	query := getDB().Model(&models.Conversation{}).Where("user_id = ?", userID)
	
	// 获取总数
	query.Count(&total)
	
	// 获取会话列表
	err := query.Preload("User").
		Preload("Peer").
		Preload("LastMessage").
		Order("is_pinned DESC, updated_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&conversations).Error
	
	return conversations, total, err
}

// MarkMessagesAsRead 标记消息为已读
func MarkMessagesAsRead(userID, peerID string) error {
	// 标记对方发给自己的消息为已读
	result := getDB().Model(&models.Message{}).
		Where("sender_id = ? AND receiver_id = ? AND is_read = ?", peerID, userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
		})
	
	if result.Error != nil {
		return result.Error
	}
	
	// 重置会话的未读数
	if result.RowsAffected > 0 {
		getDB().Model(&models.Conversation{}).
			Where("user_id = ? AND peer_id = ?", userID, peerID).
			Update("unread_count", 0)
	}
	
	return nil
}

// PinConversation 置顶会话
func PinConversation(userID, peerID string) error {
	result := getDB().Model(&models.Conversation{}).
		Where("user_id = ? AND peer_id = ?", userID, peerID).
		Update("is_pinned", true)
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	
	return nil
}

// UnpinConversation 取消置顶会话
func UnpinConversation(userID, peerID string) error {
	result := getDB().Model(&models.Conversation{}).
		Where("user_id = ? AND peer_id = ?", userID, peerID).
		Update("is_pinned", false)
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	
	return nil
}

// MuteConversation 静音会话
func MuteConversation(userID, peerID string) error {
	result := getDB().Model(&models.Conversation{}).
		Where("user_id = ? AND peer_id = ?", userID, peerID).
		Update("is_muted", true)
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	
	return nil
}

// UnmuteConversation 取消静音会话
func UnmuteConversation(userID, peerID string) error {
	result := getDB().Model(&models.Conversation{}).
		Where("user_id = ? AND peer_id = ?", userID, peerID).
		Update("is_muted", false)
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	
	return nil
}

// DeleteConversation 删除会话
func DeleteConversation(userID, peerID string) error {
	// 删除会话记录
	result := getDB().Model(&models.Conversation{}).
		Where("user_id = ? AND peer_id = ?", userID, peerID).
		Delete(&models.Conversation{})
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	
	// 注意：这里不删除消息记录，只是删除会话列表显示
	// 如果需要彻底删除消息，可以另外提供接口
	
	return nil
}

// GetUnreadCount 获取未读消息数
func GetUnreadCount(userID string) (int, error) {
	var total int64
	
	// 统计所有会话的未读消息总数
	getDB().Model(&models.Conversation{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(unread_count), 0)").
		Scan(&total)
	
	return int(total), nil
}