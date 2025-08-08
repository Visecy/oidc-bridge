package handler

import (
	"net/http"
	"oidc-bridge/model"
	"oidc-bridge/service"
	"oidc-bridge/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandleToken(c *gin.Context) {
	// 1. 解析请求参数
	var req model.TokenRequest
	if err := c.ShouldBind(&req); err != nil {
		utils.ErrorLogger.Printf("Failed to bind token request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "error_description": err.Error()})
		return
	}

	// 获取 scope 参数
	scope := c.PostForm("scope")

	// 2. 向 OP 代理请求
	opResp, err := service.ProxyToOPTokenEndpoint(req)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to proxy to OP token endpoint: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "error_description": err.Error()})
		return
	}

	// 3. 构建响应
	resp := model.TokenResponse{
		AccessToken:  opResp.AccessToken,
		TokenType:    opResp.TokenType,
		RefreshToken: opResp.RefreshToken,
		ExpiresIn:    opResp.ExpiresIn,
	}

	// 4. 如果 scope 包含 openid，则生成 ID Token
	if strings.Contains(scope, "openid") {
		// 获取用户信息
		userInfo, err := service.GetUserInfoFromOP(opResp.AccessToken)
		if err != nil {
			utils.ErrorLogger.Printf("Failed to get user info: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "error_description": "failed to get user info"})
			return
		}

		// 生成 ID Token
		idToken, err := service.GenerateIDToken(req.ClientID, req.RedirectURI, userInfo)
		if err != nil {
			utils.ErrorLogger.Printf("Failed to generate ID token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "error_description": "failed to generate ID token"})
			return
		}

		resp.IDToken = idToken
	}

	c.JSON(http.StatusOK, resp)
}