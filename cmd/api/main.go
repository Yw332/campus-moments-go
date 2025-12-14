package main

import (
	"log"
	"os"

	"github.com/Yw332/campus-moments-go/internal/models"
	"github.com/Yw332/campus-moments-go/internal/routes"
	"github.com/Yw332/campus-moments-go/pkg/config"
	"github.com/Yw332/campus-moments-go/pkg/database"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 加载环境变量 - 添加详细日志
	log.Println("=== Campus Moments Go 启动 ===")

	// 先检查文件是否存在
	if _, err := os.Stat(".env"); err == nil {
		log.Println("找到 .env 文件")
		if err := godotenv.Load(); err != nil {
			log.Printf("⚠️  加载 .env 文件失败: %v", err)
		} else {
			log.Println("✅ .env 文件加载成功")
		}
	} else {
		log.Println("⚠️  未找到 .env 文件")
	}

	// 2. 初始化配置
	config.Init()

	// 3. 初始化数据库（连接云服务器）
	database.Init()
	defer database.Close()

	// 4. 检查数据库连接状态
	if database.IsConnected() {
		log.Println("✅ 数据库连接正常")
		// 自动迁移数据库表结构
		models.AutoMigrate()
	} else {
		log.Println("⚠️  数据库未连接，某些功能可能不可用")
	}

	// 5. 设置Gin模式
	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
		log.Println("🚀 生产环境模式启动")
	} else {
		gin.SetMode(gin.DebugMode)
		log.Println("🔧 开发环境模式启动")
	}

	// 6. 创建Gin应用
	router := gin.Default()

	// 添加静态文件服务 - 提供上传文件的访问
	router.Static("/uploads", "./uploads")

	// 7. 注册路由（使用内部路由注册，保证返回格式一致）
	routes.SetupRoutes(router)

	// 8. 启动服务器
	port := config.Cfg.Server.Port
	log.Printf("✅ Campus Moments Go 启动成功")
	log.Printf("📡 服务器地址: http://106.52.165.122:%s", port)
	log.Printf("🌐 本地访问: http://localhost:%s", port)
	log.Printf("👤 GitHub: Yw332")
	log.Printf("🗄️  数据库: %s@%s:%s/%s",
		config.Cfg.Database.User,
		config.Cfg.Database.Host,
		config.Cfg.Database.Port,
		config.Cfg.Database.Name)

	// 监听所有网络接口以支持服务器访问
	if err := router.Run("0.0.0.0:" + port); err != nil {
		log.Fatal("❌ 服务器启动失败:", err)
	}
}
