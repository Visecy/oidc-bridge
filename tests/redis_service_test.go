package tests

import (
	"oidc-bridge/config"
	"oidc-bridge/model"
	"oidc-bridge/service"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

func init() {
	// 加载配置
	viper.SetConfigName("config_test")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	config.AppConfig = &model.Config{}
	if err := viper.Unmarshal(config.AppConfig); err != nil {
		panic(err)
	}

	// 初始化 Redis 客户端
	service.RedisClient = redis.NewClient(&redis.Options{
		Addr: config.AppConfig.RedisAddr,
	})
}

func TestSetAndGetNonce(t *testing.T) {
	// 设置测试数据
	clientID := "test_client"
	redirectURI := "https://example.com/callback"
	nonce := "test_nonce"

	// 调用 SetNonce
	err := service.SetNonce(clientID, redirectURI, nonce)
	if err != nil {
		t.Errorf("Failed to set nonce: %v", err)
	}

	// 调用 GetNonce
	retrievedNonce, err := service.GetNonce(clientID, redirectURI)
	if err != nil {
		t.Errorf("Failed to get nonce: %v", err)
	}

	// 验证结果
	if retrievedNonce != nonce {
		t.Errorf("Expected nonce %s, got %s", nonce, retrievedNonce)
	}
}

func TestGetNonceNotFound(t *testing.T) {
	// 使用不存在的键调用 GetNonce
	_, err := service.GetNonce("nonexistent_client", "https://example.com/callback")
	if err == nil {
		t.Error("Expected error for nonexistent nonce, got nil")
	}
}
