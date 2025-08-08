package tests

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"oidc-bridge/handler"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	// 设置测试环境
	gin.SetMode(gin.TestMode)
}

func TestHandleAuthorize(t *testing.T) {
	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 设置查询参数
	c.Request = &http.Request{
		URL: &url.URL{
			RawQuery: "client_id=test_client&redirect_uri=https://example.com/callback&response_type=code&scope=openid&state=test_state&nonce=test_nonce",
		},
	}

	// 调用处理函数
	handler.HandleAuthorize(c)

	// 验证响应
	if w.Code != http.StatusFound {
		t.Errorf("Expected status code %d, got %d", http.StatusFound, w.Code)
	}

	// 验证重定向 URL
	location := w.Header().Get("Location")
	if location == "" {
		t.Error("Expected redirect location, got empty")
	}
}

func TestHandleAuthorizeInvalidResponseType(t *testing.T) {
	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 设置查询参数
	c.Request = &http.Request{
		URL: &url.URL{
			RawQuery: "client_id=test_client&redirect_uri=https://example.com/callback&response_type=invalid&scope=openid&state=test_state&nonce=test_nonce",
		},
	}

	// 调用处理函数
	handler.HandleAuthorize(c)

	// 验证响应
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandleAuthorizeMissingNonce(t *testing.T) {
	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 设置查询参数
	c.Request = &http.Request{
		URL: &url.URL{
			RawQuery: "client_id=test_client&redirect_uri=https://example.com/callback&response_type=code&scope=openid&state=test_state",
		},
	}

	// 调用处理函数
	handler.HandleAuthorize(c)

	// 验证响应
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}