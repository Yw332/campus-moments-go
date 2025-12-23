package models

import (
	"time"
)

type User struct {
	ID             string    `json:"userId" gorm:"primaryKey;type:char(10)"`
	Username       string    `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Phone          string    `json:"phone" gorm:"uniqueIndex;size:20;not null"`
	Nickname       string    `json:"-" gorm:"-"`                    // 完全忽略此字段
	Avatar         string    `json:"avatar" gorm:"size:500"`
	AvatarType     int       `json:"avatarType" gorm:"type:TINYINT;default:0"` // 0-默认 1-自定义 2-系统生成
	AvatarUpdatedAt *time.Time `json:"avatarUpdatedAt" gorm:"type:DATETIME"`
	Signature      string    `json:"signature" gorm:"size:200"`
	Gender         int       `json:"-" gorm:"-"`                       // 完全忽略此字段
	Password       string    `json:"-" gorm:"size:255;not null"`       // 密码不返回给前端
	Status         int       `json:"status" gorm:"default:1"`         // 1-正常 2-禁用 3-锁定
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// 表名
func (User) TableName() string {
	return "users"
}
