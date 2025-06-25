package plugin

import (
	"encoding/json"
	"fmt"
	"gitee.com/jamespi/drone_dispatch/config"
	"net/http"
	"sync"
	"time"
)

// 令牌管理器
type TokenManager struct {
	accessToken    string       // 存储访问令牌
	refreshToken   string       // 存储刷新令牌
	expiresAt      time.Time    // 令牌过期时间
	encryptionKey  []byte       // 加密密钥
	mutex          sync.RWMutex // 互斥锁，确保线程安全
	refreshEnabled bool         // 是否启用自动刷新
}

// NewTokenManager 创建一个新的令牌管理器
func NewTokenManager(encryptionKey string) *TokenManager {
	return &TokenManager{
		encryptionKey:  []byte(encryptionKey),
		refreshEnabled: true, // 默认启用自动刷新
	}

}

// 设置令牌
func (tm *TokenManager) SetAccessToken(token string, refreshToken string, expiresIn int) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	tm.accessToken = token
	tm.refreshToken = refreshToken
	tm.expiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
}

// 获取令牌
func (tm *TokenManager) GetAccessToken() string {
	tm.mutex.RLock()
	url := fmt.Sprintf("%s?appId=%s&appKey=%s", config.DjiSettings["url"], config.DjiSettings["appId"], config.DjiSettings["appKey"])

	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	// 解析响应获取token
	var result struct {
		AccessToken string `json:"accessToken"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	tm.mutex.RUnlock()
	return result.AccessToken
}

// 令牌刷新调度
func (tm *TokenManager) RefreshAccessToken() (string, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	// 检查令牌是否过期
	if time.Until(tm.expiresAt) < 5*time.Minute && tm.refreshEnabled {
		newToken := tm.GetAccessToken()
		if newToken == "" {
			return "", nil
		}
		tm.accessToken = newToken
		tm.refreshToken = newToken
		tm.expiresAt = time.Now().Add(time.Duration(config.TokenExpiresInSettings) * time.Second)
	}
	return tm.accessToken, nil

}
