package auth

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

var (
	securityKey     string
	securityKeyOnce sync.Once
	keyMutex        sync.RWMutex
)

// GenerateSecurityKey 生成一个随机的安全密钥
func GenerateSecurityKey() string {
	securityKeyOnce.Do(func() {
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			// 如果随机生成失败，使用基于时间的密钥
			securityKey = hex.EncodeToString([]byte(time.Now().String()))[:32]
			return
		}
		securityKey = hex.EncodeToString(b)
	})
	return securityKey
}

// GetSecurityKey 获取当前的安全密钥
func GetSecurityKey() string {
	keyMutex.RLock()
	defer keyMutex.RUnlock()
	if securityKey == "" {
		return GenerateSecurityKey()
	}
	return securityKey
}

// ValidateSecurityKey 验证安全密钥是否正确
func ValidateSecurityKey(key string) bool {
	keyMutex.RLock()
	defer keyMutex.RUnlock()
	return key == securityKey
}

// ResetSecurityKey 重置安全密钥（生成新的）
func ResetSecurityKey() string {
	keyMutex.Lock()
	defer keyMutex.Unlock()
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		securityKey = hex.EncodeToString([]byte(time.Now().String()))[:32]
	} else {
		securityKey = hex.EncodeToString(b)
	}
	return securityKey
}
