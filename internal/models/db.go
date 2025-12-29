package models

import (
	"log"

	"github.com/Yw332/campus-moments-go/pkg/database"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() {
	db := database.GetDB()
	if db == nil {
		log.Println("⚠️  数据库未连接，跳过迁移")
		return
	}

	// 迁移所有核心表
	tables := []interface{}{
		&User{},           // users表
		&Moment{},         // posts表
		&Comment{},        // comments表
		&Like{},           // likes表
		&Message{},        // messages表
		&Conversation{},   // conversations表
		&FriendRelation{}, // friend_relations表
		&FriendRequest{},  // friend_requests表
		&Tag{},            // tags表
		&SearchHistory{},  // search_history表
		&VerificationCode{}, // verification_codes表
		&ResetPasswordLog{}, // reset_password_logs表
	}

	for _, table := range tables {
		if err := db.AutoMigrate(table); err != nil {
			log.Printf("❌ %T 表迁移失败: %v", table, err)
		} else {
			log.Printf("✅ %T 表迁移成功", table)
		}
	}
}

// CreateTables 如果表不存在则创建
func CreateTables() {
	db := database.GetDB()
	if db == nil {
		log.Println("⚠️  数据库未连接，无法创建表")
		return
	}

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

	// 检查并创建动态表
	if !db.Migrator().HasTable(&Moment{}) {
		if err := db.AutoMigrate(&Moment{}); err != nil {
			log.Printf("创建动态表失败: %v", err)
		} else {
			log.Println("✅ 动态表创建成功")
		}
	} else {
		log.Println("✅ 动态表已存在")
	}
}