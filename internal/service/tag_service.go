package service

import (
	"github.com/Yw332/campus-moments-go/internal/models"
	"gorm.io/gorm"
	"time"
)

// CreateTag 创建标签
func CreateTag(name, color, icon, description string) (*models.Tag, error) {
	// 检查标签名是否已存在
	var existingTag models.Tag
	err := getDB().Where("name = ?", name).First(&existingTag).Error
	if err == nil {
		return nil, gorm.ErrInvalidTransaction // 标签已存在
	}
	
	tag := &models.Tag{
		Name:        name,
		Color:       color,
		Icon:        icon,
		Description: description,
		UsageCount:  0,
		Status:      0, // 正常状态
		CreatedAt:   time.Now(),
	}
	
	if err := getDB().Create(tag).Error; err != nil {
		return nil, err
	}
	
	return tag, nil
}

// GetTagList 获取标签列表
func GetTagList(status string, page, pageSize int) ([]models.Tag, int64, error) {
	var tags []models.Tag
	var total int64
	
	offset := (page - 1) * pageSize
	
	query := getDB().Model(&models.Tag{})
	
	// 根据状态过滤
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	// 获取总数
	query.Count(&total)
	
	// 获取标签列表
	err := query.Order("usage_count DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&tags).Error
	
	return tags, total, err
}

// GetTagDetail 获取标签详情
func GetTagDetail(tagID int64) (*models.Tag, error) {
	var tag models.Tag
	
	err := getDB().First(&tag, "id = ?", tagID).Error
	if err != nil {
		return nil, err
	}
	
	return &tag, nil
}

// UpdateTag 更新标签
func UpdateTag(tagID int64, name, color, icon, description string, status int) (*models.Tag, error) {
	var tag models.Tag
	
	// 查找标签
	if err := getDB().First(&tag, "id = ?", tagID).Error; err != nil {
		return nil, err
	}
	
	// 检查标签名是否被其他标签使用
	if name != tag.Name {
		var existingTag models.Tag
		err := getDB().Where("name = ? AND id != ?", name, tagID).First(&existingTag).Error
		if err == nil {
			return nil, gorm.ErrInvalidTransaction // 标签名已存在
		}
	}
	
	// 更新字段
	updates := map[string]interface{}{
		"name":        name,
		"color":       color,
		"icon":        icon,
		"description": description,
		"status":      status,
		"updated_at":  time.Now(),
	}
	
	if err := getDB().Model(&tag).Updates(updates).Error; err != nil {
		return nil, err
	}
	
	// 重新加载数据
	getDB().First(&tag, tagID)
	
	return &tag, nil
}

// DeleteTag 删除标签
func DeleteTag(tagID int64) error {
	// 检查标签是否正在被使用
	var postCount int64
	getDB().Model(&models.Post{}).
		Where("JSON_CONTAINS(tags, ?)", `"`+string(tagID)+`"`).
		Count(&postCount)
	
	if postCount > 0 {
		return gorm.ErrInvalidTransaction // 标签正在使用中，不能删除
	}
	
	// 删除标签（软删除：更新状态）
	result := getDB().Model(&models.Tag{}).
		Where("id = ?", tagID).
		Update("status", 1) // 标记为禁用
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	
	return nil
}

// GetHotTags 获取热门标签
func GetHotTags(limit int) ([]models.Tag, error) {
	var tags []models.Tag
	
	err := getDB().Where("status = 0").
		Order("usage_count DESC, last_used_at DESC").
		Limit(limit).
		Find(&tags).Error
	
	return tags, err
}

// SearchTags 搜索标签
func SearchTags(keyword string, limit int) ([]models.Tag, error) {
	var tags []models.Tag
	
	err := getDB().Where("status = 0 AND name LIKE ?", "%"+keyword+"%").
		Order("usage_count DESC, created_at DESC").
		Limit(limit).
		Find(&tags).Error
	
	return tags, err
}

// GetTagPosts 获取标签相关帖子
func GetTagPosts(tagName, userID string, page, pageSize int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64
	
	offset := (page - 1) * pageSize
	
	// 查询包含该标签的帖子
	query := getDB().Model(&models.Post{}).
		Where("status = 0 AND JSON_CONTAINS(tags, ?)", `"`+tagName+`"`)
	
	// 根据可见性过滤
	if userID != "" {
		// 获取好友ID
		friendIDs := GetFriendIDs(userID)
		query = query.Where("(visibility = 0 OR (visibility = 1 AND user_id IN (?)) OR (visibility = 2 AND user_id = ?))", 
			append(friendIDs, userID), userID)
	} else {
		// 只显示公开帖子
		query = query.Where("visibility = 0")
	}
	
	// 获取总数
	query.Count(&total)
	
	// 获取帖子列表
	err := query.Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&posts).Error
	
	return posts, total, err
}

// UpdateTagUsage 更新标签使用统计
func UpdateTagUsage(tagName string) error {
	// 查找标签
	var tag models.Tag
	err := getDB().Where("name = ?", tagName).First(&tag).Error
	if err != nil {
		// 如果标签不存在，自动创建
		if err == gorm.ErrRecordNotFound {
			newTag := models.Tag{
				Name:        tagName,
				Color:       "#" + GenerateRandomColor(),
				Status:      0,
				UsageCount:  1,
				LastUsedAt:  &time.Time{},
				CreatedAt:   time.Now(),
			}
			now := time.Now()
			newTag.LastUsedAt = &now
			return getDB().Create(&newTag).Error
		}
		return err
	}
	
	// 更新使用统计
	now := time.Now()
	return getDB().Model(&tag).Updates(map[string]interface{}{
		"usage_count":  gorm.Expr("usage_count + ?", 1),
		"last_used_at": &now,
	}).Error
}