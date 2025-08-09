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

// SetupTests 初始化测试环境
func SetupTests() {
	// 加载测试配置
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
	if config.AppConfig.RedisAddr != "" {
		service.RedisClient = redis.NewClient(&redis.Options{
			Addr: config.AppConfig.RedisAddr,
		})
	}
}

// TestMain 是测试入口点
func TestMain(m *testing.M) {
	SetupTests()

	// 运行测试
	code := m.Run()
	os.Exit(code)
}
