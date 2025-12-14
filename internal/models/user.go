package models

import (
	"time"
)

type User struct {
	ID        int64     `json:"userId" gorm:"primaryKey;autoIncrement"`
	Username  string    `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Phone     string    `json:"phone" gorm:"uniqueIndex;size:20;not null"`
	Password  string    `json:"-" gorm:"size:255;not null"` // 密码不返回给前端
	Status    int       `json:"status" gorm:"default:1"`    // 1-正常 2-禁用 3-锁定
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// 表名
func (User) TableName() string {
	return "users"
}