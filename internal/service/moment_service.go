package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Yw332/campus-moments-go/internal/models"
	"github.com/Yw332/campus-moments-go/pkg/database"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// MomentService 动态服务
type MomentService struct {
	db *gorm.DB
}

// NewMomentService 创建动态服务实例
func NewMomentService() *MomentService {
	return &MomentService{
		db: database.GetDB(),
	}
}

// CreateMomentRequest 创建动态请求
type CreateMomentRequest struct {
	Content    string             `json:"content" binding:"required"`
	Tags       []string           `json:"tags"`
	Media      []models.MediaItem `json:"media"`
	Visibility int                `json:"visibility"` // 0公开/1好友/2私密
}

// UpdateMomentRequest 更新动态请求
type UpdateMomentRequest struct {
	Content    *string            `json:"content,omitempty"`
	Tags       []string           `json:"tags,omitempty"`
	Media      []models.MediaItem `json:"media,omitempty"`
	Visibility *int               `json:"visibility,omitempty"`
}

// CreateMoment 创建动态
func (s *MomentService) CreateMoment(userID string, req *CreateMomentRequest) (*models.Moment, error) {
	if s.db == nil {
		s.db = database.GetDB()
	}
	if s.db == nil {
		return nil, errors.New("数据库未连接")
	}

	// 从media中提取图片URL，存储到images字段
	var images []string
	for _, item := range req.Media {
		if item.Type == "image" && item.URL != "" {
			images = append(images, item.URL)
		}
	}
	fmt.Printf("调试: CreateMoment - 提取的images: %v\n", images)

	moment := &models.Moment{
		UserID:          userID,
		AuthorID:        userID,
		Title:           req.Content, // 使用内容作为标题
		Content:         req.Content,
		Tags:           models.Tags(req.Tags),
		Images:          models.Tags(images), // 将media中的图片转换为images字段
		CommentsSummary:  "[]", // 空JSON字符串
		LikedUsers:      models.Tags{}, // 空JSON数组
		Visibility:      req.Visibility,
		LikeCount:       0,
		CommentCount:    0,
		Status:          1, // 正常状态
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.db.Create(moment).Error; err != nil {
		return nil, fmt.Errorf("创建动态失败: %w", err)
	}

	// 加载作者信息
	if err := s.db.Preload("Author").First(moment, moment.ID).Error; err != nil {
		return nil, fmt.Errorf("加载作者信息失败: %w", err)
	}

	return moment, nil
}

// GetMomentByID 根据ID获取动态详情 - 修复JSON字段问题
func (s *MomentService) GetMomentByID(id int64) (*models.Moment, error) {
	// 使用GORM但排除有问题的字段
	var result struct {
		ID           int64     `gorm:"column:id"`
		UserID       string    `gorm:"column:user_id"`
		AuthorID     string    `gorm:"column:author_id"`
		Title        string    `gorm:"column:title"`
		Content      string    `gorm:"column:content"`
		Images       string    `gorm:"column:images"`
		Video        string    `gorm:"column:video"`
		Visibility   int       `gorm:"column:visibility"`
		Status       int       `gorm:"column:status"`
		LikeCount    int64     `gorm:"column:like_count"`
		CommentCount int64     `gorm:"column:comment_count"`
		ViewCount    int64     `gorm:"column:view_count"`
		CreatedAt    time.Time `gorm:"column:created_at"`
		UpdatedAt    time.Time `gorm:"column:updated_at"`
	}

	query := `SELECT id, user_id, author_id, title, content, images, video, 
	              visibility, status, like_count, comment_count, view_count, created_at, updated_at
	              FROM posts WHERE id = ? AND status = 1`

	if err := database.GetDB().Raw(query, id).Scan(&result).Error; err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("动态不存在")
		}
		return nil, fmt.Errorf("查询动态失败: %v", err)
	}

	// 手动构建Moment对象
	moment := &models.Moment{
		ID:           result.ID,
		UserID:       result.UserID,
		AuthorID:     result.AuthorID,
		Title:        result.Title,
		Content:      result.Content,
		Visibility:   result.Visibility,
		Status:       result.Status,
		LikeCount:    result.LikeCount,
		CommentCount: result.CommentCount,
		ViewCount:    result.ViewCount,
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
	}

	// 手动处理images字段
	if result.Images != "" && result.Images != "null" {
		json.Unmarshal([]byte(result.Images), &moment.Images)
	}

	// 处理video字段
	if result.Video != "" && result.Video != "null" {
		moment.Video = &result.Video
	}

	// 查询作者信息
	var user models.User
	database.GetDB().Where("user_id = ?", result.UserID).First(&user)
	moment.Author = &user

	return moment, nil
}

// ListMoments 获取动态列表（支持分页）
func (s *MomentService) ListMoments(page, pageSize int, userID *string) ([]models.Moment, int64, error) {
	if s.db == nil {
		s.db = database.GetDB()
	}
	if s.db == nil {
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

	query := s.db.Preload("Author").Where("status = 1")

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
	if s.db == nil {
		s.db = database.GetDB()
	}
	if s.db == nil {
		return nil, errors.New("数据库未连接")
	}

	// 首先检查动态是否存在且属于当前用户
	var moment models.Moment
	err := s.db.Where("id = ? AND author_id = ? AND status = 1", momentID, userID).First(&moment).Error
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
	if err := s.db.Model(&moment).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("更新动态失败: %w", err)
	}

	// 重新加载完整数据
	if err := s.db.Preload("Author").First(&moment, moment.ID).Error; err != nil {
		return nil, fmt.Errorf("重新加载动态失败: %w", err)
	}

	return &moment, nil
}

// DeleteMoment 删除动态（软删除）
func (s *MomentService) DeleteMoment(userID string, momentID int64) error {
	if s.db == nil {
		s.db = database.GetDB()
	}
	if s.db == nil {
		return errors.New("数据库未连接")
	}

	// 检查动态是否存在且属于当前用户
	var moment models.Moment
	err := s.db.Where("id = ? AND author_id = ? AND status = 1", momentID, userID).First(&moment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("动态不存在或无权限删除")
		}
		return fmt.Errorf("查询动态失败: %w", err)
	}

	// 软删除：更新状态为2（删除）
	if err := s.db.Model(&moment).Updates(map[string]interface{}{
		"status":     2,
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
