package service

import (
	"github.com/Yw332/campus-moments-go/internal/models"
	"gorm.io/gorm"
	"time"
)

// SendFriendRequest 发送好友请求
func SendFriendRequest(fromUserID, toUserID, message string) (*models.FriendRequest, error) {
	// 检查是否已经是好友
	if IsFriend(fromUserID, toUserID) {
		return nil, gorm.ErrInvalidTransaction // 已是好友
	}
	
	// 检查是否已有待处理的请求
	var existingRequest models.FriendRequest
	err := getDB().Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?) AND status = 0",
		fromUserID, toUserID, toUserID, fromUserID).First(&existingRequest).Error
	if err == nil {
		return nil, gorm.ErrInvalidTransaction // 已有待处理请求
	}
	
	// 创建好友请求
	request := &models.FriendRequest{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Message:    message,
		Status:     0, // 待处理
		ExpiresAt:  time.Now().Add(7 * 24 * time.Hour), // 7天后过期
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	
	if err := getDB().Create(request).Error; err != nil {
		return nil, err
	}
	
	// 关联用户信息
	getDB().Preload("FromUser").Preload("ToUser").First(request, request.ID)
	
	return request, nil
}

// GetFriendRequests 获取好友请求列表
func GetFriendRequests(userID, requestType string, page, pageSize int) ([]models.FriendRequest, int64, error) {
	var requests []models.FriendRequest
	var total int64
	
	offset := (page - 1) * pageSize
	
	var query *gorm.DB
	
	if requestType == "sent" {
		// 我发送的请求
		query = getDB().Model(&models.FriendRequest{}).Where("from_user_id = ?", userID)
	} else {
		// 我收到的请求
		query = getDB().Model(&models.FriendRequest{}).Where("to_user_id = ?", userID)
	}
	
	// 只获取未过期的请求
	query = query.Where("expires_at > ?", time.Now())
	
	// 获取总数
	query.Count(&total)
	
	// 获取请求列表
	err := query.Preload("FromUser").
		Preload("ToUser").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&requests).Error
	
	return requests, total, err
}

// HandleFriendRequest 处理好友请求
func HandleFriendRequest(requestID int64, userID, action string) error {
	var request models.FriendRequest
	
	// 查找请求（只能处理发给自己的请求）
	if err := getDB().First(&request, "id = ? AND to_user_id = ? AND status = 0 AND expires_at > ?", 
		requestID, userID, time.Now()).Error; err != nil {
		return err
	}
	
	if action == "accept" {
		// 同意好友请求，创建双向好友关系
		now := time.Now()
		
		// 创建好友关系（发送者 -> 接收者）
		relation1 := models.FriendRelation{
			UserID:       request.FromUserID,
			FriendID:     request.ToUserID,
			RelationType: 1, // 好友
			Status:       0,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		
		// 创建好友关系（接收者 -> 发送者）
		relation2 := models.FriendRelation{
			UserID:       request.ToUserID,
			FriendID:     request.FromUserID,
			RelationType: 1, // 好友
			Status:       0,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		
		// 开启事务
		tx := getDB().Begin()
		
		if err := tx.Create(&relation1).Error; err != nil {
			tx.Rollback()
			return err
		}
		
		if err := tx.Create(&relation2).Error; err != nil {
			tx.Rollback()
			return err
		}
		
		// 更新请求状态
		if err := tx.Model(&request).Updates(map[string]interface{}{
			"status":     1, // 已同意
			"updated_at": now,
		}).Error; err != nil {
			tx.Rollback()
			return err
		}
		
		tx.Commit()
		
	} else if action == "reject" {
		// 拒绝好友请求
		if err := getDB().Model(&request).Updates(map[string]interface{}{
			"status":     2, // 已拒绝
			"updated_at": time.Now(),
		}).Error; err != nil {
			return err
		}
	}
	
	return nil
}

// GetFriendList 获取好友列表
func GetFriendList(userID, keyword string, page, pageSize int) ([]models.FriendRelation, int64, error) {
	var friends []models.FriendRelation
	var total int64
	
	offset := (page - 1) * pageSize
	
	query := getDB().Model(&models.FriendRelation{}).
		Where("user_id = ? AND relation_type = 1 AND status = 0", userID)
	
	// 关键词搜索
	if keyword != "" {
		query = query.Joins("LEFT JOIN users ON friend_id = users.id").
			Where("users.username LIKE ? OR users.wechat_nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	
	// 获取总数
	query.Count(&total)
	
	// 获取好友列表
	err := query.Preload("Friend").
		Order("updated_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&friends).Error
	
	return friends, total, err
}

// DeleteFriend 删除好友
func DeleteFriend(userID, friendID string) error {
	// 检查是否为好友关系
	if !IsFriend(userID, friendID) {
		return gorm.ErrRecordNotFound
	}
	
	// 开启事务删除双向关系
	tx := getDB().Begin()
	
	// 删除用户对好友的关系
	if err := tx.Model(&models.FriendRelation{}).
		Where("user_id = ? AND friend_id = ? AND relation_type = 1", userID, friendID).
		Update("status", 1).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// 删除好友对用户的关系
	if err := tx.Model(&models.FriendRelation{}).
		Where("user_id = ? AND friend_id = ? AND relation_type = 1", friendID, userID).
		Update("status", 1).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	tx.Commit()
	
	return nil
}

// UpdateFriendRemark 更新好友备注
func UpdateFriendRemark(userID, friendID, remarkName string) error {
	// 检查是否为好友关系
	if !IsFriend(userID, friendID) {
		return gorm.ErrRecordNotFound
	}
	
	// 更新备注
	result := getDB().Model(&models.FriendRelation{}).
		Where("user_id = ? AND friend_id = ? AND relation_type = 1 AND status = 0", userID, friendID).
		Updates(map[string]interface{}{
			"remark_name": remarkName,
			"updated_at":  time.Now(),
		})
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	
	return nil
}

// SearchFriends 搜索好友
func SearchFriends(userID, keyword string, page, pageSize int) ([]models.FriendRelation, int64, error) {
	return GetFriendList(userID, keyword, page, pageSize)
}