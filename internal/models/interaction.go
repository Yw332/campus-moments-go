package models

import (
	"time"
)

// Comment 评论模型
type Comment struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	MomentID      int64     `json:"momentId" gorm:"index;not null"`             // 关联的动态ID
	UserID        string    `json:"userId" gorm:"type:char(10);index;not null"` // 评论者ID
	Content       string    `json:"content" gorm:"type:text;not null"`          // 评论内容
	ParentID      *int64    `json:"parentId" gorm:"index"`                      // 父评论ID（如果是回复评论）
	ReplyToUserID *string   `json:"replyToUserId" gorm:"type:char(10);index"`   // 被回复的用户ID
	LikeCount     int       `json:"likeCount" gorm:"default:0"`                 // 点赞数
	Status        int       `json:"status" gorm:"default:1;comment:1正常/2删除"`
	CreatedAt     time.Time `json:"createdAt"`

	// 关联
	User        *User      `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:false"`
	ReplyToUser *User      `json:"replyToUser,omitempty" gorm:"foreignKey:ReplyToUserID;constraint:false"`
	Replies     []*Comment `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
}

// TableName 表名
func (Comment) TableName() string {
	return "moment_comments"
}

// Like 点赞模型
type Like struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     string    `json:"userId" gorm:"type:char(10);index;not null"`       // 点赞者ID
	TargetID   int64     `json:"targetId" gorm:"index;not null"`                   // 目标ID（动态ID或评论ID）
	TargetType int       `json:"targetType" gorm:"index;not null;comment:1动态/2评论"` // 目标类型
	CreatedAt  time.Time `json:"createdAt"`

	// 关联
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:false"`
}

// TableName 表名
func (Like) TableName() string {
	return "likes"
}

// 常量定义
const (
	LikeTargetTypeMoment  = 1
	LikeTargetTypeComment = 2
)
