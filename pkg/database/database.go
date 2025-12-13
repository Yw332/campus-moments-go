package database

import (
"database/sql"
"log"

"github.com/Yw332/campus-moments-go/pkg/config"

_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// Init 初始化数据库连接
func Init() {
cfg := config.Cfg.Database

log.Printf("🔌 正在连接数据库: %s@%s:%s/%s", cfg.User, cfg.Host, cfg.Port, cfg.Name)

var err error
DB, err = sql.Open("mysql", cfg.DSN)
if err != nil {
log.Fatal("❌ 数据库连接失败:", err)
}

// 测试连接
err = DB.Ping()
if err != nil {
log.Fatal("❌ 数据库连接测试失败:", err)
}

log.Println("✅ 成功连接到云数据库")
}

// Close 关闭数据库连接
func Close() {
if DB != nil {
DB.Close()
}
}

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

