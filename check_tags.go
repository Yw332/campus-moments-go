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

	// 查看动态表中tags字段的实际值
	rows, err := db.Query("SELECT id, tags FROM posts WHERE id IN (1, 4, 5) ORDER BY id")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	fmt.Println("=== posts表中tags字段的实际值 ===")
	for rows.Next() {
		var id int
		var tags sql.NullString
		err := rows.Scan(&id, &tags)
		if err != nil {
			panic(err)
		}
		
		if tags.Valid {
			fmt.Printf("ID: %d, Tags: %s (类型: %T)\n", id, tags.String, tags.String)
		} else {
			fmt.Printf("ID: %d, Tags: NULL\n", id)
		}
	}

	// 查看tags表结构
	fmt.Println("\n=== tags表数据 ===")
	tagRows, err := db.Query("SELECT id, name FROM tags LIMIT 5")
	if err != nil {
		panic(err)
	}
	defer tagRows.Close()

	for tagRows.Next() {
		var id int
		var name string
		err := tagRows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Tag ID: %d, Name: %s\n", id, name)
	}
}