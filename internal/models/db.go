package models

import (
	"log"

	"github.com/Yw332/campus-moments-go/pkg/database"
)

// AutoMigrate 自动迁移数据库表结构（已禁用）
func AutoMigrate() {
	log.Println("🚫 AutoMigrate 已禁用，使用现有数据库结构")
	return
	db := database.GetDB()
	if db == nil {
		log.Println("⚠️  数据库未连接，跳过迁移")
		return
	}

	// 迁移动态表
	if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci").AutoMigrate(&Moment{}); err != nil {
		log.Printf("❌ 动态表迁移失败: %v", err)
	} else {
		log.Println("✅ 动态表迁移成功")
	}

	// 完全跳过用户表的迁移，使用现有表结构
	log.Println("✅ 跳过用户表迁移，使用现有表结构")

	// 迁移验证码相关表
	if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci").AutoMigrate(&VerificationCode{}); err != nil {
		log.Printf("❌ 验证码表迁移失败: %v", err)
	} else {
		log.Println("✅ 验证码表迁移成功")
	}

	// 迁移搜索历史表
	if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci").AutoMigrate(&SearchHistory{}); err != nil {
		log.Printf("❌ 搜索历史表迁移失败: %v", err)
	} else {
		log.Println("✅ 搜索历史表迁移成功")
	}

	// 迁移重置密码日志表
	if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci").AutoMigrate(&ResetPasswordLog{}); err != nil {
		log.Printf("❌ 重置密码日志表迁移失败: %v", err)
	} else {
		log.Println("✅ 重置密码日志表迁移成功")
	}

	// 迁移评论表
	if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci").AutoMigrate(&Comment{}); err != nil {
		log.Printf("❌ 评论表迁移失败: %v", err)
	} else {
		log.Println("✅ 评论表迁移成功")
	}

	// 迁移点赞表
	if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci").AutoMigrate(&Like{}); err != nil {
		log.Printf("❌ 点赞表迁移失败: %v", err)
	} else {
		log.Println("✅ 点赞表迁移成功")
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

	// 检查并创建评论表
	if !db.Migrator().HasTable(&Comment{}) {
		if err := db.AutoMigrate(&Comment{}); err != nil {
			log.Printf("创建评论表失败: %v", err)
		} else {
			log.Println("✅ 评论表创建成功")
		}
	} else {
		log.Println("✅ 评论表已存在")
	}

	// 检查并创建点赞表
	if !db.Migrator().HasTable(&Like{}) {
		if err := db.AutoMigrate(&Like{}); err != nil {
			log.Printf("创建点赞表失败: %v", err)
		} else {
			log.Println("✅ 点赞表创建成功")
		}
	} else {
		log.Println("✅ 点赞表已存在")
	}
}
