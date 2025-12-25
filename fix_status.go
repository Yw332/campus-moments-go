package main

import (
	"database/sql"
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

	// 更新所有status=1的记录为status=0（正常状态）
	// 将所有status=2的记录更新为status=1（删除状态）
	
	// 先更新状态为1的记录改为0
	result1, err := db.Exec("UPDATE posts SET status = 0 WHERE status = 1")
	if err != nil {
		log.Printf("更新正常状态失败: %v", err)
	} else {
		affected, _ := result1.RowsAffected()
		log.Printf("更新正常状态: %d 行受影响", affected)
	}

	// 查看更新后的数据
	rows, err := db.Query("SELECT id, status, content FROM posts ORDER BY created_at DESC LIMIT 5")
	if err != nil {
		log.Printf("查询失败: %v", err)
		return
	}
	defer rows.Close()

	log.Println("更新后的数据:")
	for rows.Next() {
		var id int
		var status int
		var content string
		err = rows.Scan(&id, &status, &content)
		if err != nil {
			log.Printf("扫描失败: %v", err)
			continue
		}
		log.Printf("ID: %d, 状态: %d, 内容: %s", id, status, content)
	}
}