package tests

import (
	"testing"
	"time"
	"oidc-bridge/service"
)

func TestMemoryCache(t *testing.T) {
	// 初始化内存缓存
	service.InitMemoryCache()

	// 测试设置和获取缓存项
	key := "test_key"
	value := "test_value"
	ttl := 10 * time.Second

	service.GlobalMemoryCache.Set(key, value, ttl)

	// 测试获取存在的缓存项
	retrievedValue, exists := service.GlobalMemoryCache.Get(key)
	if !exists {
		t.Error("Expected cache item to exist")
	}

	if retrievedValue != value {
		t.Errorf("Expected value %s, got %s", value, retrievedValue)
	}

	// 测试获取不存在的缓存项
	_, exists = service.GlobalMemoryCache.Get("nonexistent_key")
	if exists {
		t.Error("Expected cache item to not exist")
	}
}