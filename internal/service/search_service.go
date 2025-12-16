package service

import (
	"fmt"

	"github.com/Yw332/campus-moments-go/internal/models"
	"github.com/Yw332/campus-moments-go/pkg/database"
	"gorm.io/gorm"
)

type SearchService struct {
	db *gorm.DB
}

func NewSearchService() *SearchService {
	return &SearchService{
		db: database.GetDB(),
	}
}

type SearchRequest struct {
	Keyword  string `json:"keyword" binding:"required"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

type SearchResponse struct {
	Moments    []models.Moment `json:"moments"`
	Users      []models.User   `json:"users"`
	Pagination Pagination     `json:"pagination"`
}

type FilterRequest struct {
	Page       int      `json:"page"`
	PageSize   int      `json:"pageSize"`
	Visibility string   `json:"visibility"`
	Tags       []string `json:"tags"`
	StartDate  string   `json:"startDate"`
	EndDate    string   `json:"endDate"`
}

// Pagination 分页信息
type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
}

// SearchContent 搜索内容
func (s *SearchService) SearchContent(keyword string, page, pageSize int) (*SearchResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var moments []models.Moment
	var totalMoments int64
	var users []models.User
	var totalUsers int64

	// 搜索动态
	momentQuery := s.db.Where("content LIKE ? AND status = 1", "%"+keyword+"%").
		Preload("Author").
		Order("created_at DESC")

	// 统计动态总数
	momentQuery.Count(&totalMoments)

	// 分页查询动态
	if err := momentQuery.Offset(offset).Limit(pageSize).Find(&moments).Error; err != nil {
		return nil, fmt.Errorf("搜索动态失败: %v", err)
	}

	// 搜索用户
	userQuery := s.db.Where("username LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	// 统计用户总数
	userQuery.Count(&totalUsers)

	// 分页查询用户
	if err := userQuery.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("搜索用户失败: %v", err)
	}

	return &SearchResponse{
		Moments: moments,
		Users:   users,
		Pagination: Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    totalMoments + totalUsers, // 总数只是示例
		},
	}, nil
}

// GetHotWords 获取热词
func (s *SearchService) GetHotWords() ([]string, error) {
	// 这里简化实现，实际应该基于搜索历史统计
	hotWords := []string{
		"校园",
		"活动",
		"学习",
		"美食",
		"运动",
		"兼职",
		"考试",
		"社团",
		"室友",
		"考研",
	}
	return hotWords, nil
}

// GetSearchHistory 获取搜索历史
func (s *SearchService) GetSearchHistory(userID string) ([]string, error) {
	// 简化实现，实际应该有搜索历史表
	// 这里返回一些模拟数据
	history := []string{
		"校园活动",
		"学习方法",
		"美食推荐",
	}
	return history, nil
}

// SaveSearchHistory 保存搜索历史
func (s *SearchService) SaveSearchHistory(userID, keyword string) error {
	// 简化实现，实际应该保存到搜索历史表
	// 这里只是记录日志
	fmt.Printf("用户 %s 搜索了: %s\n", userID, keyword)
	return nil
}

// GetFilteredContent 获取筛选内容
func (s *SearchService) GetFilteredContent(filter *FilterRequest) (*SearchResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 || filter.PageSize > 100 {
		filter.PageSize = 10
	}

	offset := (filter.Page - 1) * filter.PageSize

	query := s.db.Where("status = 1").Preload("Author")

	// 按可见性筛选
	if filter.Visibility != "" {
		query = query.Where("visibility = ?", filter.Visibility)
	}

	// 按标签筛选
	if len(filter.Tags) > 0 {
		for _, tag := range filter.Tags {
			query = query.Where("JSON_CONTAINS(tags, ?)", fmt.Sprintf(`"%s"`, tag))
		}
	}

	// 按时间范围筛选
	if filter.StartDate != "" {
		query = query.Where("created_at >= ?", filter.StartDate)
	}
	if filter.EndDate != "" {
		query = query.Where("created_at <= ?", filter.EndDate+" 23:59:59")
	}

	var moments []models.Moment
	var total int64

	// 统计总数
	query.Count(&total)

	// 分页查询
	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.PageSize).Find(&moments).Error; err != nil {
		return nil, fmt.Errorf("获取筛选内容失败: %v", err)
	}

	return &SearchResponse{
		Moments: moments,
		Users:   []models.User{},
		Pagination: Pagination{
			Page:     filter.Page,
			PageSize: filter.PageSize,
			Total:    total,
		},
	}, nil
}

// GetTrendingKeywords 获取趋势关键词
func (s *SearchService) GetTrendingKeywords(days int) ([]string, error) {
	// 简化实现，实际应该基于搜索历史和内容分析
	trendingKeywords := []string{
		"期末考试",
		"寒假安排",
		"社团招新",
		"实习机会",
	}
	return trendingKeywords, nil
}

// GetSearchSuggestions 获取搜索建议
func (s *SearchService) GetSearchSuggestions(keyword string) ([]string, error) {
	if keyword == "" {
		return s.GetHotWords()
	}

	var suggestions []string
	// 从用户名和内容中获取建议
	if err := s.db.Model(&models.User{}).
		Where("username LIKE ?", "%"+keyword+"%").
		Limit(5).
		Pluck("username", &suggestions).Error; err != nil {
		return nil, fmt.Errorf("获取搜索建议失败: %v", err)
	}

	return suggestions, nil
}