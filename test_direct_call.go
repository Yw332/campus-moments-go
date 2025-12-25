package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
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
	CreatedAt    string    `json:"createdAt"`
}

func main() {
	// 直接测试HTTP接口
	resp, err := http.Get("http://localhost:8080/moments?page=1&pageSize=3")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result struct {
		Code    int `json:"code"`
		Message string `json:"message"`
		Data    struct {
			List       []SimpleMoment `json:"list"`
			Pagination struct {
				Page     int `json:"page"`
				PageSize int `json:"pageSize"`
				Total    int `json:"total"`
			} `json:"pagination"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		panic(err)
	}

	fmt.Printf("=== HTTP响应 ===\n")
	fmt.Printf("Code: %d\n", result.Code)
	fmt.Printf("Message: %s\n", result.Message)
	fmt.Printf("List is nil: %v\n", result.Data.List == nil)
	fmt.Printf("List length: %d\n", len(result.Data.List))
	fmt.Printf("Total: %d\n", result.Data.Pagination.Total)

	// 直接调用数据库函数
	db, err := sql.Open("mysql", "workbench_user:ruangong7@tcp(106.52.165.122:3306)/campus_moments?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 查询数据
	rows, err := db.Query("SELECT id, title, content, images, user_id, like_count, comment_count, created_at FROM posts WHERE status = 1 ORDER BY created_at DESC LIMIT 3 OFFSET 0")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int64
		var title, content, userID string
		var images sql.NullString
		var likeCount, commentCount int64
		var createdAt string

		err := rows.Scan(&id, &title, &content, &images, &userID, &likeCount, &commentCount, &createdAt)
		if err != nil {
			fmt.Printf("Scan错误: %v\n", err)
			continue
		}
		count++
		fmt.Printf("第%d条: ID=%d, Title=%s, Images=%v\n", count, id, title, images.String)
	}
}