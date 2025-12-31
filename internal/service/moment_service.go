package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Yw332/campus-moments-go/internal/models"
	"github.com/Yw332/campus-moments-go/pkg/database"
	"gorm.io/gorm"
)

// MomentService 动态服务
type MomentService struct{}

// NewMomentService 创建动态服务实例
func NewMomentService() *MomentService {
	return &MomentService{}
}

// getDB 获取数据库连接
func (s *MomentService) getDB() *gorm.DB {
	return database.GetDB()
}

// CreateMomentRequest 创建动态请求
type CreateMomentRequest struct {
	Title      string                  `json:"title"`       // 标题（可选）
	Content    string                  `json:"content" binding:"required"`
	Tags       []string                `json:"tags"`
	Images     []string                `json:"images"`      // 图片URL数组（前端格式）
	Media      []models.MediaItem      `json:"media"`       // 媒体项（后端格式）
	Visibility int                     `json:"visibility"`  // 0公开/1好友/2私密
}

// UpdateMomentRequest 更新动态请求
type UpdateMomentRequest struct {
	Content    *string                 `json:"content,omitempty"`
	Tags       []string                `json:"tags,omitempty"`
	Media      []models.MediaItem      `json:"media,omitempty"`
	Visibility *int                    `json:"visibility,omitempty"`
}

// CreateMoment 创建动态
func (s *MomentService) CreateMoment(userID string, req *CreateMomentRequest) (*models.Moment, error) {
	db := s.getDB()
	if db == nil {
		return nil, errors.New("数据库未连接")
	}

	// 转换Tags到JSON格式
	tagsJSON, _ := json.Marshal(models.Tags(req.Tags))

	// 处理图片：优先使用Images数组，如果没有则从Media中提取
	var imagesJSON json.RawMessage
	if len(req.Images) > 0 {
		// 使用前端传来的Images数组
		imagesJSON, _ = json.Marshal(req.Images)
	} else if len(req.Media) > 0 {
		// 从Media中提取图片URL
		var imageURLs []string
		for _, media := range req.Media {
			if media.Type == "image" {
				imageURLs = append(imageURLs, media.URL)
			}
		}
		if len(imageURLs) > 0 {
			imagesJSON, _ = json.Marshal(imageURLs)
		}
	}

	moment := &models.Moment{
		UserID:     userID,
		Title:      req.Title,
		Content:    req.Content,
		Images:     imagesJSON,
		Tags:       tagsJSON,
		Visibility: req.Visibility,
		LikeCount:  0,
		CommentCount: 0,
		Status:     0, // 正常状态（根据数据库结构：0-正常 1-删除）
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := db.Create(moment).Error; err != nil {
		return nil, fmt.Errorf("创建动态失败: %w", err)
	}

	// 加载作者信息
	if err := db.Preload("User").First(moment, moment.ID).Error; err != nil {
		return nil, fmt.Errorf("加载作者信息失败: %w", err)
	}

	return moment, nil
}

// GetMomentByID 根据ID获取动态详情
func (s *MomentService) GetMomentByID(id int64) (*models.Moment, error) {
	db := s.getDB()
	if db == nil {
		return nil, errors.New("数据库未连接")
	}

	var moment models.Moment
	err := db.Preload("User").Where("id = ? AND status = 0", id).First(&moment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("动态不存在")
		}
		return nil, fmt.Errorf("查询动态失败: %w", err)
	}

	return &moment, nil
}

// ListMoments 获取动态列表（支持分页）
func (s *MomentService) ListMoments(page, pageSize int, userID *string) ([]models.Moment, int64, error) {
	db := s.getDB()
	if db == nil {
		return nil, 0, errors.New("数据库未连接")
	}

	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // 限制最大分页
	}

	offset := (page - 1) * pageSize

	var moments []models.Moment
	var total int64

	query := db.Preload("User").Where("status = 0")

	// 如果指定了用户ID，则查询该用户的动态
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	// 查询总数
	if err := query.Model(&models.Moment{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询总数失败: %w", err)
	}

	// 查询列表
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&moments).Error; err != nil {
		return nil, 0, fmt.Errorf("查询动态列表失败: %w", err)
	}

	return moments, total, nil
}

// UpdateMoment 更新动态
func (s *MomentService) UpdateMoment(userID string, momentID int64, req *UpdateMomentRequest) (*models.Moment, error) {
	db := s.getDB()
	if db == nil {
		return nil, errors.New("数据库未连接")
	}

	// 首先检查动态是否存在且属于当前用户
	var moment models.Moment
	err := db.Where("id = ? AND user_id = ? AND status = 0", momentID, userID).First(&moment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("动态不存在或无权限修改")
		}
		return nil, fmt.Errorf("查询动态失败: %w", err)
	}

	// 更新字段
	updates := make(map[string]interface{})
	updates["updated_at"] = time.Now()

	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Tags != nil {
		updates["tags"] = models.Tags(req.Tags)
	}
	if req.Media != nil {
		updates["media"] = models.MediaItems(req.Media)
	}
	if req.Visibility != nil {
		updates["visibility"] = *req.Visibility
	}

	// 执行更新
	if err := db.Model(&moment).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("更新动态失败: %w", err)
	}

	// 重新加载完整数据
	if err := db.Preload("User").First(&moment, moment.ID).Error; err != nil {
		return nil, fmt.Errorf("重新加载动态失败: %w", err)
	}

	return &moment, nil
}

// DeleteMoment 删除动态（软删除）
func (s *MomentService) DeleteMoment(userID string, momentID int64) error {
	db := s.getDB()
	if db == nil {
		return errors.New("数据库未连接")
	}

	// 检查动态是否存在且属于当前用户
	var moment models.Moment
	err := db.Where("id = ? AND user_id = ? AND status = 0", momentID, userID).First(&moment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("动态不存在或无权限删除")
		}
		return fmt.Errorf("查询动态失败: %w", err)
	}

	// 软删除：更新状态为1（删除）
	if err := db.Model(&moment).Updates(map[string]interface{}{
		"status":     1,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("删除动态失败: %w", err)
	}

	return nil
}

// AdminDeleteMoment 管理员删除动态（可删除任意用户的动态）
func (s *MomentService) AdminDeleteMoment(momentID int64) error {
	db := s.getDB()
	if db == nil {
		return errors.New("数据库未连接")
	}

	// 检查动态是否存在（status=0表示正常）
	var moment models.Moment
	err := db.Where("id = ? AND status = 0", momentID).First(&moment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("动态不存在")
		}
		return fmt.Errorf("查询动态失败: %w", err)
	}

	// 软删除：更新状态为1（已删除）
	if err := db.Model(&moment).Updates(map[string]interface{}{
		"status":     1,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("删除动态失败: %w", err)
	}

	return nil
}

// GetUserMoments 获取用户的所有动态
func (s *MomentService) GetUserMoments(userID string, page, pageSize int) ([]models.Moment, int64, error) {
	return s.ListMoments(page, pageSize, &userID)
}