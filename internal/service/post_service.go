package service

import (
	"encoding/json"
	"github.com/Yw332/campus-moments-go/internal/models"
	"gorm.io/gorm"
	"time"
)

// CreatePost 创建帖子
func CreatePost(userID, title, content string, images []string, video string, visibility int, tags []string) (*models.Post, error) {
	post := &models.Post{
		UserID:     userID,
		Title:       title,
		Content:     content,
		Visibility:  visibility,
		Status:      0, // 正常状态
		ViewCount:   0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	// 处理图片
	if len(images) > 0 {
		imagesJSON, _ := json.Marshal(images)
		post.Images = imagesJSON
	}
	
	// 处理视频
	if video != "" {
		post.Video = video
	}
	
	// 处理标签
	if len(tags) > 0 {
		tagsJSON, _ := json.Marshal(tags)
		post.Tags = tagsJSON
	}
	
	// 保存到数据库
	if err := getDB().Create(post).Error; err != nil {
		return nil, err
	}
	
	// 关联用户信息
	var user models.User
	if err := getDB().First(&user, "id = ?", userID).Error; err == nil {
		post.User = &user
	}
	
	// 更新用户发帖数
	getDB().Model(&models.User{}).Where("id = ?", userID).Update("post_count", gorm.Expr("post_count + ?", 1))
	
	// 更新标签使用统计
	for _, tagName := range tags {
		updateTagUsage(tagName)
	}
	
	return post, nil
}

// GetHomePagePosts 获取主页帖子（公开和好友帖子）
func GetHomePagePosts(userID string, page, pageSize int) ([]models.Post, int64, error) {
	return GetPostList(userID, "all", page, pageSize)
}

// GetPostList 获取帖子列表
func GetPostList(userID, visibility string, page, pageSize int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64
	
	offset := (page - 1) * pageSize
	
	query := getDB().Model(&models.Post{}).Where("status = ?", 0)
	
	// 根据可见性过滤
	if visibility == "0" {
		// 公开帖子
		query = query.Where("visibility = ?", 0)
	} else if visibility == "1" && userID != "" {
		// 好友可见的帖子（需要检查好友关系）
		friendIDs := GetFriendIDs(userID)
		query = query.Where("visibility IN (0, 1) AND (user_id = ? OR user_id IN (?))", userID, friendIDs)
	} else if visibility == "2" && userID != "" {
		// 包括私有帖子（只看自己的）
		query = query.Where("user_id = ? OR visibility = ?", userID, 0)
	} else if visibility == "all" && userID != "" {
		// 所有可见帖子：公开 + 好友帖子 + 自己的帖子
		friendIDs := GetFriendIDs(userID)
		query = query.Where("visibility = ? OR (visibility = ? AND user_id IN (?)) OR user_id = ?", 0, 1, friendIDs, userID)
	} else if visibility == "all" {
		// 未登录用户只能看公开帖子
		query = query.Where("visibility = ?", 0)
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

// GetPostDetail 获取帖子详情
func GetPostDetail(postID int64, userID string) (*models.Post, error) {
	var post models.Post
	
	err := getDB().Preload("User").
		First(&post, "id = ? AND status = ?", postID, 0).Error
	
	if err != nil {
		return nil, err
	}
	
	// 检查可见性
	if post.Visibility == 2 && post.UserID != userID {
		return nil, gorm.ErrRecordNotFound // 私密帖子只有作者可见
	}
	
	if post.Visibility == 1 && post.UserID != userID {
		// 好友可见，需要检查好友关系
		if !IsFriend(userID, post.UserID) {
			return nil, gorm.ErrRecordNotFound
		}
	}
	
	return &post, nil
}

// UpdatePost 更新帖子
func UpdatePost(postID int64, userID, title, content string, images []string, video string, visibility int, tags []string) (*models.Post, error) {
	var post models.Post
	
	// 检查帖子是否存在且属于当前用户
	if err := getDB().First(&post, "id = ? AND user_id = ? AND status = ?", postID, userID, 0).Error; err != nil {
		return nil, err
	}
	
	// 更新字段
	post.Title = title
	post.Content = content
	post.Visibility = visibility
	post.UpdatedAt = time.Now()
	
	// 更新图片
	if len(images) > 0 {
		imagesJSON, _ := json.Marshal(images)
		post.Images = imagesJSON
	}
	
	// 更新视频
	post.Video = video
	
	// 更新标签
	if len(tags) > 0 {
		tagsJSON, _ := json.Marshal(tags)
		post.Tags = tagsJSON
	}
	
	// 保存更新
	if err := getDB().Save(&post).Error; err != nil {
		return nil, err
	}
	
	// 更新标签使用统计
	for _, tagName := range tags {
		updateTagUsage(tagName)
	}
	
	// 重新加载用户信息
	getDB().Preload("User").First(&post, postID)
	
	return &post, nil
}

// DeletePost 删除帖子
func DeletePost(postID int64, userID string) error {
	// 软删除：更新状态为删除
	result := getDB().Model(&models.Post{}).
		Where("id = ? AND user_id = ?", postID, userID).
		Update("status", 1)
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	
	// 更新用户发帖数
	getDB().Model(&models.User{}).Where("id = ?", userID).Update("post_count", gorm.Expr("post_count - ?", 1))
	
	return nil
}

// GetUserPosts 获取用户帖子列表
func GetUserPosts(currentUserID, targetUserID string, page, pageSize int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64
	
	offset := (page - 1) * pageSize
	
	query := getDB().Model(&models.Post{}).Where("user_id = ? AND status = ?", targetUserID, 0)
	
	// 如果不是查看自己的帖子，需要过滤可见性
	if currentUserID != targetUserID {
		// 只显示公开帖子
		query = query.Where("visibility = ?", 0)
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

// IncrementViewCount 增加浏览量
func IncrementViewCount(postID int64) error {
	return getDB().Model(&models.Post{}).
		Where("id = ?", postID).
		Update("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// 辅助函数：更新标签使用统计
func updateTagUsage(tagName string) {
	// 查找或创建标签
	var tag models.Tag
	result := getDB().FirstOrCreate(&tag, models.Tag{
		Name:      tagName,
		Color:     "#" + GenerateRandomColor(),
		Status:    0,
		CreatedAt: time.Now(),
	})
	
	if result.Error == nil {
		// 更新使用统计
		now := time.Now()
		getDB().Model(&tag).Updates(map[string]interface{}{
			"usage_count":  gorm.Expr("usage_count + ?", 1),
			"last_used_at": &now,
		})
	}
}