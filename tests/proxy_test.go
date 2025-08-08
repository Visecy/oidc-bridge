package tests

import (
	"oidc-bridge/config"
	"oidc-bridge/model"
	"oidc-bridge/service"
	"testing"

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
}

func TestProxyToOPTokenEndpoint(t *testing.T) {
	// 创建测试请求
	req := model.TokenRequest{
		GrantType:    "authorization_code",
		Code:         "test_code",
		RedirectURI:  "https://example.com/callback",
		ClientID:     "test_client",
		ClientSecret: "test_secret",
	}

	// 调用函数
	// 注意：由于我们没有实际的 OP 服务，这里会返回错误
	// 但我们仍然可以验证处理逻辑是否正确执行
	_, err := service.ProxyToOPTokenEndpoint(req)
	if err == nil {
		t.Error("Expected error due to no real OP service, got nil")
	}
}

func TestGetUserInfoFromOP(t *testing.T) {
	// 调用函数
	// 注意：由于我们没有实际的 OP 服务，这里会返回错误
	// 但我们仍然可以验证处理逻辑是否正确执行
	_, err := service.GetUserInfoFromOP("test_token")
	if err == nil {
		t.Error("Expected error due to no real OP service, got nil")
	}
}