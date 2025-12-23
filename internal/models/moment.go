package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Moment 动态模型
type Moment struct {
	ID           int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	Content      string          `json:"content" gorm:"type:text;not null"`
	AuthorID     string          `json:"authorId" gorm:"type:char(10);not null;index"`
	Tags         Tags            `json:"tags" gorm:"type:json"`
	Media        MediaItems      `json:"media" gorm:"type:json"`
	Visibility   int             `json:"visibility" gorm:"default:0;comment:0公开/1好友/2私密"`
	LikeCount    int             `json:"likeCount" gorm:"default:0"`
	CommentCount int             `json:"commentCount" gorm:"default:0"`
	Status       int             `json:"status" gorm:"default:1;comment:1正常/2删除"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
	
	// 关联字段
	Author       *User           `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
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