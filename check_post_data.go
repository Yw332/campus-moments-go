package main

import (
	"database/sql"
	"encoding/json"
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

	// 查询最新的一条动态
	query := `
		SELECT id, content, images, tags, created_at 
		FROM posts 
		WHERE status = 0 
		ORDER BY created_at DESC 
		LIMIT 1
	`
	
	var id int
	var content string
	var imagesRaw string
	var tagsRaw string
	var createdAt string
	
	err = db.QueryRow(query).Scan(&id, &content, &imagesRaw, &tagsRaw, &createdAt)
	if err != nil {
		log.Fatal(err)
	}

	// 解析images字段
	var images []string
	if imagesRaw != "" && imagesRaw != "null" {
		json.Unmarshal([]byte(imagesRaw), &images)
	}

	// 解析tags字段
	var tags []string
	if tagsRaw != "" && tagsRaw != "null" {
		json.Unmarshal([]byte(tagsRaw), &tags)
	}

	fmt.Printf("最新动态信息:\n")
	fmt.Printf("ID: %d\n", id)
	fmt.Printf("内容: %s\n", content)
	fmt.Printf("图片: %v\n", images)
	fmt.Printf("标签: %v\n", tags)
	fmt.Printf("创建时间: %s\n", createdAt)
}