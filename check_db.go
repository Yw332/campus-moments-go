package main

import (
	"fmt"
	"log"

	"github.com/Yw332/campus-moments-go/pkg/database"
)

func main() {
	// 连接数据库
	database.Init()
	defer database.Close()

	// 直接查询posts表
	db := database.GetDB()
	var results []struct {
		ID    int
		Title string
	}
	err := db.Raw("SELECT id, title FROM posts ORDER BY id DESC LIMIT 10").Scan(&results).Error
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 直接查询posts表 ===")
	for _, row := range results {
		fmt.Printf("ID: %d, Title: %s\n", row.ID, row.Title)
	}
}