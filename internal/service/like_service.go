package service

import (
	"encoding/json"
	"github.com/Yw332/campus-moments-go/internal/models"
	"gorm.io/gorm"
	"time"
)

// ToggleLikePost 点赞/取消点赞帖子
func ToggleLikePost(postID int64, userID string) (bool, error) {
	// 检查帖子是否存在
	var post models.Post
	if err := getDB().First(&post, "id = ? AND status = ?", postID, 0).Error; err != nil {
		return false, err
	}
	
	// 查找是否已点赞
	var like models.Like
	err := getDB().Where("user_id = ? AND target_type = 1 AND target_id = ?", userID, postID).First(&like).Error
	
	if err == gorm.ErrRecordNotFound {
		// 没有点赞记录，添加点赞
		newLike := models.Like{
			UserID:     userID,
			TargetType: 1, // 帖子
			TargetID:   postID,
			CreatedAt:  time.Now(),
		}
		
		if err := getDB().Create(&newLike).Error; err != nil {
			return false, err
		}
		
		// 更新帖子点赞数
		getDB().Model(&post).Update("like_count", gorm.Expr("like_count + ?", 1))
		
		// 更新用户获赞数
		getDB().Model(&models.User{}).Where("id = ?", post.UserID).Update("like_count", gorm.Expr("like_count + ?", 1))
		
		// 更新帖子的点赞用户列表
		updatePostLikedUsers(postID, userID, true)
		
		return true, nil
	} else if err == nil {
		// 已有点赞记录，删除点赞
		if err := getDB().Delete(&like).Error; err != nil {
			return false, err
		}
		
		// 更新帖子点赞数
		getDB().Model(&post).Update("like_count", gorm.Expr("GREATEST(like_count - ?, 0)", 1))
		
		// 更新用户获赞数
		getDB().Model(&models.User{}).Where("id = ?", post.UserID).Update("like_count", gorm.Expr("GREATEST(like_count - ?, 0)", 1))
		
		// 更新帖子的点赞用户列表
		updatePostLikedUsers(postID, userID, false)
		
		return false, nil
	}
	
	return false, err
}

// GetPostLikes 获取帖子点赞列表
func GetPostLikes(postID int64, page, pageSize int) ([]models.Like, int64, error) {
	var likes []models.Like
	var total int64
	
	offset := (page - 1) * pageSize
	
	query := getDB().Model(&models.Like{}).Where("target_type = 1 AND target_id = ?", postID)
	
	// 获取总数
	query.Count(&total)
	
	// 获取点赞列表
	err := query.Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&likes).Error
	
	return likes, total, err
}

// GetCommentLikes 获取评论点赞列表
func GetCommentLikes(commentID int64, page, pageSize int) ([]models.Like, int64, error) {
	var likes []models.Like
	var total int64
	
	offset := (page - 1) * pageSize
	
	query := getDB().Model(&models.Like{}).Where("target_type = 2 AND target_id = ?", commentID)
	
	// 获取总数
	query.Count(&total)
	
	// 获取点赞列表
	err := query.Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&likes).Error
	
	return likes, total, err
}

// GetUserLikes 获取用户点赞列表
func GetUserLikes(userID, targetType string, page, pageSize int) ([]models.Like, int64, error) {
	var likes []models.Like
	var total int64

	offset := (page - 1) * pageSize

	query := getDB().Model(&models.Like{}).Where("user_id = ?", userID)

	if targetType != "" {
		query = query.Where("target_type = ?", targetType)
	}

	// 获取总数
	query.Count(&total)

	// 获取点赞列表,预加载User和Post信息
	err := query.Preload("User").
		Preload("Post").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&likes).Error

	return likes, total, err
}

// updatePostLikedUsers 更新帖子的点赞用户列表
func updatePostLikedUsers(postID int64, userID string, isAdd bool) {
	var post models.Post
	if err := getDB().First(&post, "id = ?", postID).Error; err != nil {
		return
	}
	
	var likedUsers []string
	if post.LikedUsers != nil {
		json.Unmarshal(post.LikedUsers, &likedUsers)
	}
	
	if isAdd {
		// 添加用户到点赞列表
		found := false
		for _, uid := range likedUsers {
			if uid == userID {
				found = true
				break
			}
		}
		if !found {
			likedUsers = append(likedUsers, userID)
		}
	} else {
		// 从点赞列表中移除用户
		newLikedUsers := make([]string, 0, len(likedUsers))
		for _, uid := range likedUsers {
			if uid != userID {
				newLikedUsers = append(newLikedUsers, uid)
			}
		}
		likedUsers = newLikedUsers
	}
	
	likedUsersJSON, _ := json.Marshal(likedUsers)
	getDB().Model(&post).Update("liked_users", likedUsersJSON)
}