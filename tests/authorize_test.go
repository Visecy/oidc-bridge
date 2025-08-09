package tests

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"oidc-bridge/config"
	"oidc-bridge/handler"
	"oidc-bridge/model"
	"oidc-bridge/service"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func init() {
	// 设置测试环境
	gin.SetMode(gin.TestMode)
}

// loadTestConfig 加载指定的测试配置文件
func loadTestConfig(configFile string) {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	config.AppConfig = &model.Config{}
	if err := viper.Unmarshal(config.AppConfig); err != nil {
		panic(err)
	}
}

// setupTestWithConfig 设置测试环境并加载指定配置
func setupTestWithConfig(configFile string) func() {
	originalConfig := config.AppConfig
	loadTestConfig(configFile)
	return func() { config.AppConfig = originalConfig }
}

func TestHandleAuthorize(t *testing.T) {
	// 使用默认测试配置
	defer setupTestWithConfig("config_test.yaml")()

	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 设置查询参数
	c.Request = &http.Request{
		Method: "GET",
		URL: &url.URL{
			RawQuery: "client_id=test_client&redirect_uri=https://example.com/callback&response_type=code&scope=openid&state=test_state&nonce=test_nonce",
		},
	}

	// 使用内存缓存替代 Redis 以避免连接问题
	service.InitMemoryCache()

	// 调用处理函数
	handler.HandleAuthorize(c)

	// 验证响应
	// 在 Gin 的测试模式中，重定向不会自动发生，所以我们需要检查状态码和 Location 头
	if w.Code != http.StatusFound {
		t.Errorf("Expected status code %d, got %d. Response body: %s", http.StatusFound, w.Code, w.Body.String())
	}

	// 验证重定向 URL
	location := w.Header().Get("Location")
	if location == "" {
		t.Error("Expected redirect location, got empty")
	}

	// 验证 nonce 是否存储在缓存中
	storedNonce, exists := service.GlobalMemoryCache.Get("nonce:test_client:https://example.com/callback")
	if !exists {
		t.Error("nonce should be stored in cache")
	}
	if storedNonce != "test_nonce" {
		t.Errorf("Expected nonce %s, got %s", "test_nonce", storedNonce)
	}
}

func TestHandleAuthorizeScopeMapping(t *testing.T) {
	// 使用 scope mapping 测试配置
	defer setupTestWithConfig("scope_mapping_test.yaml")()

	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 设置查询参数，包含需要映射的 scope
	c.Request = &http.Request{
		Method: "GET",
		URL: &url.URL{
			RawQuery: "client_id=test_client&redirect_uri=https://example.com/callback&response_type=code&scope=openid profile email custom_scope&state=test_state&nonce=test_nonce",
		},
	}

	// 使用内存缓存替代 Redis 以避免连接问题
	service.InitMemoryCache()

	// 调用处理函数
	handler.HandleAuthorize(c)

	// 验证响应
	// 在 Gin 的测试模式中，重定向不会自动发生，所以我们需要检查状态码和 Location 头
	if w.Code != http.StatusFound {
		t.Errorf("Expected status code %d, got %d. Response body: %s", http.StatusFound, w.Code, w.Body.String())
	}

	// 验证重定向 URL
	location := w.Header().Get("Location")
	if location == "" {
		t.Error("Expected redirect location, got empty")
	}

	// 解析重定向 URL 并验证 scope 映射
	parsedURL, err := url.Parse(location)
	if err != nil {
		t.Errorf("Failed to parse redirect URL: %v", err)
	}

	queryParams := parsedURL.Query()
	mappedScopes := queryParams.Get("scope")

	// 验证 openid 被正确移除（因为代码中会移除 openid）
	if strings.Contains(mappedScopes, "openid") {
		t.Error("openid scope should be removed from mapped scopes")
	}

	// 验证其他 scope 被正确映射
	expectedScopes := []string{"user_profile", "user_email", "mapped_custom_scope"}
	for _, expected := range expectedScopes {
		if !strings.Contains(mappedScopes, expected) {
			t.Errorf("Expected mapped scope %s not found in %s", expected, mappedScopes)
		}
	}
}

func TestHandleAuthorizeScopeMappingUnmapped(t *testing.T) {
	// 使用 scope mapping 测试配置
	defer setupTestWithConfig("scope_mapping_test.yaml")()

	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 设置查询参数，包含未映射的 scope
	c.Request = &http.Request{
		Method: "GET",
		URL: &url.URL{
			RawQuery: "client_id=test_client&redirect_uri=https://example.com/callback&response_type=code&scope=openid unmapped_scope&state=test_state&nonce=test_nonce",
		},
	}

	// 使用内存缓存替代 Redis 以避免连接问题
	service.InitMemoryCache()

	// 调用处理函数
	handler.HandleAuthorize(c)

	// 验证响应
	// 在 Gin 的测试模式中，重定向不会自动发生，所以我们需要检查状态码和 Location 头
	if w.Code != http.StatusFound {
		t.Errorf("Expected status code %d, got %d. Response body: %s", http.StatusFound, w.Code, w.Body.String())
	}

	// 验证重定向 URL
	location := w.Header().Get("Location")
	if location == "" {
		t.Error("Expected redirect location, got empty")
	}

	// 解析重定向 URL 并验证未映射的 scope 保持原样
	parsedURL, err := url.Parse(location)
	if err != nil {
		t.Errorf("Failed to parse redirect URL: %v", err)
	}

	queryParams := parsedURL.Query()
	mappedScopes := queryParams.Get("scope")

	// 验证未映射的 scope 保持原样
	if !strings.Contains(mappedScopes, "unmapped_scope") {
		t.Error("Unmapped scope should remain unchanged")
	}
}

func TestHandleAuthorizeInvalidResponseType(t *testing.T) {
	// 使用默认测试配置
	defer setupTestWithConfig("config_test.yaml")()

	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 设置查询参数
	c.Request = &http.Request{
		Method: "GET",
		URL: &url.URL{
			RawQuery: "client_id=test_client&redirect_uri=https://example.com/callback&response_type=invalid&scope=openid&state=test_state&nonce=test_nonce",
		},
	}

	// 使用内存缓存替代 Redis 以避免连接问题
	service.InitMemoryCache()

	// 调用处理函数
	handler.HandleAuthorize(c)

	// 验证响应
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandleAuthorizeMissingNonce(t *testing.T) {
	// 使用默认测试配置
	defer setupTestWithConfig("config_test.yaml")()

	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 设置查询参数
	c.Request = &http.Request{
		Method: "GET",
		URL: &url.URL{
			RawQuery: "client_id=test_client&redirect_uri=https://example.com/callback&response_type=code&scope=openid&state=test_state",
		},
	}

	// 使用内存缓存替代 Redis 以避免连接问题
	service.InitMemoryCache()

	// 调用处理函数
	handler.HandleAuthorize(c)

	// 验证响应
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}
