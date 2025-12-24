package handlers

import (
	"database/sql"
	"encoding/json"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

// SimpleMoment 简单的动态结构
type SimpleMoment struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Images       []string  `json:"images"`
	AuthorID     string    `json:"authorId"`
	LikeCount    int64     `json:"likeCount"`
	CommentCount int64     `json:"commentCount"`
	CreatedAt    time.Time `json:"createdAt"`
}

// getMomentsFromDB 从数据库获取动态列表
func getMomentsFromDB(page, pageSize int) ([]SimpleMoment, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	// 直接连接数据库
	db, err := sql.Open("mysql", "workbench_user:ruangong7@tcp(106.52.165.122:3306)/campus_moments?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		return nil, 0, err
	}
	defer db.Close()

	// 查询总数
	var total int
	err = db.QueryRow("SELECT COUNT(1) FROM posts WHERE status = 1").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	rows, err := db.Query("SELECT id, title, content, images, user_id, like_count, comment_count, created_at FROM posts WHERE status = 1 ORDER BY created_at DESC LIMIT ? OFFSET ?", pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []SimpleMoment
	for rows.Next() {
		var m SimpleMoment
		var imagesRaw sql.NullString
		var createdAt string

		if err := rows.Scan(&m.ID, &m.Title, &m.Content, &imagesRaw, &m.AuthorID, &m.LikeCount, &m.CommentCount, &createdAt); err != nil {
			continue
		}

		if imagesRaw.Valid && imagesRaw.String != "" {
			if err := json.Unmarshal([]byte(imagesRaw.String), &m.Images); err != nil {
				m.Images = []string{}
			}
		} else {
			m.Images = []string{}
		}

		if createdAt != "" {
			// 尝试多种时间格式
			formats := []string{
				"2006-01-02 15:04:05.000",
				"2006-01-02 15:04:05",
			}
			for _, format := range formats {
				if parsedTime, err := time.Parse(format, createdAt); err == nil {
					m.CreatedAt = parsedTime
					break
				}
			}
		}

		list = append(list, m)
	}

	return list, total, nil
}