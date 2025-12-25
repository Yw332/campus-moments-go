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

	page := 1
	pageSize := 5
	offset := (page - 1) * pageSize

	fmt.Printf("查询参数: page=%d, pageSize=%d, offset=%d\n", page, pageSize, offset)
	
	// 执行和getMomentsFromDB完全相同的查询
	fmt.Println("\n执行查询：")
	fmt.Printf("SELECT id, title, content, images, user_id, like_count, comment_count, created_at FROM posts WHERE status = 1 ORDER BY created_at DESC LIMIT %d OFFSET %d\n", pageSize, offset)
	
	rows, err := db.Query("SELECT id, title, content, images, user_id, like_count, comment_count, created_at FROM posts WHERE status = 1 ORDER BY created_at DESC LIMIT ? OFFSET ?", pageSize, offset)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	fmt.Println("查询结果：")
	count := 0
	for rows.Next() {
		count++
		var id int
		var title, content, userID string
		var likeCount, commentCount int64
		var createdAt string
		
		err := rows.Scan(&id, &title, &content, &userID, &likeCount, &commentCount, &createdAt)
		if err != nil {
			fmt.Printf("扫描失败: %v\n", err)
			continue
		}
		fmt.Printf("第%d条 - ID: %d, 标题: %s, 内容: %s, 用户: %s, 时间: %s\n", count, id, title, content, userID, createdAt)
	}
	fmt.Printf("总共扫描到: %d 行\n", count)
}