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
	// 加载可选 Issuer 的测试配置
	viper.SetConfigName("config_test_optional_issuer")
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

func TestHandleDiscoveryWithOptionalIssuer(t *testing.T) {
	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 创建一个测试请求，以便获取正确的 URL
	c.Request, _ = http.NewRequest("GET", "http://example.com/.well-known/openid-configuration", nil)

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
	expectedIssuer := "http://example.com"

	if discovery.Issuer != expectedIssuer {
		t.Errorf("Expected issuer %s, got %s", expectedIssuer, discovery.Issuer)
	}

	if discovery.AuthorizationEndpoint != expectedIssuer+"/authorize" {
		t.Errorf("Expected authorization endpoint %s, got %s", expectedIssuer+"/authorize", discovery.AuthorizationEndpoint)
	}

	if discovery.TokenEndpoint != expectedIssuer+"/token" {
		t.Errorf("Expected token endpoint %s, got %s", expectedIssuer+"/token", discovery.TokenEndpoint)
	}

	if discovery.UserInfoEndpoint != expectedIssuer+"/userinfo" {
		t.Errorf("Expected userinfo endpoint %s, got %s", expectedIssuer+"/userinfo", discovery.UserInfoEndpoint)
	}

	if discovery.JwksURI != expectedIssuer+"/.well-known/jwks.json" {
		t.Errorf("Expected jwks URI %s, got %s", expectedIssuer+"/.well-known/jwks.json", discovery.JwksURI)
	}
}