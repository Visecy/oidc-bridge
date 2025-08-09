package handler

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net/http"
	"oidc-bridge/model"
	"oidc-bridge/service"

	"github.com/gin-gonic/gin"
)

func HandleJWKS(c *gin.Context) {
	// 1. 加载公钥
	publicKey, err := service.LoadPublicKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "error_description": fmt.Sprintf("failed to load public key: %v", err)})
		return
	}

	// 2. 构建 JWKS
	// 将公钥转换为 DER 格式
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "error_description": fmt.Sprintf("failed to marshal public key: %v", err)})
		return
	}

	// 构建 JWK
	jwk := model.JWK{
		KTY: "RSA",
		Use: "sig",
		Kid: "1",
		N:   base64.RawURLEncoding.EncodeToString(pubKeyBytes),
		E:   "AQAB", // 65537 的 Base64 URL 编码
	}

	// 构建 JWKS
	jwks := model.JWKS{
		Keys: []model.JWK{jwk},
	}

	c.JSON(http.StatusOK, jwks)
}
