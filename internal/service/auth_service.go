package service

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Yw332/campus-moments-go/internal/models"
	"github.com/Yw332/campus-moments-go/pkg/database"
	"github.com/Yw332/campus-moments-go/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct{}

// NewAuthService 创建认证服务实例
func NewAuthService() *AuthService {
	return &AuthService{}
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Account  string `json:"account" binding:"required"`  // 用户名或手机号
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Token    string      `json:"token"`
	UserInfo interface{} `json:"userInfo"`
}

// validatePassword 验证密码强度
func validatePassword(password string) error {
	if len(password) < 8 || len(password) > 20 {
		return errors.New("密码长度必须在8-20位之间")
	}

	// 检查是否包含大小写字母和数字
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)

	if !hasUpper || !hasLower || !hasDigit {
		return errors.New("密码必须包含大小写字母和数字")
	}

	return nil
}

// validateUsername 验证用户名
func validateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("用户名长度必须在3-20个字符之间")
	}

	// 只允许字母、数字、中文和下划线
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_\p{Han}]+$`, username)
	if !matched {
		return errors.New("用户名只能包含字母、数字、中文和下划线")
	}

	return nil
}

// validatePhone 验证手机号
func validatePhone(phone string) error {
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
	if !matched {
		return errors.New("手机号格式不正确")
	}

	return nil
}

// Register 用户注册
func (s *AuthService) Register(req *RegisterRequest) (*models.User, error) {
	// 验证输入
	if err := validateUsername(req.Username); err != nil {
		return nil, err
	}

	if err := validatePhone(req.Phone); err != nil {
		return nil, err
	}

	if err := validatePassword(req.Password); err != nil {
		return nil, err
	}

	db := database.GetDB()

	// 检查用户名是否已存在
	var existingUser models.User
	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查手机号是否已存在
	if err := db.Where("phone = ?", req.Phone).First(&existingUser).Error; err == nil {
		return nil, errors.New("手机号已注册")
	}



	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}

	// 获取当前最大ID - 处理char类型的ID
	var maxIDStr string
	db.Raw("SELECT COALESCE(MAX(id), '0000000000') FROM users").Scan(&maxIDStr)
	log.Printf("当前最大ID字符串: %s", maxIDStr)
	
	// 将字符串ID转换为整数
	var maxID int64
	fmt.Sscanf(maxIDStr, "%d", &maxID)
	log.Printf("转换后的最大ID: %d", maxID)
	
	// 检查表结构，获取ID字段的默认值信息
	var columnInfo struct {
		Field      string `json:"Field"`
		Type       string `json:"Type"`
		Null       string `json:"Null"`
		Key        string `json:"Key"`
		Default    string `json:"Default"`
		Extra      string `json:"Extra"`
	}
	
	// 查询ID字段的详细信息
	if err := db.Raw("SHOW COLUMNS FROM users WHERE Field = 'id'").Scan(&columnInfo).Error; err == nil {
		log.Printf("ID字段信息: %+v", columnInfo)
	}

	// 创建用户 - 生成10位字符串ID，前面补零
	newID := maxID + 1
	idStr := fmt.Sprintf("%010d", newID) // 格式化为10位数字，前面补零
	log.Printf("生成的新ID: %s (数字: %d)", idStr, newID)
	
	sql := "INSERT INTO users (id, username, phone, password, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())"
	
	log.Printf("执行SQL: %s", sql)
	log.Printf("参数: id=%s, username=%s, phone=%s", idStr, req.Username, req.Phone)
	
	if err := db.Exec(sql, idStr, req.Username, req.Phone, string(hashedPassword), 1).Error; err != nil {
		return nil, fmt.Errorf("创建用户失败: %v", err)
	}

	user := &models.User{
		ID:       newID,
		Username: req.Username,
		Phone:    req.Phone,
		Status:   1,
	}

	return user, nil
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	db := database.GetDB()

	// 查找用户（支持用户名或手机号登录）
	var user models.User
	err := db.Where("username = ? OR phone = ?", req.Account, req.Account).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("账户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	// 检查用户状态
	if user.Status == 2 {
		return nil, errors.New("账户已被禁用")
	}
	if user.Status == 3 {
		return nil, errors.New("账户已被锁定，请联系管理员")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("密码错误")
	}

	// 生成JWT token
	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %v", err)
	}

	// 返回用户信息（不包含密码）
	userInfo := map[string]interface{}{
		"userId":   user.ID,
		"username": user.Username,
		"phone":    user.Phone,
	}

	return &LoginResponse{
		Token:    token,
		UserInfo: userInfo,
	}, nil
}

// GetUserByID 根据ID获取用户
func (s *AuthService) GetUserByID(userID int64) (*models.User, error) {
	db := database.GetDB()
	var user models.User

	if err := db.Where("id = ? AND status = 1", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	return &user, nil
}

// UpdatePassword 更新密码
func (s *AuthService) UpdatePassword(userID int64, oldPassword, newPassword string) error {
	// 验证新密码强度
	if err := validatePassword(newPassword); err != nil {
		return err
	}

	db := database.GetDB()
	var user models.User

	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("原密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %v", err)
	}

	// 更新密码
	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if err := db.Save(&user).Error; err != nil {
		return fmt.Errorf("更新密码失败: %v", err)
	}

	return nil
}