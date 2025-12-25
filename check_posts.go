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

	fmt.Println("检查posts表中的数据:")
	
	// 查询所有动态
	query := `
		SELECT id, content, images, status, created_at 
		FROM posts 
		ORDER BY created_at DESC 
		LIMIT 5
	`
	
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var content string
		var imagesRaw string
		var status int
		var createdAt string
		
		err = rows.Scan(&id, &content, &imagesRaw, &status, &createdAt)
		if err != nil {
			log.Printf("扫描行失败: %v", err)
			continue
		}

		fmt.Printf("ID: %d, 状态: %d, 时间: %s\n", id, status, createdAt)
		fmt.Printf("内容: %s\n", content)
		fmt.Printf("图片: %s\n", imagesRaw)
		fmt.Println("---")
	}
}