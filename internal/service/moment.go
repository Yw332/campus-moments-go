package service

import (
    "database/sql"
    "encoding/json"
    "errors"
    "time"

    "github.com/Yw332/campus-moments-go/internal/models"
    "github.com/Yw332/campus-moments-go/pkg/database"
)

// Moment 简单的业务模型，供 handler 返回给前端
type Moment struct {
    ID           int64     `json:"id"`
    Content      string    `json:"content"`
    Images       []string  `json:"images"`
    AuthorID     int64     `json:"authorId"`
    LikeCount    int       `json:"likeCount"`
    CommentCount int       `json:"commentCount"`
    CreatedAt    time.Time `json:"createdAt"`
}

// ListMoments 从数据库读取动态列表并返回总数
// 注意：此函数假设存在名为 `posts` 的表，且 `media` 字段为 JSON 字符串或 NULL。
func ListMoments(page, pageSize int) ([]Moment, int, error) {
    if database.DB == nil {
        return nil, 0, errors.New("database not initialized")
    }

    if page < 1 {
        page = 1
    }
    if pageSize <= 0 {
        pageSize = 10
    }
    offset := (page - 1) * pageSize

    // 查询总数
    var total int
    if err := database.DB.QueryRow("SELECT COUNT(1) FROM posts").Scan(&total); err != nil {
        if err != sql.ErrNoRows {
            return nil, 0, err
        }
        total = 0
    }

    // 查询具体数据
    rows, err := database.DB.Query("SELECT id, content, media, author_id, like_count, comment_count, created_at FROM posts ORDER BY created_at DESC LIMIT ? OFFSET ?", pageSize, offset)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var list []Moment
    for rows.Next() {
        var m Moment
        var mediaRaw sql.NullString
        var createdAt sql.NullTime

        if err := rows.Scan(&m.ID, &m.Content, &mediaRaw, &m.AuthorID, &m.LikeCount, &m.CommentCount, &createdAt); err != nil {
            return nil, 0, err
        }

        // 从媒体JSON中提取图片URL数组
        if mediaRaw.Valid && mediaRaw.String != "" {
            var mediaItems []models.MediaItem
            if err := json.Unmarshal([]byte(mediaRaw.String), &mediaItems); err != nil {
                // 如果解析失败，忽略并返回空切片
                m.Images = []string{}
            } else {
                // 提取所有图片的URL
                for _, item := range mediaItems {
                    if item.Type == "image" {
                        m.Images = append(m.Images, item.URL)
                    }
                }
            }
        } else {
            m.Images = []string{}
        }

        if createdAt.Valid {
            m.CreatedAt = createdAt.Time
        }

        list = append(list, m)
    }

    return list, total, nil
}
