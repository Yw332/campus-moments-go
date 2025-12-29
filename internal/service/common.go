package service

import (
	"github.com/Yw332/campus-moments-go/internal/models"
	"github.com/Yw332/campus-moments-go/pkg/database"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

// 全局数据库连接
var DB *gorm.DB

// getDB 获取数据库连接（懒加载）
func getDB() *gorm.DB {
	if DB == nil {
		DB = database.GetDB()
	}
	return DB
}

// isFriend 检查是否为好友
func IsFriend(userID1, userID2 string) bool {
	var count int64
	getDB().Model(&models.FriendRelation{}).
		Where("user_id = ? AND friend_id = ? AND relation_type = 1 AND status = 0", userID1, userID2).
		Count(&count)
	return count > 0
}

// getFriendIDs 获取好友ID列表
func GetFriendIDs(userID string) []string {
	var friendIDs []string
	getDB().Model(&models.FriendRelation{}).
		Where("user_id = ? AND relation_type = 1 AND status = 0", userID).
		Pluck("friend_id", &friendIDs)
	return friendIDs
}

// generateRandomColor 生成随机颜色
func GenerateRandomColor() string {
	colors := []string{"FF6B6B", "4ECDC4", "45B7D1", "96CEB4", "FFEAA7", "DDA0DD", "98D8C8", "F7DC6F"}
	rand.Seed(time.Now().UnixNano())
	return colors[rand.Intn(len(colors))]
}