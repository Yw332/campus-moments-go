package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Yw332/campus-moments-go/internal/models"
	"github.com/Yw332/campus-moments-go/pkg/database"
	"gorm.io/gorm"
)

type VerificationService struct {
	db *gorm.DB
}

func NewVerificationService() *VerificationService {
	return &VerificationService{
		db: database.GetDB(),
	}
}

// SendVerificationCode 发送验证码
func (s *VerificationService) SendVerificationCode(phone, verificationType string) error {
	// 1. 验证手机号格式
	if len(phone) != 11 || phone[0:1] != "1" {
		return fmt.Errorf("手机号格式不正确")
	}

	// 2. 检查发送频率（1分钟内只能发送一次）
	var lastCode models.VerificationCode
	oneMinuteAgo := time.Now().Add(-1 * time.Minute)
	
	if err := s.db.Where("phone = ? AND created_at > ? AND type = ?", phone, oneMinuteAgo, verificationType).
		Order("created_at DESC").First(&lastCode).Error; err == nil {
		return fmt.Errorf("发送过于频繁，请1分钟后再试")
	}

	// 3. 生成6位验证码
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	
	// 4. 保存验证码
	verificationCode := models.VerificationCode{
		Phone:     phone,
		Code:      code,
		Type:      verificationType,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	if err := s.db.Create(&verificationCode).Error; err != nil {
		return fmt.Errorf("保存验证码失败: %v", err)
	}

	// 5. 发送短信（这里简化实现，实际需要调用短信服务API）
	err := s.sendSMS(phone, code)
	if err != nil {
		return fmt.Errorf("短信发送失败: %v", err)
	}

	return nil
}

// VerifyCode 验证验证码
func (s *VerificationService) VerifyCode(phone, code, verificationType string) error {
	var verificationCode models.VerificationCode
	
	// 查找最新的未使用验证码
	if err := s.db.Where("phone = ? AND code = ? AND type = ? AND is_used = ? AND expires_at > ?",
		phone, code, verificationType, false, time.Now()).
		Order("created_at DESC").First(&verificationCode).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("验证码无效或已过期")
		}
		return fmt.Errorf("验证失败: %v", err)
	}

	// 标记为已使用
	if err := s.db.Model(&verificationCode).Update("is_used", true).Error; err != nil {
		return fmt.Errorf("标记验证码失败: %v", err)
	}

	return nil
}

// ResetPasswordByPhone 通过手机号重置密码
func (s *VerificationService) ResetPasswordByPhone(phone, newPassword string) error {
	// 1. 查找用户
	var user models.User
	if err := s.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("该手机号未注册")
		}
		return fmt.Errorf("查询用户失败: %v", err)
	}

	// 2. 更新密码
	if err := s.db.Model(&user).Update("password", newPassword).Error; err != nil {
		return fmt.Errorf("重置密码失败: %v", err)
	}

	// 3. 记录重置日志
	resetLog := models.ResetPasswordLog{
		UserID:  fmt.Sprintf("%010d", user.ID),
		Phone:   phone,
		ResetAt: time.Now(),
	}
	
	s.db.Create(&resetLog)

	return nil
}

// sendSMS 发送短信（简化实现）
func (s *VerificationService) sendSMS(phone, code string) error {
	// 这里应该调用真实的短信服务API，如阿里云短信、腾讯云短信等
	// 这里只是记录日志
	fmt.Printf("验证码已发送至 %s: %s\n", phone, code)
	return nil
}

// CleanupExpiredCodes 清理过期验证码（应该作为定时任务执行）
func (s *VerificationService) CleanupExpiredCodes() {
	s.db.Where("expires_at < ?", time.Now()).Delete(&models.VerificationCode{})
}