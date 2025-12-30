package main

import (
	"fmt"
	"log"

	"github.com/Yw332/campus-moments-go/pkg/config"
	"github.com/Yw332/campus-moments-go/pkg/database"
)

func main() {
	// 初始化配置
	config.Init()

	// 初始化数据库连接
	database.Init()

	// 加载数据库连接
	db := database.GetDB()
	if db == nil {
		log.Fatal("数据库连接失败")
	}

	// 先删除所有外键约束（GORM 自动生成的约束名可能不同）
	_ = db.Exec("ALTER TABLE comments DROP FOREIGN KEY comments_ibfk_1")
	_ = db.Exec("ALTER TABLE comments DROP FOREIGN KEY comments_ibfk_2")
	_ = db.Exec("ALTER TABLE comments DROP FOREIGN KEY comments_ibfk_3")
	_ = db.Exec("ALTER TABLE comments DROP FOREIGN KEY comments_ibfk_4")
	_ = db.Exec("ALTER TABLE comments DROP FOREIGN KEY fk_comments_post")
	_ = db.Exec("ALTER TABLE comments DROP FOREIGN KEY comments_post_id_foreign")

	// 先修改 posts.id 为 bigint
	_ = db.Exec("ALTER TABLE posts MODIFY COLUMN id BIGINT NOT NULL AUTO_INCREMENT")
	// 修改主键（如果有重复定义）
	_ = db.Exec("ALTER TABLE posts DROP PRIMARY KEY")
	_ = db.Exec("ALTER TABLE posts ADD PRIMARY KEY (id)")
	fmt.Println("✅ 修改 posts.id 为 BIGINT")

	// 再修改 post_id 列类型为 bigint
	if err := db.Exec("ALTER TABLE comments MODIFY COLUMN post_id BIGINT NOT NULL").Error; err != nil {
		log.Printf("修改 comments.post_id 类型失败: %v", err)
	} else {
		fmt.Println("✅ 修改 comments.post_id 为 BIGINT")
	}

	// 重新添加外键约束
	if err := db.Exec("ALTER TABLE comments ADD CONSTRAINT fk_comments_post FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE").Error; err != nil {
		log.Printf("添加 comments 外键失败: %v", err)
	} else {
		fmt.Println("✅ 添加 comments 外键成功")
	}

	// 同样修复 likes 表的 target_id 列
	if err := db.Exec("ALTER TABLE likes MODIFY COLUMN target_id BIGINT NOT NULL").Error; err != nil {
		log.Printf("修改 likes.target_id 类型失败: %v", err)
	} else {
		fmt.Println("✅ 修改 likes.target_id 为 BIGINT")
	}

	fmt.Println("\n🎉 数据库迁移完成!")
}
