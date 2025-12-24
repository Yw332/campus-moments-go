package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "workbench_user:ruangong7@tcp(106.52.165.122:3306)/campus_moments")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	// 测试总数查询
	var total int
	err = db.QueryRow("SELECT COUNT(1) FROM posts WHERE status = 1").Scan(&total)
	if err != nil {
		log.Printf("查询总数失败: %v", err)
	} else {
		fmt.Printf("总数: %d\n", total)
	}
	
	// 测试数据查询
	rows, err := db.Query("SELECT id, title, content, images, user_id, like_count, comment_count, created_at FROM posts WHERE status = 1 ORDER BY created_at DESC LIMIT 5")
	if err != nil {
		log.Printf("查询数据失败: %v", err)
		return
	}
	defer rows.Close()
	
	fmt.Println("\n数据:")
	for rows.Next() {
		var id int64
		var title, content, userID string
		var images sql.NullString
		var likeCount, commentCount int64
		var createdAt string
		
		err := rows.Scan(&id, &title, &content, &images, &userID, &likeCount, &commentCount, &createdAt)
		if err != nil {
			log.Printf("扫描数据失败: %v", err)
			continue
		}
		
		fmt.Printf("ID: %d, Title: %s, UserID: %s, Likes: %d\n", id, title, userID, likeCount)
	}
}