package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/Yw332/campus-moments-go/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DB    *sql.DB // 原生SQL连接
	GORM  *gorm.DB // GORM连接
)

// Init 初始化数据库连接
func Init() {
	cfg := config.Cfg.Database

	log.Printf("🔌 正在连接数据库: %s@%s:%s/%s", cfg.User, cfg.Host, cfg.Port, cfg.Name)

	// 初始化原生SQL连接
	var err error
	DB, err = sql.Open("mysql", cfg.DSN)
	if err != nil {
		log.Printf("⚠️  数据库连接失败: %v", err)
		log.Println("🚀 应用将在无数据库模式下启动")
		return
	}

	// 设置连接超时
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 测试原生连接（带超时）
	DB.SetConnMaxLifetime(10 * time.Second)
	err = DB.PingContext(ctx)
	if err != nil {
		log.Printf("⚠️  数据库连接测试失败: %v", err)
		log.Println("🚀 应用将在无数据库模式下启动")
		return
	}

	// 初始化GORM连接
	GORM, err = gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 生产环境静默日志
	})
	if err != nil {
		log.Printf("⚠️  GORM数据库连接失败: %v", err)
		log.Println("🚀 应用将在无数据库模式下启动")
		return
	}

	// 配置连接池
	sqlDB, err := GORM.DB()
	if err != nil {
		log.Printf("⚠️  获取GORM底层连接失败: %v", err)
		log.Println("🚀 应用将在无数据库模式下启动")
		return
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	log.Println("✅ 成功连接到数据库")
	log.Printf("📊 连接池配置: 最大连接数=%d, 空闲连接数=%d", cfg.MaxOpenConns, cfg.MaxIdleConns)
}

// Close 关闭数据库连接
func Close() {
	if DB != nil {
		DB.Close()
	}
	if GORM != nil {
		sqlDB, _ := GORM.DB()
		sqlDB.Close()
	}
}

// GetDB 获取GORM数据库连接
func GetDB() *gorm.DB {
	return GORM
}

// GetSQLDB 获取原生SQL数据库连接
func GetSQLDB() *sql.DB {
	return DB
}

// ============ 原生SQL方法（保持向后兼容）===========

// Query 执行查询
func Query(sql string, args ...interface{}) (*sql.Rows, error) {
	return DB.Query(sql, args...)
}

// QueryRow 执行单行查询
func QueryRow(sql string, args ...interface{}) *sql.Row {
	return DB.QueryRow(sql, args...)
}

// Exec 执行插入/更新/删除
func Exec(sql string, args ...interface{}) (sql.Result, error) {
	return DB.Exec(sql, args...)
}

// ============ GORM辅助方法 ============

// WithTransaction 在事务中执行操作
func WithTransaction(fn func(tx *gorm.DB) error) error {
	return GORM.Transaction(fn)
}

// BatchInsert 批量插入
func BatchInsert(data interface{}) error {
	return GORM.CreateInBatches(data, 100).Error
}

// IsConnected 检查数据库是否连接成功
func IsConnected() bool {
	if GORM == nil {
		return false
	}
	
	sqlDB, err := GORM.DB()
	if err != nil {
		return false
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	return sqlDB.PingContext(ctx) == nil
}