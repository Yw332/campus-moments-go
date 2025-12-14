package models

import (
	"log"

	"github.com/Yw332/campus-moments-go/pkg/database"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() {
	// 完全跳过用户表的迁移，使用现有表结构
	log.Println("✅ 跳过用户表迁移，使用现有表结构")
}

// CreateTables 如果表不存在则创建
func CreateTables() {
	db := database.GetDB()

	// 检查并创建用户表
	if !db.Migrator().HasTable(&User{}) {
		if err := db.AutoMigrate(&User{}); err != nil {
			log.Printf("创建用户表失败: %v", err)
		} else {
			log.Println("✅ 用户表创建成功")
		}
	} else {
		log.Println("✅ 用户表已存在")
	}
}