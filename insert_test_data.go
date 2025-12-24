package main

import (
	"log"
	"time"

	"github.com/Yw332/campus-moments-go/internal/models"
	"github.com/Yw332/campus-moments-go/pkg/config"
	"github.com/Yw332/campus-moments-go/pkg/database"
)

func main() {
	log.Println("=== 插入测试数据 ===")
	
	// 初始化配置和数据库
	config.Init()
	database.Init()
	defer database.Close()

	if !database.IsConnected() {
		log.Fatal("❌ 数据库未连接")
	}

	// 自动迁移
	models.AutoMigrate()

	// 插入测试用户
	user := models.User{
		ID:        "1000000001",
		Username:  "testuser",
		Phone:     "13800138001",
		Avatar:    "/static/avatars/default.jpg",
		Password:  "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // 密码: password
		Status:    1, // 1-正常
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.GetDB().Create(&user).Error; err != nil {
		log.Printf("⚠️  用户可能已存在: %v", err)
	} else {
		log.Printf("✅ 创建测试用户成功: ID=%s", user.ID)
	}

	// 插入测试动态
	moments := []models.Moment{
		{
			UserID:      "1000000001",
			Title:       "今天天气真好",
			Content:     "今天天气真好，适合出去走走！☀️",
			Images:      []string{"static/images/weather1.jpg"},
			Visibility:  0, // 0公开
			Status:      1, // 1正常
			Tags:        []string{"天气", "日常"},
			LikedUsers:  []string{},
			LikeCount:   5,
			CommentCount: 2,
			ViewCount:   15,
			CreatedAt:   time.Now().Add(-2 * time.Hour),
			UpdatedAt:   time.Now().Add(-2 * time.Hour),
		},
		{
			UserID:      "1000000001",
			Title:       "学习心得分享",
			Content:     "分享一下最近的学习心得，坚持就是胜利！💪",
			Images:      []string{"static/images/study1.jpg"},
			Visibility:  0,
			Status:      1,
			Tags:        []string{"学习", "心得"},
			LikedUsers:  []string{},
			LikeCount:   8,
			CommentCount: 3,
			ViewCount:   25,
			CreatedAt:   time.Now().Add(-5 * time.Hour),
			UpdatedAt:   time.Now().Add(-5 * time.Hour),
		},
		{
			UserID:      "1000000001",
			Title:       "食堂美食推荐",
			Content:     "食堂今天的菜真不错，推荐大家试试！😋",
			Images:      []string{"static/images/food1.jpg", "static/images/food2.jpg"},
			Visibility:  0,
			Status:      1,
			Tags:        []string{"美食", "食堂"},
			LikedUsers:  []string{},
			LikeCount:   12,
			CommentCount: 4,
			ViewCount:   35,
			CreatedAt:   time.Now().Add(-1 * time.Hour),
			UpdatedAt:   time.Now().Add(-1 * time.Hour),
		},
		{
			UserID:      "1000000001",
			Title:       "期末考试复习",
			Content:     "期末考试要来了，大家都在复习吗？加油！📚",
			Images:      []string{"static/images/exam1.jpg"},
			Visibility:  0,
			Status:      1,
			Tags:        []string{"考试", "复习"},
			LikedUsers:  []string{},
			LikeCount:   15,
			CommentCount: 6,
			ViewCount:   50,
			CreatedAt:   time.Now().Add(-30 * time.Minute),
			UpdatedAt:   time.Now().Add(-30 * time.Minute),
		},
		{
			UserID:      "1000000001",
			Title:       "篮球赛邀请",
			Content:     "周末组织篮球赛，有人要一起吗？🏀",
			Video:       "static/videos/basketball.mp4",
			Visibility:  0,
			Status:      1,
			Tags:        []string{"运动", "篮球"},
			LikedUsers:  []string{},
			LikeCount:   3,
			CommentCount: 1,
			ViewCount:   20,
			CreatedAt:   time.Now().Add(-10 * time.Minute),
			UpdatedAt:   time.Now().Add(-10 * time.Minute),
		},
	}

	for i, moment := range moments {
		if err := database.GetDB().Create(&moment).Error; err != nil {
			log.Printf("⚠️  动态 %d 可能已存在: %v", i+1, err)
		} else {
			log.Printf("✅ 创建测试动态 %d 成功: ID=%d", i+1, moment.ID)
		}
	}

	log.Println("✅ 测试数据插入完成！")
	log.Printf("🌐 现在可以访问: http://localhost:8080/moments")
}