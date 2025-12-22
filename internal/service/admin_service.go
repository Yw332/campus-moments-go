package service

import (
	"errors"
	"time"

	"github.com/Yw332/campus-moments-go/internal/models"
	"github.com/Yw332/campus-moments-go/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

// AdminService 管理员相关业务
type AdminService struct{}

// NewAdminService 创建 AdminService
func NewAdminService() *AdminService {
	return &AdminService{}
}

// Authenticate 验证管理员账号密码
func (s *AdminService) Authenticate(username, password string) (*models.Admin, error) {
	db := database.GetDB()
	var admin models.Admin
	if err := db.Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, errors.New("管理员不存在或密码错误")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return nil, errors.New("管理员不存在或密码错误")
	}

	// 更新最后登录时间（异步，不阻塞主流程）
	go func() {
		admin.LastLoginAt = ptrTime(time.Now())
		db.Save(&admin)
	}()

	return &admin, nil
}

// IsAdminByUsername 判断 username 是否为管理员
func (s *AdminService) IsAdminByUsername(username string) (bool, error) {
	db := database.GetDB()
	var admin models.Admin
	if err := db.Where("username = ?", username).First(&admin).Error; err != nil {
		return false, nil
	}
	return true, nil
}

func ptrTime(t time.Time) *time.Time { return &t }
