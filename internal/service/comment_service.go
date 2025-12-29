package service

import (
	"encoding/json"
	"github.com/Yw332/campus-moments-go/internal/models"
	"gorm.io/gorm"
	"time"
)

// CreateComment 创建评论
func CreateComment(postID int, userID, content string, replies []map[string]interface{}) (*models.Comment, error) {
	// 检查帖子是否存在
	var post models.Post
	if err := getDB().First(&post, "id = ? AND status = ?", postID, 0).Error; err != nil {
		return nil, err
	}
	
	comment := &models.Comment{
		PostID:    postID,
		UserID:    userID,
		Content:   content,
		Status:    0, // 正常状态
		IsAuthor:  userID == post.UserID, // 检查是否为作者
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// 处理回复
	if len(replies) > 0 {
		repliesJSON, _ := json.Marshal(replies)
		comment.Replies = repliesJSON
	}
	
	// 保存评论
	if err := getDB().Create(comment).Error; err != nil {
		return nil, err
	}
	
	// 关联用户信息
	var user models.User
	if err := getDB().First(&user, "id = ?", userID).Error; err == nil {
		comment.User = &user
	}
	
	// 更新帖子评论数
	getDB().Model(&models.Post{}).Where("id = ?", postID).Update("comment_count", gorm.Expr("comment_count + ?", 1))
	
	// 更新用户评论数
	getDB().Model(&models.User{}).Where("id = ?", userID).Update("comment_count", gorm.Expr("comment_count + ?", 1))
	
	return comment, nil
}

// GetCommentList 获取评论列表
func GetCommentList(postID int64, page, pageSize int) ([]models.Comment, int64, error) {
	var comments []models.Comment
	var total int64
	
	offset := (page - 1) * pageSize
	
	query := getDB().Model(&models.Comment{}).Where("post_id = ? AND status = ?", postID, 0)
	
	// 获取总数
	query.Count(&total)
	
	// 获取评论列表
	err := query.Preload("User").
		Order("created_at ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&comments).Error
	
	return comments, total, err
}

// UpdateComment 更新评论
func UpdateComment(commentID int64, userID, content string) (*models.Comment, error) {
	var comment models.Comment
	
	// 检查评论是否存在且属于当前用户
	if err := getDB().First(&comment, "id = ? AND user_id = ? AND status = ?", commentID, userID, 0).Error; err != nil {
		return nil, err
	}
	
	// 更新内容
	comment.Content = content
	comment.UpdatedAt = time.Now()
	
	// 保存更新
	if err := getDB().Save(&comment).Error; err != nil {
		return nil, err
	}
	
	// 重新加载用户信息
	getDB().Preload("User").First(&comment, commentID)
	
	return &comment, nil
}

// DeleteComment 删除评论
func DeleteComment(commentID int64, userID string) error {
	// 软删除：更新状态
	result := getDB().Model(&models.Comment{}).
		Where("id = ? AND user_id = ?", commentID, userID).
		Update("status", 1)
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	
	// 获取评论信息用于更新计数
	var comment models.Comment
	if err := getDB().First(&comment, "id = ?", commentID).Error; err == nil {
		// 更新帖子评论数
		getDB().Model(&models.Post{}).Where("id = ?", comment.PostID).Update("comment_count", gorm.Expr("GREATEST(comment_count - ?, 0)", 1))
		
		// 更新用户评论数
		getDB().Model(&models.User{}).Where("id = ?", userID).Update("comment_count", gorm.Expr("GREATEST(comment_count - ?, 0)", 1))
	}
	
	return nil
}

// ToggleLikeComment 点赞/取消点赞评论
func ToggleLikeComment(commentID int64, userID string) (bool, error) {
	// 检查评论是否存在
	var comment models.Comment
	if err := getDB().First(&comment, "id = ? AND status = ?", commentID, 0).Error; err != nil {
		return false, err
	}
	
	// 查找是否已点赞
	var like models.Like
	err := getDB().Where("user_id = ? AND target_type = 2 AND target_id = ?", userID, commentID).First(&like).Error
	
	if err == gorm.ErrRecordNotFound {
		// 没有点赞记录，添加点赞
		newLike := models.Like{
			UserID:     userID,
			TargetType: 2, // 评论
			TargetID:   commentID,
			CreatedAt:  time.Now(),
		}
		
		if err := getDB().Create(&newLike).Error; err != nil {
			return false, err
		}
		
		// 更新评论点赞数
		getDB().Model(&comment).Update("like_count", gorm.Expr("like_count + ?", 1))
		
		return true, nil
	} else if err == nil {
		// 已有点赞记录，删除点赞
		if err := getDB().Delete(&like).Error; err != nil {
			return false, err
		}
		
		// 更新评论点赞数
		getDB().Model(&comment).Update("like_count", gorm.Expr("GREATEST(like_count - ?, 0)", 1))
		
		return false, nil
	}
	
	return false, err
}

// ReplyComment 回复评论
func ReplyComment(commentID int64, userID, content string) (*models.Comment, error) {
	// 获取原评论
	var parentComment models.Comment
	if err := getDB().First(&parentComment, "id = ? AND status = ?", commentID, 0).Error; err != nil {
		return nil, err
	}
	
	// 创建回复评论
	reply := &models.Comment{
		PostID:    parentComment.PostID,
		UserID:    userID,
		Content:   content,
		Status:    0,
		IsAuthor:  userID == parentComment.UserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// 保存回复
	if err := getDB().Create(reply).Error; err != nil {
		return nil, err
	}
	
	// 更新原评论的回复列表
	var replies []map[string]interface{}
	if parentComment.Replies != nil {
		json.Unmarshal(parentComment.Replies, &replies)
	}
	
	newReply := map[string]interface{}{
		"id":         reply.ID,
		"userId":     reply.UserID,
		"content":    reply.Content,
		"isAuthor":   reply.IsAuthor,
		"createdAt":  reply.CreatedAt,
	}
	
	replies = append(replies, newReply)
	repliesJSON, _ := json.Marshal(replies)
	
	// 更新原评论
	getDB().Model(&parentComment).Updates(map[string]interface{}{
		"replies":    repliesJSON,
		"updated_at": time.Now(),
	})
	
	// 关联用户信息
	var user models.User
	if err := getDB().First(&user, "id = ?", userID).Error; err == nil {
		reply.User = &user
	}
	
	// 更新帖子评论数
	getDB().Model(&models.Post{}).Where("id = ?", parentComment.PostID).Update("comment_count", gorm.Expr("comment_count + ?", 1))
	
	// 更新用户评论数
	getDB().Model(&models.User{}).Where("id = ?", userID).Update("comment_count", gorm.Expr("comment_count + ?", 1))
	
	return reply, nil
}