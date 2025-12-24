package main

import (
	"database/sql"
	"log"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "workbench_user:ruangong7@tcp(106.52.165.122:3306)/campus_moments?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	// 先确保用户存在
	userQuery := `INSERT IGNORE INTO users (id, username, phone, password, status, created_at, updated_at) 
	               VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(userQuery, "1000000001", "testuser", "13800138001", "password", 1, time.Now(), time.Now())
	if err != nil {
		log.Printf("用户创建失败: %v", err)
	} else {
		log.Printf("✅ 用户创建成功")
	}
	
	// 插入测试数据
	query := `INSERT INTO posts (user_id, author_id, title, content, images, visibility, status, tags, liked_users, comments_summary, like_count, comment_count, view_count, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	// 先查找一个存在的用户ID
	var userID string
	err = db.QueryRow("SELECT id FROM users LIMIT 1").Scan(&userID)
	if err != nil {
		log.Printf("没有找到用户: %v", err)
		return
	}
	
	result, err := db.Exec(query, 
		userID, userID,  // user_id, author_id
		"测试动态标题", 
		"这是一条测试动态内容，用于验证前端对接是否正常！🎉", 
		"[]", 
		0, 1, 
		"[\"测试\", \"前端\"]", 
		"[]", 
		"{}", 
		10, 3, 25, 
		time.Now(), time.Now())
	
	if err != nil {
		log.Printf("插入失败: %v", err)
	} else {
		id, _ := result.LastInsertId()
		log.Printf("✅ 插入测试数据成功: ID=%d", id)
	}
}