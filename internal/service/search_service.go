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
	Posts      []models.Moment `json:"posts"`
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
func (s *SearchService) SearchContent(keyword string, page, pageSize int) (*SearchResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var posts []models.Moment
	var totalPosts int64
	var users []models.User
	var totalUsers int64

	// 检查数据库连接
	if s.db == nil {
		return &SearchResponse{
			Posts:      []models.Post{},
			Users:      []models.User{},
			Pagination: Pagination{Page: page, PageSize: pageSize, Total: 0},
		}, nil
	}

	// 搜索动态（使用Moment模型）
	postQuery := s.db.Model(&models.Moment{}).
		Where("(title LIKE ? OR content LIKE ?) AND status = ?", 
			"%"+keyword+"%", "%"+keyword+"%", 0).
		Preload("User").
		Order("created_at DESC")

	// 统计动态总数
	postQuery.Count(&totalPosts)

	// 分页查询动态
	if err := postQuery.Offset(offset).Limit(pageSize).Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("搜索动态失败: %v", err)
	}

	// 搜索用户
	userQuery := s.db.Model(&models.User{}).
		Where("username LIKE ? OR signature LIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	// 统计用户总数
	userQuery.Count(&totalUsers)

	// 分页查询用户
	if err := userQuery.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("搜索用户失败: %v", err)
	}

	return &SearchResponse{
		Posts: posts,
		Users: users,
		Pagination: Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    totalPosts + totalUsers, // 总数只是示例
		},
	}, nil
}

// GetHotWords 获取热词
func (s *SearchService) GetHotWords() ([]string, error) {
	// 获取今日热词（基于搜索历史）
	var keywords []string
	
	// 获取最近24小时的搜索历史，按出现频率排序
	var results []struct {
		Keyword string `json:"keyword"`
		Count   int    `json:"count"`
	}
	
	if err := s.db.Model(&models.SearchHistory{}).
		Select("keyword, COUNT(*) as count").
		Where("created_at >= DATE_SUB(NOW(), INTERVAL 24 HOUR)").
		Group("keyword").
		Order("count DESC, created_at DESC").
		Limit(10).
		Find(&results).Error; err != nil {
		// 如果查询失败，返回默认热词
		return []string{"校园", "活动", "学习", "美食", "运动", "兼职", "考试", "社团", "室友", "考研"}, nil
	}
	
	for _, result := range results {
		keywords = append(keywords, result.Keyword)
	}
	
	// 如果没有足够的热词，添加默认热词
	if len(keywords) < 10 {
		defaultKeywords := []string{"校园", "活动", "学习", "美食", "运动", "兼职", "考试", "社团", "室友", "考研"}
		for _, kw := range defaultKeywords {
			if !contains(keywords, kw) {
				keywords = append(keywords, kw)
				if len(keywords) >= 10 {
					break
				}
			}
		}
	}
	
	return keywords, nil
}

// GetSearchHistory 获取搜索历史
func (s *SearchService) GetSearchHistory(userID string) ([]string, error) {
	var histories []models.SearchHistory
	if err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(20).
		Find(&histories).Error; err != nil {
		return nil, fmt.Errorf("获取搜索历史失败: %v", err)
	}

	var keywords []string
	for _, history := range histories {
		keywords = append(keywords, history.Keyword)
	}
	return keywords, nil
}

// SaveSearchHistory 保存搜索历史
func (s *SearchService) SaveSearchHistory(userID, keyword string) error {
	// 先删除重复的搜索记录
	s.db.Where("user_id = ? AND keyword = ?", userID, keyword).Delete(&models.SearchHistory{})
	
	// 保存新的搜索记录
	history := models.SearchHistory{
		UserID:  userID,
		Keyword: keyword,
	}
	return s.db.Create(&history).Error
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

	query := s.db.Model(&models.Moment{}).Where("status = ?", 0).Preload("User")

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

	var posts []models.Moment
	var total int64

	// 统计总数
	query.Count(&total)

	// 分页查询
	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.PageSize).Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("获取筛选内容失败: %v", err)
	}

	return &SearchResponse{
		Posts: posts,
		Users: []models.User{},
		Pagination: Pagination{
			Page:     filter.Page,
			PageSize: filter.PageSize,
			Total:    total,
		},
	}, nil
}

// GetTrendingKeywords 获取趋势关键词
func (s *SearchService) GetTrendingKeywords(days int) ([]string, error) {
	// 基于搜索历史获取趋势关键词
	var keywords []string
	
	// 获取最近days天的搜索历史，按出现频率排序
	var results []struct {
		Keyword string `json:"keyword"`
		Count   int    `json:"count"`
	}
	
	if err := s.db.Model(&models.SearchHistory{}).
		Select("keyword, COUNT(*) as count").
		Where("created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)", days).
		Group("keyword").
		Order("count DESC, created_at DESC").
		Limit(10).
		Find(&results).Error; err != nil {
		// 如果查询失败，返回默认热词
		return []string{"校园", "活动", "学习", "美食", "运动"}, nil
	}
	
	for _, result := range results {
		keywords = append(keywords, result.Keyword)
	}
	
	// 如果没有足够的热词，添加默认热词
	if len(keywords) < 5 {
		defaultKeywords := []string{"期末考试", "寒假安排", "社团招新", "实习机会", "考研"}
		for _, kw := range defaultKeywords {
			if !contains(keywords, kw) {
				keywords = append(keywords, kw)
				if len(keywords) >= 10 {
					break
				}
			}
		}
	}
	
	return keywords, nil
}

// GetSearchSuggestions 获取搜索建议
func (s *SearchService) GetSearchSuggestions(keyword string) ([]string, error) {
	if keyword == "" {
		return s.GetHotWords()
	}

	var suggestions []string
	
	// 从用户名中获取建议
	if err := s.db.Model(&models.User{}).
		Where("username LIKE ?", "%"+keyword+"%").
		Limit(3).
		Pluck("username", &suggestions).Error; err != nil {
		return nil, fmt.Errorf("获取搜索建议失败: %v", err)
	}
	
	// 从搜索历史中获取建议
	var historySuggestions []string
	if err := s.db.Model(&models.SearchHistory{}).
		Where("keyword LIKE ?", "%"+keyword+"%").
		Order("created_at DESC").
		Limit(2).
		Pluck("keyword", &historySuggestions).Error; err != nil {
		return nil, fmt.Errorf("获取搜索建议失败: %v", err)
	}
	
	// 合并建议并去重
	for _, hs := range historySuggestions {
		if !contains(suggestions, hs) {
			suggestions = append(suggestions, hs)
		}
	}
	
	return suggestions, nil
}

// contains 检查字符串是否在切片中
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}