package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Yw332/campus-moments-go/internal/models"
	"github.com/Yw332/campus-moments-go/pkg/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct{}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{}
}

// UpdateProfileRequest 更新资料请求
type UpdateProfileRequest struct {
	Username   string `json:"username"`
	Phone      string `json:"phone"`
	Avatar     string `json:"avatar"`
	AvatarType int    `json:"avatarType"`
	Signature  string `json:"signature"`
}

// PublicUserInfo 公开用户信息（不含隐私字段）
type PublicUserInfo struct {
	ID              string     `json:"id"`
	Username        string     `json:"username"`
	Avatar          string     `json:"avatar"`
	AvatarType      int        `json:"avatarType"`
	AvatarUpdatedAt *time.Time `json:"avatarUpdatedAt"`
	PostCount       int        `json:"postCount"`
	LikeCount       int        `json:"likeCount"`
	CommentCount    int        `json:"commentCount"`
	Signature       string     `json:"signature"`
	LastActiveAt    *time.Time `json:"lastActiveAt"`
}

// GetUserByID 根据ID获取用户信息
func (s *UserService) GetUserByID(userID string) (*models.User, error) {
	db := database.GetDB()
	var user models.User

	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	return &user, nil
}

// GetPublicUserInfo 获取公开用户信息
func (s *UserService) GetPublicUserInfo(userID string) (*PublicUserInfo, error) {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return &PublicUserInfo{
		ID:              user.ID,
		Username:        user.Username,
		Avatar:          user.Avatar,
		AvatarType:      user.AvatarType,
		AvatarUpdatedAt: user.AvatarUpdatedAt,
		PostCount:       user.PostCount,
		LikeCount:       user.LikeCount,
		CommentCount:    user.CommentCount,
		Signature:       user.Signature,
		LastActiveAt:    user.LastActiveAt,
	}, nil
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(userID string, req *UpdateProfileRequest) (*models.User, error) {
	db := database.GetDB()
	var user models.User

	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
		now := time.Now()
		updates["avatar_updated_at"] = &now
	}
	if req.AvatarType != 0 {
		updates["avatar_type"] = req.AvatarType
	}
	if req.Signature != "" {
		updates["signature"] = req.Signature
	}

	if len(updates) == 0 {
		return &user, nil
	}

	if err := db.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 刷新数据
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdatePassword 更新密码
func (s *UserService) UpdatePassword(userID, oldPassword, newPassword string) error {
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
		return err
	}

	// 更新密码
	if err := db.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		return err
	}

	return nil
}

// UpdateAvatar 更新头像
func (s *UserService) UpdateAvatar(userID, avatarURL string, avatarType int) (*models.User, error) {
	db := database.GetDB()
	var user models.User

	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	now := time.Now()
	if err := db.Model(&user).Updates(map[string]interface{}{
		"avatar":           avatarURL,
		"avatar_type":      avatarType,
		"avatar_updated_at": &now,
	}).Error; err != nil {
		return nil, err
	}

	// 刷新数据
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateSignature 更新个性签名
func (s *UserService) UpdateSignature(userID, signature string) (*models.User, error) {
	db := database.GetDB()
	var user models.User

	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	if err := db.Model(&user).Update("signature", signature).Error; err != nil {
		return nil, err
	}

	// 刷新数据
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateLastActive 更新最后活跃时间
func (s *UserService) UpdateLastActive(userID string) error {
	db := database.GetDB()
	now := time.Now()

	if err := db.Model(&models.User{}).Where("id = ?", userID).Update("last_active_at", &now).Error; err != nil {
		return err
	}

	return nil
}

// SearchUsers 搜索用户
func (s *UserService) SearchUsers(keyword string, page, pageSize int) ([]PublicUserInfo, int64, error) {
	db := database.GetDB()

	var total int64
	query := db.Model(&models.User{}).
		Where("username LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []models.User
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// 转换为公开信息
	result := make([]PublicUserInfo, len(users))
	for i, user := range users {
		result[i] = PublicUserInfo{
			ID:              user.ID,
			Username:        user.Username,
			Avatar:          user.Avatar,
			AvatarType:      user.AvatarType,
			AvatarUpdatedAt: user.AvatarUpdatedAt,
			PostCount:       user.PostCount,
			LikeCount:       user.LikeCount,
			CommentCount:    user.CommentCount,
			Signature:       user.Signature,
			LastActiveAt:    user.LastActiveAt,
		}
	}

	return result, total, nil
}

// IncrementPostCount 增加帖子数
func (s *UserService) IncrementPostCount(userID string) error {
	db := database.GetDB()
	return db.Model(&models.User{}).Where("id = ?", userID).UpdateColumn("post_count", gorm.Expr("post_count + 1")).Error
}

// DecrementPostCount 减少帖子数
func (s *UserService) DecrementPostCount(userID string) error {
	db := database.GetDB()
	return db.Model(&models.User{}).Where("id = ?", userID).Where("post_count > 0").UpdateColumn("post_count", gorm.Expr("post_count - 1")).Error
}

// IncrementLikeCount 增加点赞数
func (s *UserService) IncrementLikeCount(userID string) error {
	db := database.GetDB()
	return db.Model(&models.User{}).Where("id = ?", userID).UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
}

// DecrementLikeCount 减少点赞数
func (s *UserService) DecrementLikeCount(userID string) error {
	db := database.GetDB()
	return db.Model(&models.User{}).Where("id = ?", userID).Where("like_count > 0").UpdateColumn("like_count", gorm.Expr("like_count - 1")).Error
}

// IncrementCommentCount 增加评论数
func (s *UserService) IncrementCommentCount(userID string) error {
	db := database.GetDB()
	return db.Model(&models.User{}).Where("id = ?", userID).UpdateColumn("comment_count", gorm.Expr("comment_count + 1")).Error
}

// DecrementCommentCount 减少评论数
func (s *UserService) DecrementCommentCount(userID string) error {
	db := database.GetDB()
	return db.Model(&models.User{}).Where("id = ?", userID).Where("comment_count > 0").UpdateColumn("comment_count", gorm.Expr("comment_count - 1")).Error
}

// ResetUserPassword 管理员重置用户密码
func (s *UserService) ResetUserPassword(targetUserID, newPassword string) error {
	db := database.GetDB()
	var user models.User

	if err := db.Where("id = ?", targetUserID).First(&user).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新密码
	if err := db.Model(&user).Updates(map[string]interface{}{
		"password":   string(hashedPassword),
		"updated_at": time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	return nil
}

// AdminGetAllUsers 管理员获取所有用户列表
func (s *UserService) AdminGetAllUsers(page, pageSize int, keyword string) ([]models.User, int64, error) {
	db := database.GetDB()

	var total int64
	query := db.Model(&models.User{})

	// 关键词搜索
	if keyword != "" {
		query = query.Where("username LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []models.User
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// AdminBanUser 管理员封禁用户（status: 0-正常, 1-封禁）
func (s *UserService) AdminBanUser(targetUserID string) error {
	db := database.GetDB()

	var user models.User
	if err := db.Where("id = ?", targetUserID).First(&user).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 设置status为1表示封禁
	if err := db.Model(&user).Updates(map[string]interface{}{
		"status":     int64(1),
		"updated_at": time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("封禁用户失败: %w", err)
	}

	return nil
}

// AdminUnbanUser 管理员解封用户
func (s *UserService) AdminUnbanUser(targetUserID string) error {
	db := database.GetDB()

	var user models.User
	if err := db.Where("id = ?", targetUserID).First(&user).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 设置status为0表示正常
	if err := db.Model(&user).Updates(map[string]interface{}{
		"status":     int64(0),
		"updated_at": time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("解封用户失败: %w", err)
	}

	return nil
}

// AdminDeleteUser 管理员删除用户
func (s *UserService) AdminDeleteUser(targetUserID string) error {
	db := database.GetDB()

	var user models.User
	if err := db.Where("id = ?", targetUserID).First(&user).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 删除用户（级联删除相关数据）
	if err := db.Delete(&user).Error; err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}

	return nil
}
