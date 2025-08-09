package service

import (
	"sync"
	"time"
)

// MemoryCache 本地内存缓存
type MemoryCache struct {
	data  map[string]*cacheItem
	mutex sync.RWMutex
}

// cacheItem 缓存项
type cacheItem struct {
	value      string
	expireTime time.Time
}

// NewMemoryCache 创建新的内存缓存实例
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		data: make(map[string]*cacheItem),
	}
}

// Set 设置缓存项
func (m *MemoryCache) Set(key, value string, ttl time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.data[key] = &cacheItem{
		value:      value,
		expireTime: time.Now().Add(ttl),
	}
}

// Get 获取缓存项
func (m *MemoryCache) Get(key string) (string, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	item, exists := m.data[key]
	if !exists {
		return "", false
	}

	// 检查是否过期
	if time.Now().After(item.expireTime) {
		// 删除过期项
		delete(m.data, key)
		return "", false
	}

	return item.value, true
}

// ClearExpired 清理过期项
func (m *MemoryCache) ClearExpired() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	for key, item := range m.data {
		if now.After(item.expireTime) {
			delete(m.data, key)
		}
	}
}

// GlobalMemoryCache 全局内存缓存实例
var GlobalMemoryCache *MemoryCache

// InitMemoryCache 初始化内存缓存
func InitMemoryCache() {
	GlobalMemoryCache = NewMemoryCache()

	// 启动定时清理过期项的goroutine
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			GlobalMemoryCache.ClearExpired()
		}
	}()
}
