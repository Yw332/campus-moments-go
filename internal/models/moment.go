package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Moment 动态模型 (对应 posts 表)
type Moment struct {
	ID               int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID           string     `json:"userId" gorm:"column:user_id;type:char(10);not null;index"`
	AuthorID         string     `json:"authorId" gorm:"column:author_id;type:char(10);index"`  // 保持兼容性
	Title            string     `json:"title" gorm:"column:title;type:varchar(100)"`
	Content          string     `json:"content" gorm:"column:content;type:text;not null"`
	Images           Tags       `json:"images" gorm:"column:images;type:json"`
	Video            string     `json:"video" gorm:"column:video;type:varchar(200)"`
	Visibility       int        `json:"visibility" gorm:"column:visibility;type:TINYINT;default:0"`  // 0公开/1好友/2私密
	Status           int        `json:"status" gorm:"column:status;type:TINYINT;default:1"`         // 1正常/2删除
	Tags             Tags       `json:"tags" gorm:"column:tags;type:json"`
	LikedUsers       Tags       `json:"likedUsers" gorm:"column:liked_users;type:json"`
	CommentsSummary  string     `json:"commentsSummary" gorm:"column:comments_summary;type:json"`
	LikeCount        int64      `json:"likeCount" gorm:"column:like_count;default:0"`
	CommentCount     int64      `json:"commentCount" gorm:"column:comment_count;default:0"`
	ViewCount        int64      `json:"viewCount" gorm:"column:view_count;default:0"`
	CreatedAt        time.Time  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt        time.Time  `json:"updatedAt" gorm:"column:updated_at"`

	// 兼容字段
	Media            Tags       `json:"media" gorm:"column:media;type:json"`  // 兼容旧的media字段

	// 关联字段
	Author *User `json:"author,omitempty" gorm:"foreignKey:AuthorID;references:ID;constraint:false"`
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
