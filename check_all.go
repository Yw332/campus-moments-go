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

	// 查询所有status=1的记录
	fmt.Println("查询status=1的记录：")
	rows, err := db.Query("SELECT COUNT(1) FROM posts WHERE status = 1")
	if err != nil {
		panic(err)
	}
	
	var count int
	if rows.Next() {
		rows.Scan(&count)
	}
	fmt.Printf("status=1的总数: %d\n", count)
	
	// 查询前5条status=1的记录ID
	fmt.Println("\n前5条status=1的记录：")
	rows, err = db.Query("SELECT id, created_at FROM posts WHERE status = 1 ORDER BY created_at DESC LIMIT 5")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var createdAt string
		err := rows.Scan(&id, &createdAt)
		if err != nil {
			continue
		}
		fmt.Printf("ID: %d, 时间: %s\n", id, createdAt)
	}
}