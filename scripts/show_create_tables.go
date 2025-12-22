package main

import (
	"fmt"
	"log"

	"github.com/Yw332/campus-moments-go/pkg/config"
	"github.com/Yw332/campus-moments-go/pkg/database"
)

func main() {
	config.Init()
	database.Init()

	db := database.GetDB()
	if db == nil {
		log.Fatalln("数据库未连接，无法检测表结构")
	}

	tables := []string{"moments", "users", "posts", "posts"}
	for _, t := range tables {
		fmt.Printf("--- SHOW CREATE TABLE %s ---\n", t)
		row := db.Raw("SHOW CREATE TABLE " + t).Row()
		var tbl string
		var createSQL string
		err := row.Scan(&tbl, &createSQL)
		if err != nil {
			log.Printf("表 %s 获取失败: %v\n", t, err)
			continue
		}
		fmt.Println(createSQL)
		fmt.Println()
	}

	// Show column types
	fmt.Println("--- INFORMATION_SCHEMA.COLUMNS for moments and users ---")
	rows, err := db.Raw("SELECT TABLE_NAME, COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_KEY FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME IN ('moments','users')", config.Cfg.Database.Name).Rows()
	if err != nil {
		log.Printf("查询列信息失败: %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var table, col, ctype, isnull, ckey string
		err := rows.Scan(&table, &col, &ctype, &isnull, &ckey)
		if err != nil {
			log.Printf("scan err: %v", err)
			continue
		}
		fmt.Printf("%s.%s: %s, nullable=%s, key=%s\n", table, col, ctype, isnull, ckey)
	}
}
