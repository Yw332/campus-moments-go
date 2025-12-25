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

	// 查询时间格式
	rows, err := db.Query("SELECT id, created_at FROM posts WHERE status = 1 ORDER BY created_at DESC LIMIT 3")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var createdAt string
		
		err := rows.Scan(&id, &createdAt)
		if err != nil {
			log.Printf("扫描失败: %v", err)
			continue
		}
		fmt.Printf("ID: %d, 时间字符串: '%s'\n", id, createdAt)
	}
}