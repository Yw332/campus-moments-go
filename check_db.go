package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "workbench_user:ruangong7@tcp(106.52.165.122:3306)/campus_moments")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	// 查看users表结构
	rows, err := db.Query("DESCRIBE users")
	if err != nil {
		log.Printf("查询users表失败: %v", err)
	} else {
		fmt.Println("=== users表结构 ===")
		for rows.Next() {
			var field, typ, null, key, extra string
			var def interface{}
			rows.Scan(&field, &typ, &null, &key, &def, &extra)
			fmt.Printf("%s: %s\n", field, typ)
		}
	}
	
	// 查看moments表结构（如果存在）
	rows2, err := db.Query("DESCRIBE moments")
	if err != nil {
		log.Printf("查询moments表失败: %v", err)
	} else {
		fmt.Println("\n=== moments表结构 ===")
		for rows2.Next() {
			var field, typ, null, key, extra string
			var def interface{}
			rows2.Scan(&field, &typ, &null, &key, &def, &extra)
			fmt.Printf("%s: %s\n", field, typ)
		}
	}
}