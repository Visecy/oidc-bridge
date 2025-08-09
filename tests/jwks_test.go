package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"oidc-bridge/config"
	"oidc-bridge/handler"
	"oidc-bridge/model"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	// 确保密钥文件存在
	if _, err := os.Stat(config.AppConfig.PrivateKeyPath); os.IsNotExist(err) {
		os.Exit(0)
	}

	if _, err := os.Stat(config.AppConfig.PublicKeyPath); os.IsNotExist(err) {
		os.Exit(0)
	}
}

func TestHandleJWKS(t *testing.T) {
	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 调用处理函数
	handler.HandleJWKS(c)

	// 验证响应
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// 解析响应内容
	var jwks model.JWKS
	if err := json.Unmarshal(w.Body.Bytes(), &jwks); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	// 验证关键字段
	if len(jwks.Keys) == 0 {
		t.Error("Expected at least one key in JWKS")
	}

	key := jwks.Keys[0]
	if key.KTY != "RSA" {
		t.Errorf("Expected key type RSA, got %s", key.KTY)
	}

	if key.Use != "sig" {
		t.Errorf("Expected key use sig, got %s", key.Use)
	}

	if key.Kid == "" {
		t.Error("Expected key ID, got empty")
	}

	if key.N == "" {
		t.Error("Expected key modulus, got empty")
	}

	if key.E == "" {
		t.Error("Expected key exponent, got empty")
	}
}
