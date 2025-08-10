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

func TestGenerateIDTokenWithoutNonce(t *testing.T) {
	// 确保密钥文件存在
	if _, err := os.Stat(config.AppConfig.PrivateKeyPath); os.IsNotExist(err) {
		t.Skip("Private key file not found, skipping ID token generation tests")
	}

	// 设置测试数据
	clientID := "test_client_no_nonce"
	redirectURI := "https://example.com/callback"
	userInfo := map[string]interface{}{
		"sub":   "test_user",
		"name":  "Test User",
		"email": "test@example.com",
	}

	// 注意：我们不设置 nonce 到 Redis

	// 调用函数
	issuer := "http://localhost:8080"
	token, err := service.GenerateIDToken(issuer, clientID, redirectURI, userInfo)
	if err != nil {
		t.Errorf("Failed to generate ID token without nonce: %v", err)
	}

	if token == "" {
		t.Error("Expected ID token, got empty string")
	}
}