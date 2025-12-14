package models

import (
	"log"

	"github.com/Yw332/campus-moments-go/pkg/database"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() {
	db := database.GetDB()
	
	// 仅添加avatar字段到现有用户表
	if err := db.AutoMigrate(&User{}); err != nil {
		log.Printf("用户表迁移失败: %v", err)
	} else {
		log.Println("✅ 用户表结构更新成功")
	}
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