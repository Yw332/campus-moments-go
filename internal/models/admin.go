package models

import "time"

// Admin 管理员模型（只读对应已有数据库表，不做迁移）
type Admin struct {
	ID          int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	Username    string     `json:"username" gorm:"size:20;not null;uniqueIndex"`
	Password    string     `json:"-" gorm:"size:80;not null"`
	Role        string     `json:"role" gorm:"size:15;default:admin"`
	LastLoginAt *time.Time `json:"lastLoginAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

func (Admin) TableName() string {
	return "admins"
}
