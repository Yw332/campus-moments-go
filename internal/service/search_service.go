package service

import (
	"fmt"
	"time"

	"github.com/Yw332/campus-moments-go/internal/models"
	"github.com/Yw332/campus-moments-go/pkg/database"
)

type SearchService struct {
}

func NewSearchService() *SearchService {
	return &SearchService{}
}

type SearchRequest struct {
	Keyword  string `json:"keyword" binding:"required"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

type SearchResponse struct {
	Moments    []models.Moment `json:"moments"`
	Users      []models.User   `json:"users"`
	Pagination Pagination      `json:"pagination"`
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
func (s *SearchService) SearchContent(keyword string, page, pageSize int, sortBy string) (*SearchResponse, error) {
	// 获取数据库连接
	db := database.GetDB()
	if db == nil {
		return &SearchResponse{
			Moments: []models.Moment{},
			Users:   []models.User{},
			Pagination: Pagination{
				Page:     page,
				PageSize: pageSize,
				Total:    0,
			},
		}, nil
	}
	


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
	momentQuery := db.Table("moments").Where("content LIKE ?", "%"+keyword+"%")

	// 根据排序方式设置排序
	switch sortBy {
	case "hottest":
		// 最热：按点赞数和评论数综合排序
		momentQuery = momentQuery.Order("(like_count + comment_count * 2) DESC, created_at DESC")
	case "comprehensive":
		// 综合：考虑时间、点赞、评论数
		momentQuery = momentQuery.Order("(like_count * 3 + comment_count * 2 + TIMESTAMPDIFF(HOUR, created_at, NOW()) / 24) DESC")
	case "latest", "":
		// 最新：按时间排序（默认）
		momentQuery = momentQuery.Order("created_at DESC")
	default:
		momentQuery = momentQuery.Order("created_at DESC")
	}

	// 统计动态总数
	momentQuery.Count(&totalMoments)

	// 分页查询动态
	if err := momentQuery.Offset(offset).Limit(pageSize).Find(&moments).Error; err != nil {
		return nil, fmt.Errorf("搜索动态失败: %v", err)
	}

	// 搜索用户
	userQuery := db.Table("users").Where("username LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%")

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
			Total:    totalMoments + totalUsers,
		},
	}, nil
}

// GetHotWords 获取热词
func (s *SearchService) GetHotWords() ([]string, error) {
	// 获取数据库连接
	db := database.GetDB()
	if db == nil {
		// 返回默认热词
		return []string{
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
		}, nil
	}

	// 基于今日搜索历史统计热词
	var hotWords []string
	
	// 获取今日搜索最多的关键词
	err := db.Table("search_history").
		Select("keyword, COUNT(*) as count").
		Where("created_at >= DATE_FORMAT(NOW(), '%Y-%m-%d')").
		Group("keyword").
		Order("count DESC").
		Limit(10).
		Pluck("keyword", &hotWords).Error
	
	// 如果没有搜索历史，返回默认热词
	if err != nil || len(hotWords) == 0 {
		hotWords = []string{
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
	}
	
	return hotWords, nil
}

// GetSearchHistory 获取搜索历史
func (s *SearchService) GetSearchHistory(userID string) ([]string, error) {
	// 获取数据库连接
	db := database.GetDB()
	if db == nil {
		// 返回默认空列表
		return []string{}, nil
	}

	var keywords []string
	
	// 获取用户最近30天的搜索历史，去重并按时间倒序
	type HistoryResult struct {
		Keyword string
		LatestTime time.Time
	}
	var results []HistoryResult
	
	err := db.Table("search_history").
		Select("keyword, MAX(created_at) as latest_time").
		Where("user_id = ? AND created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)", userID).
		Group("keyword").
		Order("latest_time DESC").
		Limit(20).
		Find(&results).Error
		
	if err != nil {
		// 如果表不存在或其他错误，返回空列表而不是错误
		return []string{}, nil
	}
	
	// 提取关键词
	for _, result := range results {
		keywords = append(keywords, result.Keyword)
	}
	
	if err != nil {
		// 如果表不存在或其他错误，返回空列表而不是错误
		return []string{}, nil
	}
	
	return keywords, nil
}

// SaveSearchHistory 保存搜索历史
func (s *SearchService) SaveSearchHistory(userID, keyword string) error {
	// 获取数据库连接
	db := database.GetDB()
	if db == nil {
		// 记录日志但不报错
		fmt.Printf("数据库不可用，跳过保存搜索历史: %s\n", keyword)
		return nil
	}

	// 检查是否已存在相同的搜索记录，避免重复
	var existingHistory models.SearchHistory
	err := db.Table("search_history").Where("user_id = ? AND keyword = ? AND DATE(created_at) = CURDATE()", 
		userID, keyword).First(&existingHistory).Error
	
	if err == nil {
		// 今天已搜索过相同关键词，只更新时间
		return db.Table("search_history").Where("id = ?", existingHistory.ID).
			Update("created_at", time.Now()).Error
	}
	
	// 创建新的搜索历史记录
	history := models.SearchHistory{
		UserID:    userID,
		Keyword:   keyword,
		CreatedAt: time.Now(),
	}
	
	return db.Table("search_history").Create(&history).Error
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

	db := database.GetDB()
	if db == nil {
		return &SearchResponse{
			Moments: []models.Moment{},
			Users:   []models.User{},
			Pagination: Pagination{
				Page:     filter.Page,
				PageSize: filter.PageSize,
				Total:    0,
			},
		}, nil
	}

	query := db.Table("moments").Where("1=1")

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
	// 获取数据库连接
	db := database.GetDB()
	if db == nil {
		// 返回默认热词
		return []string{
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
		}, nil
	}

	if keyword == "" {
		return s.GetHotWords()
	}

	var suggestions []string
	// 从用户名和内容中获取建议
	if err := db.Table("users").
		Where("username LIKE ?", "%"+keyword+"%").
		Limit(5).
		Pluck("username", &suggestions).Error; err != nil {
		return nil, fmt.Errorf("获取搜索建议失败: %v", err)
	}

	return suggestions, nil
}
