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
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Images       []string  `json:"images"`
	AuthorID     string    `json:"authorId"`
	LikeCount    int64     `json:"likeCount"`
	CommentCount int64     `json:"commentCount"`
	CreatedAt    time.Time `json:"createdAt"`
}

// ListMoments 从数据库读取动态列表并返回总数
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
	if err := database.DB.QueryRow("SELECT COUNT(1) FROM posts WHERE status = 1").Scan(&total); err != nil {
		if err != sql.ErrNoRows {
			return nil, 0, err
		}
		total = 0
	}

	// 查询具体数据
	rows, err := database.DB.Query("SELECT id, title, content, images, user_id, like_count, comment_count, created_at FROM posts WHERE status = 1 ORDER BY created_at DESC LIMIT ? OFFSET ?", pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []Moment
	for rows.Next() {
		var m Moment
		var imagesRaw sql.NullString
		var createdAt string

		if err := rows.Scan(&m.ID, &m.Title, &m.Content, &imagesRaw, &m.AuthorID, &m.LikeCount, &m.CommentCount, &createdAt); err != nil {
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

		if createdAt != "" {
			if parsedTime, err := time.Parse("2006-01-02 15:04:05.000", createdAt); err == nil {
				m.CreatedAt = parsedTime
			}
		}

		list = append(list, m)
	}

	return list, total, nil
}