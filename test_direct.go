package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 连接数据库
	db, err := sql.Open("mysql", "workbench_user:ruangong7@tcp(106.52.165.122:3306)/campus_moments?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 测试查询
	query := "SELECT id, title, content, images, user_id, like_count, comment_count, created_at FROM posts WHERE status = 1 ORDER BY created_at DESC LIMIT 3 OFFSET 0"
	rows, err := db.Query(query)
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
			fmt.Printf("扫描错误: %v\n", err)
			continue
		}

		count++
		fmt.Printf("第%d条: ID=%d, Title=%s, Content=%s, Images=%v, UserID=%s\n", 
			count, id, title, content[:min(len(content), 20)]+"...", images.String, userID)
	}

	fmt.Printf("总共找到 %d 条记录\n", count)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}