package config

import (
"fmt"
"log"
"os"
"strconv"
"time"

"github.com/joho/godotenv"
)

// Config 全局配置
type Config struct {
App      AppConfig
Server   ServerConfig
Database DatabaseConfig
JWT      JWTConfig
}

type AppConfig struct {
Name    string
Version string
Env     string
}

type ServerConfig struct {
Port         string
Mode         string
ReadTimeout  time.Duration
WriteTimeout time.Duration
}

type DatabaseConfig struct {
Host         string
Port         string
Name         string
User         string
Password     string
MaxOpenConns int
MaxIdleConns int
DSN          string
}

type JWTConfig struct {
Secret      string
ExpireHours int
}

var Cfg *Config

// Init 初始化配置
func Init() {
// 加载.env文件 - 这里添加详细日志
log.Println("正在加载环境变量...")
if err := godotenv.Load(); err != nil {
log.Printf("⚠️  godotenv加载失败: %v", err)
log.Println("使用系统环境变量")
} else {
log.Println("✅ .env文件加载成功")
}

// 测试读取一个环境变量
testVal := os.Getenv("DB_HOST")
log.Printf("测试读取DB_HOST: %s", testVal)

Cfg = &Config{
App: AppConfig{
Name:    getEnv("APP_NAME", "Campus Moments Go API"),
Version: getEnv("APP_VERSION", "1.0.0"),
Env:     getEnv("APP_ENV", "development"),
},
Server: ServerConfig{
Port: getEnv("PORT", "3000"),
Mode: getEnv("GIN_MODE", "debug"),
},
Database: DatabaseConfig{
Host:         getEnv("DB_HOST", "localhost"),
Port:         getEnv("DB_PORT", "3306"),
Name:         getEnv("DB_NAME", "campus_moments"),
User:         getEnv("DB_USER", "root"),
Password:     getEnv("DB_PASSWORD", ""),
MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
},
JWT: JWTConfig{
Secret:      getEnv("JWT_SECRET", "your-default-secret-key"),
ExpireHours: getEnvAsInt("JWT_EXPIRE_HOURS", 24),
},
}

// 构建数据库连接字符串（云服务器）
Cfg.Database.DSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
Cfg.Database.User, Cfg.Database.Password, 
Cfg.Database.Host, Cfg.Database.Port, Cfg.Database.Name)

log.Printf("配置初始化完成:")
log.Printf("  环境: %s", Cfg.App.Env)
log.Printf("  数据库: %s@%s:%s", Cfg.Database.User, Cfg.Database.Host, Cfg.Database.Port)
}

func getEnv(key, defaultValue string) string {
value := os.Getenv(key)
if value == "" {
return defaultValue
}
return value
}

func getEnvAsInt(key string, defaultValue int) int {
if value, exists := os.LookupEnv(key); exists {
if intValue, err := strconv.Atoi(value); err == nil {
return intValue
}
}
return defaultValue
}

// IsProduction 是否为生产环境
func IsProduction() bool {
return Cfg.App.Env == "production"
}
