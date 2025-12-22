package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/Yw332/campus-moments-go/internal/models"
	"github.com/Yw332/campus-moments-go/pkg/config"
	"github.com/Yw332/campus-moments-go/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	var username string
	var password string
	var cleanup bool

	flag.StringVar(&username, "username", "itestadmin", "admin username")
	flag.StringVar(&password, "password", "Admin12345", "admin password")
	flag.BoolVar(&cleanup, "cleanup", false, "delete the test admin and exit")
	flag.Parse()

	config.Init()
	database.Init()
	db := database.GetDB()
	if db == nil {
		log.Fatalln("数据库未连接")
	}

	if cleanup {
		if err := db.Where("username = ?", username).Delete(&models.Admin{}).Error; err != nil {
			log.Fatalf("删除测试管理员失败: %v", err)
		}
		fmt.Println("已删除测试管理员:", username)
		return
	}

	// insert or get
	pwd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	admin := models.Admin{Username: username, Password: string(pwd), CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := db.Where(models.Admin{Username: admin.Username}).FirstOrCreate(&admin).Error; err != nil {
		log.Fatalf("创建或获取测试管理员失败: %v", err)
	}
	fmt.Printf("测试管理员已准备：username=%s password=%s id=%d\n", admin.Username, password, admin.ID)
}
