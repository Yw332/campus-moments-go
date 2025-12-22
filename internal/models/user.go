package models

import (
	"time"
)

type User struct {
	ID        string    `json:"userId" gorm:"primaryKey;type:char(10)"`
	Username  string    `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Phone     string    `json:"phone" gorm:"uniqueIndex;size:20;not null"`
	Nickname  string    `json:"nickname" gorm:"size:50"`
	Avatar    string    `json:"avatar" gorm:"size:255"`
	Gender    int       `json:"gender" gorm:"default:0;comment:0未知/1男/2女"`
	Bio       string    `json:"bio" gorm:"size:500"`
	Password  string    `json:"-" gorm:"size:255;not null"` // 密码不返回给前端
	Status    int       `json:"status" gorm:"default:1"`    // 1-正常 2-禁用 3-锁定
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// 表名
func (User) TableName() string {
	return "users"
}
