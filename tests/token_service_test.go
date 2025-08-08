package tests

import (
	"oidc-bridge/config"
	"oidc-bridge/model"
	"oidc-bridge/service"
	"os"
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

func TestGenerateIDToken(t *testing.T) {
	// 确保密钥文件存在
	if _, err := os.Stat(config.AppConfig.PrivateKeyPath); os.IsNotExist(err) {
		t.Skip("Private key file not found, skipping ID token generation tests")
	}

	// 设置测试数据
	clientID := "test_client"
	redirectURI := "https://example.com/callback"
	userInfo := map[string]interface{}{
		"sub": "test_user",
		"name": "Test User",
		"email": "test@example.com",
	}

	// 存储 nonce 到 Redis
	err := service.SetNonce(clientID, redirectURI, "test_nonce")
	if err != nil {
		t.Fatalf("Failed to set nonce: %v", err)
	}

	// 调用函数
	token, err := service.GenerateIDToken(clientID, redirectURI, userInfo)
	if err != nil {
		t.Errorf("Failed to generate ID token: %v", err)
	}

	if token == "" {
		t.Error("Expected ID token, got empty string")
	}
}