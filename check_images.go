package main

import (
	"database/sql"
	"encoding/json"
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

	// 查询有图片的动态
	rows, err := db.Query("SELECT id, content, images FROM posts WHERE status = 1 AND images IS NOT NULL AND images != '[]' AND images != '' ORDER BY created_at DESC")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	fmt.Println("有图片的动态：")
	for rows.Next() {
		var id int
		var content string
		var imagesRaw string
		
		err := rows.Scan(&id, &content, &imagesRaw)
		if err != nil {
			continue
		}
		
		var images []string
		if imagesRaw != "" && imagesRaw != "null" {
			json.Unmarshal([]byte(imagesRaw), &images)
		}
		
		fmt.Printf("ID: %d\n内容: %s\n图片: %v\n---\n", id, content, images)
	}
}