package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
	"gorm.io/gorm"
)

// Moment 动态模型 - 与posts表对应
type Moment struct {
	ID              int             `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID          string          `json:"userId" gorm:"column:user_id;type:char(10);not null;index"`
	Title           string          `json:"title" gorm:"column:title;type:varchar(100)"`
	Content         string          `json:"content" gorm:"column:content;type:text;not null"`
	Images          json.RawMessage `json:"images" gorm:"column:images;type:json"`
	Video           string          `json:"video" gorm:"column:video;type:varchar(200)"`
	Visibility      int             `json:"visibility" gorm:"column:visibility;type:tinyint;default:0;comment:0-公开 1-好友 2-仅自己"`
	Status          int             `json:"status" gorm:"column:status;type:tinyint;default:0;comment:0-正常 1-删除"`
	Tags            json.RawMessage `json:"tags" gorm:"column:tags;type:json"`
	LikedUsers      json.RawMessage `json:"likedUsers" gorm:"column:liked_users;type:json"`
	CommentsSummary json.RawMessage `json:"commentsSummary" gorm:"column:comments_summary;type:json"`
	LikeCount       int             `json:"likeCount" gorm:"column:like_count;type:int;default:0"`
	CommentCount    int             `json:"commentCount" gorm:"column:comment_count;type:int;default:0"`
	ViewCount       int             `json:"viewCount" gorm:"column:view_count;type:int;default:0"`
	CreatedAt       time.Time       `json:"createdAt" gorm:"column:created_at;type:datetime"`
	UpdatedAt       time.Time       `json:"updatedAt" gorm:"column:updated_at;type:datetime"`
	
	// 关联字段
	User            *User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	
	// 兼容字段（保持向后兼容）
	AuthorID        string          `json:"authorId,omitempty" gorm:"-"`
	Author          *User           `json:"author,omitempty" gorm:"-"`
	Media           MediaItems      `json:"media,omitempty" gorm:"-"`
}

// Tags 标签数组类型
type Tags []string

// Value 实现 driver.Valuer 接口，用于数据库存储
func (t Tags) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return json.Marshal(t)
}

// Scan 实现 sql.Scanner 接口，用于数据库读取
func (t *Tags) Scan(value interface{}) error {
	if value == nil {
		*t = Tags{}
		return nil
	}
	
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, t)
	case string:
		return json.Unmarshal([]byte(v), t)
	}
	return nil
}

// MediaItem 媒体项
type MediaItem struct {
	URL      string `json:"url"`
	Type     string `json:"type"` // image/video
	Size     int64  `json:"size"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Duration int    `json:"duration"` // 视频时长(秒)，图片为0
}

// MediaItems 媒体项数组类型
type MediaItems []MediaItem

// Value 实现 driver.Valuer 接口
func (m MediaItems) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

// Scan 实现 sql.Scanner 接口
func (m *MediaItems) Scan(value interface{}) error {
	if value == nil {
		*m = MediaItems{}
		return nil
	}
	
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, m)
	case string:
		return json.Unmarshal([]byte(v), m)
	}
	return nil
}

// 表名
func (Moment) TableName() string {
	return "posts"
}

// BeforeFind GORM钩子 - 数据查询前处理
func (m *Moment) BeforeFind(tx *gorm.DB) (err error) {
	// 设置兼容字段
	m.AuthorID = m.UserID
	return nil
}

// AfterFind GORM钩子 - 数据查询后处理
func (m *Moment) AfterFind(tx *gorm.DB) (err error) {
	// 转换图片和视频到Media格式
	var media MediaItems
	
	// 处理图片
	if m.Images != nil && len(m.Images) > 0 {
		var images []string
		if err := json.Unmarshal(m.Images, &images); err == nil {
			for _, img := range images {
				media = append(media, MediaItem{
					URL:      img,
					Type:     "image",
					Size:     0,
					Width:    0,
					Height:   0,
					Duration: 0,
				})
			}
		}
	}
	
	// 处理视频
	if m.Video != "" {
		media = append(media, MediaItem{
			URL:      m.Video,
			Type:     "video",
			Size:     0,
			Width:    0,
			Height:   0,
			Duration: 0,
		})
	}
	
	m.Media = media
	
	// 设置兼容字段
	m.AuthorID = m.UserID
	m.Author = m.User
	
	return nil
}