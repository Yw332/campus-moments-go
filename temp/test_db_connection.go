package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 数据库配置
	dsn := "workbench_user:ruangong7@tcp(106.52.165.122:3306)/campus_moments?charset=utf8mb4&parseTime=True&loc=Local"
	
	log.Println("=== 数据库连接测试 ===")
	log.Printf("DSN: %s", dsn)
	
	// 尝试连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("❌ 打开连接失败: %v", err)
		return
	}
	defer db.Close()
	
	// 设置连接超时
	db.SetConnMaxLifetime(10 * time.Second)
	
	// 测试连接（带超时）
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	log.Println("正在测试连接...")
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("❌ 连接测试失败: %v", err)
	} else {
		log.Println("✅ 连接测试成功")
	}
}