package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

func main() {
	// 连接数据库
	db, err := sql.Open("mysql", "workbench_user:ruangong7@tcp(106.52.165.122:3306)/campus_moments?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	page := 1
	pageSize := 3
	offset := (page - 1) * pageSize

	// 查询总数
	var total int
	err = db.QueryRow("SELECT COUNT(1) FROM posts WHERE status = 1").Scan(&total)
	if err != nil {
		panic(err)
	}
	fmt.Printf("总数: %d\n", total)

	// 查询数据
	rows, err := db.Query("SELECT id, title, content, images, user_id, like_count, comment_count, created_at FROM posts WHERE status = 1 ORDER BY created_at DESC LIMIT ? OFFSET ?", pageSize, offset)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var list []SimpleMoment
	for rows.Next() {
		var m SimpleMoment
		var imagesRaw sql.NullString
		var createdAt string

		if err := rows.Scan(&m.ID, &m.Title, &m.Content, &imagesRaw, &m.AuthorID, &m.LikeCount, &m.CommentCount, &createdAt); err != nil {
			panic(err)
		}

		// 处理 images 字段
		if imagesRaw.Valid && imagesRaw.String != "" && imagesRaw.String != "null" {
			if err := json.Unmarshal([]byte(imagesRaw.String), &m.Images); err != nil {
				m.Images = []string{}
			}
		} else {
			m.Images = []string{}
		}

		// 设置创建时间
		m.CreatedAt = time.Now()

		fmt.Printf("ID: %d, 标题: %s, 内容: %s, 图片: %v\n", m.ID, m.Title, m.Content, m.Images)
		list = append(list, m)
	}

	fmt.Printf("列表长度: %d\n", len(list))
	
	// 转换为JSON看看
	// 创建返回结果
	result := map[string]interface{}{
		"list": list,
		"pagination": map[string]int{
			"page":     page,
			"pageSize": pageSize,
			"total":    total,
		},
	}
	
	// 编码为JSON
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("返回的JSON:\n%s\n", string(jsonData))
}