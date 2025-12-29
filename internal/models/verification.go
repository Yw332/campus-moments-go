package models

import (
	"time"
)

// VerificationCode 验证码模型
type VerificationCode struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Phone       string    `json:"phone" gorm:"type:varchar(11);not null;index"`
	Code        string    `json:"code" gorm:"type:varchar(6);not null"`
	Type        string    `json:"type" gorm:"type:varchar(20);not null;default:'reset_password'"` // reset_password, login_verify
	IsUsed      bool      `json:"isUsed" gorm:"default:false"`
	ExpiresAt   time.Time `json:"expiresAt" gorm:"not null"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// SearchHistory 搜索历史模型
type SearchHistory struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    string    `json:"userId" gorm:"type:char(10);not null;index"`
	Keyword   string    `json:"keyword" gorm:"type:varchar(100);not null"`
	CreatedAt time.Time `json:"createdAt"`
}

// ResetPasswordLog 重置密码日志
type ResetPasswordLog struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    string    `json:"userId" gorm:"type:char(10);not null;index"`
	Phone     string    `json:"phone" gorm:"type:varchar(11);not null"`
	ResetAt   time.Time `json:"resetAt"`
	IP        string    `json:"ip" gorm:"type:varchar(45)"`
	UserAgent string    `json:"userAgent" gorm:"type:text"`
}

// 表名
func (VerificationCode) TableName() string {
	return "verification_codes"
}

func (SearchHistory) TableName() string {
	return "search_history"
}

func (ResetPasswordLog) TableName() string {
	return "reset_password_logs"
}