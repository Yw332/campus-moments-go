package service

import (
    "database/sql"
    "encoding/json"
    "errors"
    "time"

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
// 注意：此函数假设存在名为 `moments` 的表，且 `images` 字段为 JSON 字符串或 NULL。
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
    if err := database.DB.QueryRow("SELECT COUNT(1) FROM moments").Scan(&total); err != nil {
        if err != sql.ErrNoRows {
            return nil, 0, err
        }
        total = 0
    }

    // 查询具体数据
    rows, err := database.DB.Query("SELECT id, content, images, author_id, like_count, comment_count, created_at FROM moments ORDER BY created_at DESC LIMIT ? OFFSET ?", pageSize, offset)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var list []Moment
    for rows.Next() {
        var m Moment
        var imagesRaw sql.NullString
        var createdAt sql.NullTime

        if err := rows.Scan(&m.ID, &m.Content, &imagesRaw, &m.AuthorID, &m.LikeCount, &m.CommentCount, &createdAt); err != nil {
            return nil, 0, err
        }

        if imagesRaw.Valid && imagesRaw.String != "" {
            if err := json.Unmarshal([]byte(imagesRaw.String), &m.Images); err != nil {
                // 如果解析失败，忽略并返回空切片
                m.Images = []string{}
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
