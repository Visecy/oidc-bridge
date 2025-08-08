package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"oidc-bridge/config"
	"oidc-bridge/handler"
	"oidc-bridge/model"
	"testing"

	"github.com/gin-gonic/gin"
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

func TestHandleDiscovery(t *testing.T) {
	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 调用处理函数
	handler.HandleDiscovery(c)

	// 验证响应
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// 解析响应内容
	var discovery model.Discovery
	if err := json.Unmarshal(w.Body.Bytes(), &discovery); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	// 验证关键字段
	if discovery.Issuer != config.AppConfig.Issuer {
		t.Errorf("Expected issuer %s, got %s", config.AppConfig.Issuer, discovery.Issuer)
	}

	if discovery.AuthorizationEndpoint != config.AppConfig.Issuer+"/authorize" {
		t.Errorf("Expected authorization endpoint %s, got %s", config.AppConfig.Issuer+"/authorize", discovery.AuthorizationEndpoint)
	}

	if discovery.TokenEndpoint != config.AppConfig.Issuer+"/token" {
		t.Errorf("Expected token endpoint %s, got %s", config.AppConfig.Issuer+"/token", discovery.TokenEndpoint)
	}

	if discovery.UserInfoEndpoint != config.AppConfig.Issuer+"/userinfo" {
		t.Errorf("Expected userinfo endpoint %s, got %s", config.AppConfig.Issuer+"/userinfo", discovery.UserInfoEndpoint)
	}

	if discovery.JwksURI != config.AppConfig.Issuer+"/.well-known/jwks.json" {
		t.Errorf("Expected jwks URI %s, got %s", config.AppConfig.Issuer+"/.well-known/jwks.json", discovery.JwksURI)
	}
}