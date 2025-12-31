package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Comment 评论模型
type Comment struct {
	ID            int           `json:"id" gorm:"primaryKey;autoIncrement"`
	PostID        int64         `json:"postId" gorm:"column:post_id;type:bigint;not null;index"`
	UserID        string        `json:"userId" gorm:"column:user_id;type:char(10);not null;index"`
	Content       string        `json:"content" gorm:"column:content;type:varchar(1000);not null"`
	Replies       json.RawMessage `json:"replies" gorm:"column:replies;type:json"`
	LikeCount     int           `json:"likeCount" gorm:"column:like_count;type:int;default:0"`
	IsAuthor      bool          `json:"isAuthor" gorm:"column:is_author;type:tinyint(1);default:0"`
	Status        int           `json:"status" gorm:"column:status;type:tinyint;default:0"`
	CreatedAt     time.Time     `json:"createdAt" gorm:"column:created_at;type:datetime"`
	UpdatedAt     time.Time     `json:"updatedAt" gorm:"column:updated_at;type:datetime"`

	// 关联字段（不设置外键约束）
	User          *User         `json:"user,omitempty" gorm:"-"`
	Post          *Post         `json:"post,omitempty" gorm:"-"`
}

// 表名
func (Comment) TableName() string {
	return "comments"
}

// Like 点赞模型
type Like struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     string    `json:"userId" gorm:"column:user_id;type:char(10);not null;index"`
	TargetType int       `json:"targetType" gorm:"column:target_type;type:tinyint;not null;comment:1-帖子 2-评论"`
	TargetID   int64     `json:"targetId" gorm:"column:target_id;type:bigint;not null;index"`
	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;type:datetime"`

	// 关联字段（不设置外键约束）
	User *User `json:"user,omitempty" gorm:"-"`
	Post *Post `json:"post,omitempty" gorm:"-"`
}

// 表名
func (Like) TableName() string {
	return "likes"
}

// Message 消息模型
type Message struct {
	ID             int64     `json:"id" gorm:"primaryKey;column:id;autoIncrement"`
	SenderID       string    `json:"senderId" gorm:"column:sender_id;type:char(10);not null;index"`
	ReceiverID     string    `json:"receiverId" gorm:"column:receiver_id;type:char(10);not null;index"`
	MsgType        int       `json:"msgType" gorm:"column:msg_type;type:tinyint;not null;comment:1-文本 2-图片 3-视频 4-文件"`
	ContentPreview string    `json:"contentPreview" gorm:"column:content_preview;type:varchar(255)"`
	FileURL        string    `json:"fileUrl" gorm:"column:file_url;type:varchar(500)"`
	FileSize       int       `json:"fileSize" gorm:"column:file_size;type:int"`
	IsEncrypted    bool      `json:"isEncrypted" gorm:"column:is_encrypted;type:tinyint(1);default:0"`
	IsRead         bool      `json:"isRead" gorm:"column:is_read;type:tinyint(1);default:0"`
	DeviceID       string    `json:"deviceId" gorm:"column:device_id;type:varchar(64)"`
	ServerMsgID    string    `json:"serverMsgId" gorm:"column:server_msg_id;type:varchar(64)"`
	CreatedAt      time.Time `json:"createdAt" gorm:"column:created_at;type:datetime"`

	// 关联字段（不设置外键约束）
	Sender         *User     `json:"sender,omitempty" gorm:"-"`
	Receiver       *User     `json:"receiver,omitempty" gorm:"-"`
}

// 表名
func (Message) TableName() string {
	return "messages"
}

// Conversation 会话模型
type Conversation struct {
	ID              int        `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID          string     `json:"userId" gorm:"column:user_id;type:char(10);not null;index"`
	PeerID          string     `json:"peerId" gorm:"column:peer_id;type:char(10);not null;index"`
	LastMsgID       int64      `json:"lastMsgId" gorm:"column:last_msg_id;type:bigint"`
	LastMsgPreview  string     `json:"lastMsgPreview" gorm:"column:last_msg_preview;type:varchar(255)"`
	UnreadCount     int        `json:"unreadCount" gorm:"column:unread_count;type:int;default:0"`
	IsPinned        bool       `json:"isPinned" gorm:"column:is_pinned;type:tinyint(1);default:0"`
	IsMuted         bool       `json:"isMuted" gorm:"column:is_muted;type:tinyint(1);default:0"`
	UpdatedAt       time.Time  `json:"updatedAt" gorm:"column:updated_at;type:datetime"`

	// 关联字段（不设置外键约束）
	User            *User      `json:"user,omitempty" gorm:"-"`
	Peer            *User      `json:"peer,omitempty" gorm:"-"`
	LastMessage     *Message   `json:"lastMessage,omitempty" gorm:"-"`
}

// 表名
func (Conversation) TableName() string {
	return "conversations"
}

// FriendRelation 好友关系模型
type FriendRelation struct {
	ID           int        `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID       string     `json:"userId" gorm:"column:user_id;type:char(10);not null;index"`
	FriendID     string     `json:"friendId" gorm:"column:friend_id;type:char(10);not null;index"`
	RelationType int        `json:"relationType" gorm:"column:relation_type;type:tinyint;comment:1-好友 2-黑名单"`
	RemarkName   string     `json:"remarkName" gorm:"column:remark_name;type:varchar(50)"`
	Status       int        `json:"status" gorm:"column:status;type:tinyint;default:0;comment:0-正常 1-已删除"`
	CreatedAt    time.Time  `json:"createdAt" gorm:"column:created_at;type:datetime"`
	UpdatedAt    time.Time  `json:"updatedAt" gorm:"column:updated_at;type:datetime"`

	// 关联字段（不设置外键约束）
	User         *User      `json:"user,omitempty" gorm:"-"`
	Friend       *User      `json:"friend,omitempty" gorm:"-"`
}

// 表名
func (FriendRelation) TableName() string {
	return "friend_relations"
}

// FriendRequest 好友请求模型
type FriendRequest struct {
	ID         int        `json:"id" gorm:"primaryKey;autoIncrement"`
	FromUserID string     `json:"fromUserId" gorm:"column:from_user_id;type:char(10);not null;index"`
	ToUserID   string     `json:"toUserId" gorm:"column:to_user_id;type:char(10);not null;index"`
	Message    string     `json:"message" gorm:"column:message;type:varchar(200)"`
	Status     int        `json:"status" gorm:"column:status;type:tinyint;default:0;comment:0-待处理 1-已同意 2-已拒绝"`
	ExpiresAt  time.Time  `json:"expiresAt" gorm:"column:expires_at;type:datetime"`
	CreatedAt  time.Time  `json:"createdAt" gorm:"column:created_at;type:datetime"`
	UpdatedAt  time.Time  `json:"updatedAt" gorm:"column:updated_at;type:datetime"`

	// 关联字段（不设置外键约束）
	FromUser   *User      `json:"fromUser,omitempty" gorm:"-"`
	ToUser     *User      `json:"toUser,omitempty" gorm:"-"`
}

// 表名
func (FriendRequest) TableName() string {
	return "friend_requests"
}

// Tag 标签模型
type Tag struct {
	ID          int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string     `json:"name" gorm:"column:name;type:varchar(20);not null;uniqueIndex"`
	Color       string     `json:"color" gorm:"column:color;type:varchar(7);comment:hex颜色值"`
	Icon        string     `json:"icon" gorm:"column:icon;type:varchar(50)"`
	Description string     `json:"description" gorm:"column:description;type:varchar(200)"`
	UsageCount  int        `json:"usageCount" gorm:"column:usage_count;type:int;default:0"`
	LastUsedAt  *time.Time `json:"lastUsedAt" gorm:"column:last_used_at;type:datetime"`
	Status      int        `json:"status" gorm:"column:status;type:tinyint;default:0;comment:0-正常 1-禁用"`
	CreatedAt   time.Time  `json:"createdAt" gorm:"column:created_at;type:datetime"`
}

// 表名
func (Tag) TableName() string {
	return "tags"
}



// Post 帖子模型（完整版本，与现有Moment合并）
type Post struct {
	ID              int64           `json:"id" gorm:"primaryKey;autoIncrement"`
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
	
	// 关联字段（不设置外键约束）
	User            *User           `json:"user,omitempty" gorm:"-"`
}

// 表名
func (Post) TableName() string {
	return "posts"
}

// ImageURLs 图片URL数组类型
type ImageURLs []string

// Value 实现 driver.Valuer 接口
func (i ImageURLs) Value() (driver.Value, error) {
	if i == nil {
		return nil, nil
	}
	return json.Marshal(i)
}

// Scan 实现 sql.Scanner 接口
func (i *ImageURLs) Scan(value interface{}) error {
	if value == nil {
		*i = ImageURLs{}
		return nil
	}
	
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, i)
	case string:
		return json.Unmarshal([]byte(v), i)
	}
	return nil
}

// TagArray 标签数组类型
type TagArray []string

// Value 实现 driver.Valuer 接口
func (t TagArray) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return json.Marshal(t)
}

// Scan 实现 sql.Scanner 接口
func (t *TagArray) Scan(value interface{}) error {
	if value == nil {
		*t = TagArray{}
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