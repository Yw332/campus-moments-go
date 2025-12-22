package service

import (
	"errors"
	"time"

	"github.com/Yw332/campus-moments-go/internal/models"
	"github.com/Yw332/campus-moments-go/pkg/database"
	"gorm.io/gorm"
)

type InteractionService struct{}

func NewInteractionService() *InteractionService {
	return &InteractionService{}
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	MomentID      int64  `json:"momentId" binding:"required"`
	Content       string `json:"content" binding:"required"`
	ParentID      *int64 `json:"parentId"`      // 回复评论时传
	ReplyToUserID *string `json:"replyToUserId"` // 被回复的用户ID
}

// CreateComment 发表评论
func (s *InteractionService) CreateComment(userID string, req *CreateCommentRequest) (*models.Comment, error) {
	db := database.GetDB()

	comment := &models.Comment{
		MomentID:      req.MomentID,
		UserID:        userID,
		Content:       req.Content,
		ParentID:      req.ParentID,
		ReplyToUserID: req.ReplyToUserID,
		Status:        1,
		CreatedAt:     time.Now(),
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. 创建评论
		if err := tx.Create(comment).Error; err != nil {
			return err
		}

		// 2. 更新动态的评论数
		if err := tx.Model(&models.Moment{}).Where("id = ?", req.MomentID).
			UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	
	// 加载用户信息以便返回
	db.Preload("User").Preload("ReplyToUser").First(comment, comment.ID)

	return comment, nil
}

// ToggleLike 点赞/取消点赞
func (s *InteractionService) ToggleLike(userID string, targetID int64, targetType int) (bool, error) {
	db := database.GetDB()
	var isLiked bool

	err := db.Transaction(func(tx *gorm.DB) error {
		// 检查是否已点赞
		var existingLike models.Like
		result := tx.Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).First(&existingLike)

		if result.Error == nil {
			// 已点赞 -> 取消点赞
			if err := tx.Delete(&existingLike).Error; err != nil {
				return err
			}
			isLiked = false

			// 减少计数
			if targetType == models.LikeTargetTypeMoment {
				if err := tx.Model(&models.Moment{}).Where("id = ?", targetID).
					UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error; err != nil {
					return err
				}
			} else if targetType == models.LikeTargetTypeComment {
				if err := tx.Model(&models.Comment{}).Where("id = ?", targetID).
					UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error; err != nil {
					return err
				}
			}

		} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 未点赞 -> 点赞
			newLike := models.Like{
				UserID:     userID,
				TargetID:   targetID,
				TargetType: targetType,
				CreatedAt:  time.Now(),
			}
			if err := tx.Create(&newLike).Error; err != nil {
				return err
			}
			isLiked = true

			// 增加计数
			if targetType == models.LikeTargetTypeMoment {
				if err := tx.Model(&models.Moment{}).Where("id = ?", targetID).
					UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
					return err
				}
			} else if targetType == models.LikeTargetTypeComment {
				if err := tx.Model(&models.Comment{}).Where("id = ?", targetID).
					UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
					return err
				}
			}
		} else {
			return result.Error
		}

		return nil
	})

	return isLiked, err
}

// GetMomentComments 获取动态的评论列表
func (s *InteractionService) GetMomentComments(momentID int64) ([]models.Comment, error) {
	db := database.GetDB()
	var comments []models.Comment

	// 这里可以做分页，暂时全量返回
	// 预加载用户信息和被回复用户信息
	err := db.Where("moment_id = ? AND status = 1", momentID).
		Preload("User").
		Preload("ReplyToUser").
		Order("created_at asc"). // 按时间正序
		Find(&comments).Error

	return comments, err
}

// DeleteComment 删除评论
func (s *InteractionService) DeleteComment(userID string, commentID int64) error {
	db := database.GetDB()

	var comment models.Comment
	if err := db.First(&comment, commentID).Error; err != nil {
		return errors.New("评论不存在")
	}

	// 权限检查：只能删除自己的评论
	if comment.UserID != userID {
		return errors.New("无权删除此评论")
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// 软删除
		if err := tx.Model(&comment).Update("status", 2).Error; err != nil {
			return err
		}

		// 减少动态评论数
		if err := tx.Model(&models.Moment{}).Where("id = ?", comment.MomentID).
			UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})
}
