package models

import (
	"time"
)

type User struct {
	ID              string    `json:"id" gorm:"primaryKey;column:id;type:char(10)"`
	Username        string    `json:"username" gorm:"column:username;type:varchar(20);not null"`
	Password        string    `json:"-" gorm:"column:password;type:varchar(80);not null"` // 密码不返回给前端
	Phone           string    `json:"phone" gorm:"column:phone;type:varchar(15)"`
	AvatarURL       string    `json:"avatarUrl" gorm:"column:avatar;type:varchar(500)"`
	AvatarType      int       `json:"avatarType" gorm:"column:avatar_type;type:tinyint"`
	AvatarUpdatedAt *time.Time `json:"avatarUpdatedAt" gorm:"column:avatar_updated_at;type:datetime"`
	PostCount       int       `json:"postCount" gorm:"column:post_count;type:int;default:0"`
	LikeCount       int       `json:"likeCount" gorm:"column:like_count;type:int;default:0"`
	CommentCount    int       `json:"commentCount" gorm:"column:comment_count;type:int;default:0"`
	Status          int64     `json:"status" gorm:"column:status;type:bigint"`
	Role            int       `json:"role" gorm:"column:role;type:tinyint;default:0;comment:0-普通用户 1-管理员"`
	LastLoginAt     *time.Time `json:"lastLoginAt" gorm:"column:last_login_at;type:datetime"`
	LastLoginIP     string    `json:"lastLoginIP" gorm:"column:last_login_ip;type:varchar(45)"`
	LoginCount      int       `json:"loginCount" gorm:"column:login_count;type:int;default:0"`
	CreatedAt       time.Time `json:"createdAt" gorm:"column:created_at;type:datetime(3)"`
	UpdatedAt       time.Time `json:"updatedAt" gorm:"column:updated_at;type:datetime(3)"`
	OpenID          string    `json:"openId" gorm:"column:openid;type:varchar(100)"`
	UnionID         string    `json:"unionId" gorm:"column:unionid;type:varchar(100)"`
	WechatNickname  string    `json:"wechatNickname" gorm:"column:wechat_nickname;type:varchar(100)"`
	WechatAvatar    string    `json:"wechatAvatar" gorm:"column:wechat_avatar;type:varchar(500)"`
	Signature       string    `json:"signature" gorm:"column:signature;type:varchar(200)"`
	LoginType       int       `json:"loginType" gorm:"column:login_type;type:tinyint"`
	LastActiveAt    *time.Time `json:"lastActiveAt" gorm:"column:last_active_at;type:datetime"`
}

// Admin 管理员表
type Admin struct {
	ID          int        `json:"id" gorm:"primaryKey;column:id;autoIncrement"`
	Username    string     `json:"username" gorm:"column:username;type:varchar(50);not null;uniqueIndex"`
	Password    string     `json:"-" gorm:"column:password;type:varchar(255);not null"`
	Role        string     `json:"role" gorm:"column:role;type:varchar(20);comment:super_admin,moderator,editor等"`
	LastLoginAt *time.Time `json:"lastLoginAt" gorm:"column:last_login_at;type:datetime"`
	CreatedAt   time.Time  `json:"createdAt" gorm:"column:created_at;type:datetime"`
	UpdatedAt   time.Time  `json:"updatedAt" gorm:"column:updated_at;type:datetime"`
}

// IsAdmin 判断用户是否是管理员
func (u *User) IsAdmin() bool {
	return u.Role == 1
}

// IsSuperAdmin 判断管理员是否是超级管理员
func (a *Admin) IsSuperAdmin() bool {
	return a.Role == "super_admin"
}

// 表名
func (User) TableName() string {
	return "users"
}

// 管理员表名
func (Admin) TableName() string {
	return "admins"
}