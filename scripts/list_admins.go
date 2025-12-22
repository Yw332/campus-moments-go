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
		log.Fatalln("数据库未连接")
	}

	rows, err := db.Raw("SELECT id, username, role FROM admins LIMIT 100").Rows()
	if err != nil {
		log.Fatalf("查询admins失败: %v", err)
	}
	defer rows.Close()
	fmt.Println("id\tusername\trole")
	for rows.Next() {
		var id int64
		var username, role string
		rows.Scan(&id, &username, &role)
		fmt.Printf("%d\t%s\t%s\n", id, username, role)
	}
}
