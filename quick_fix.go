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

	// 查询有content但images为空的记录
	rows, err := db.Query("SELECT id, content FROM posts WHERE status = 0 AND (images IS NULL OR images = '[]' OR images = '')")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var content string
		err := rows.Scan(&id, &content)
		if err != nil {
			log.Printf("扫描失败: %v", err)
			continue
		}

		// 给这些记录添加一个示例图片
		imageUrls := []string{"http://106.52.165.122:8080/static/files/20251225013944_076749e3.png"}
		imageJSON, _ := json.Marshal(imageUrls)
		
		_, err = db.Exec("UPDATE posts SET images = ? WHERE id = ?", string(imageJSON), id)
		if err != nil {
			log.Printf("更新失败 ID %d: %v", id, err)
		} else {
			fmt.Printf("更新成功 ID %d: %s\n", id, string(imageJSON))
		}
	}
}