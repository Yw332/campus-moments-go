package token_blacklist

import (
	"sync"
	"time"
)

// TokenBlacklist Token黑名单管理
type TokenBlacklist struct {
	mu          sync.RWMutex
	blacklist   map[string]time.Time // token -> 过期时间
	cleanupTime time.Time             // 上次清理时间
}

var (
	instance *TokenBlacklist
	once     sync.Once
)

// GetInstance 获取单例实例
func GetInstance() *TokenBlacklist {
	once.Do(func() {
		instance = &TokenBlacklist{
			blacklist:   make(map[string]time.Time),
			cleanupTime: time.Now(),
		}
	})
	return instance
}

// AddToken 添加Token到黑名单
func (tb *TokenBlacklist) AddToken(token string, expiresAt time.Time) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	
	tb.blacklist[token] = expiresAt
	
	// 定期清理过期Token
	if time.Since(tb.cleanupTime) > 10*time.Minute {
		tb.cleanup()
		tb.cleanupTime = time.Now()
	}
}

// IsBlacklisted 检查Token是否在黑名单中
func (tb *TokenBlacklist) IsBlacklisted(token string) bool {
	tb.mu.RLock()
	defer tb.mu.RUnlock()
	
	expiresAt, exists := tb.blacklist[token]
	if !exists {
		return false
	}
	
	// 如果Token已过期，从黑名单中移除并返回false
	if time.Now().After(expiresAt) {
		delete(tb.blacklist, token)
		return false
	}
	
	return true
}

// cleanup 清理过期的Token
func (tb *TokenBlacklist) cleanup() {
	now := time.Now()
	for token, expiresAt := range tb.blacklist {
		if now.After(expiresAt) {
			delete(tb.blacklist, token)
		}
	}
}

// RemoveExpiredTokens 手动清理过期Token
func (tb *TokenBlacklist) RemoveExpiredTokens() {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.cleanup()
}

// GetBlacklistSize 获取黑名单大小（用于监控）
func (tb *TokenBlacklist) GetBlacklistSize() int {
	tb.mu.RLock()
	defer tb.mu.RUnlock()
	return len(tb.blacklist)
}