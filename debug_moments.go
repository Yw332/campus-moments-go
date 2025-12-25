package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 连接数据库
	db, err := sql.Open("mysql", "workbench_user:ruangong7@tcp(106.52.165.122:3306)/campus_moments?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	fmt.Println("数据库连接成功")

	// 查询总数
	var total int
	err = db.QueryRow("SELECT COUNT(1) FROM posts WHERE status = 0").Scan(&total)
	if err != nil {
		log.Fatal("查询总数失败:", err)
	}
	fmt.Printf("正常状态的动态总数: %d\n", total)

	// 查询数据
	rows, err := db.Query("SELECT id, title, content, images, user_id, like_count, comment_count, created_at FROM posts WHERE status = 0 ORDER BY created_at DESC LIMIT 3 OFFSET 0")
	if err != nil {
		log.Fatal("查询数据失败:", err)
	}
	defer rows.Close()

	fmt.Println("查询到的动态:")
	for rows.Next() {
		var id int64
		var title, content, imagesRaw, userID string
		var likeCount, commentCount int64
		var createdAt string

		if err := rows.Scan(&id, &title, &content, &imagesRaw, &userID, &likeCount, &commentCount, &createdAt); err != nil {
			log.Printf("扫描行失败: %v", err)
			continue
		}

		fmt.Printf("ID: %d, 标题: %s, 内容: %s, 图片: %s, 用户: %s\n", id, title, content, imagesRaw, userID)
	}
}