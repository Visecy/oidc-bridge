package service

import (
	"time"

	"oidc-bridge/config"
	"oidc-bridge/utils"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateIDToken(issuer, clientID, redirectURI string, userInfo map[string]interface{}) (string, error) {
	// 1. 检查 nonce
	nonce, err := GetNonce(clientID, redirectURI)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to retrieve nonce: %v", err)
		return "", err
	}

	// 2. 构建 claims
	now := time.Now().Unix()
	claims := jwt.MapClaims{
		"iss":   issuer,
		"aud":   clientID,
		"exp":   now + int64(config.AppConfig.IDTokenLifetime),
		"iat":   now,
		"nonce": nonce,
	}

	// 3. 映射用户属性
	for opAttr, oidcClaim := range config.AppConfig.AttrMapping {
		if value, ok := userInfo[opAttr]; ok {
			claims[oidcClaim] = value
		}
	}

	// 4. 加载私钥
	privateKey, err := LoadPrivateKey()
	if err != nil {
		utils.ErrorLogger.Printf("Failed to load private key: %v", err)
		return "", err
	}

	// 5. 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// 6. 签名 token
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to sign token: %v", err)
		return "", err
	}

	utils.DebugLogger.Printf("Generated ID token for client: %s", clientID)
	return signedToken, nil
}
