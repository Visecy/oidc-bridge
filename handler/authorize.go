package handler

import (
	"net/http"
	"net/url"
	"oidc-bridge/config"
	"oidc-bridge/service"
	"oidc-bridge/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandleAuthorize(c *gin.Context) {
	// 1. 验证参数
	clientID := c.Query("client_id")
	redirectURI := c.Query("redirect_uri")
	responseType := c.Query("response_type")
	scope := c.Query("scope")
	state := c.Query("state")
	nonce := c.Query("nonce")

	utils.DebugLogger.Printf("Handling authorize request for client: %s", clientID)

	if responseType != "code" {
		utils.ErrorLogger.Printf("Unsupported response type: %s for client: %s", responseType, clientID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported_response_type"})
		return
	}

	// 2. 处理 scope 映射
	hasOpenID := false
	var mappedScopes []string
	for _, s := range strings.Split(scope, ",") {
		s = strings.TrimSpace(s)
		if s == "openid" {
			hasOpenID = true
			continue
		}
		if mapped, ok := config.AppConfig.ScopeMapping[s]; ok {
			mappedScopes = append(mappedScopes, mapped)
		} else {
			mappedScopes = append(mappedScopes, s)
		}
	}

	// 3. 缓存 nonce
	if hasOpenID {
		if nonce == "" {
			utils.ErrorLogger.Printf("Nonce is required when scope includes openid for client: %s", clientID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "error_description": "nonce is required when scope includes openid"})
			return
		}
		// 存储 nonce 到 Redis
		err := service.SetNonce(clientID, redirectURI, nonce)
		if err != nil {
			utils.ErrorLogger.Printf("Failed to cache nonce for client: %s, error: %v", clientID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "error_description": "failed to cache nonce"})
			return
		}
		utils.DebugLogger.Printf("Nonce cached for client: %s", clientID)
	}

	// 4. 构建重定向 URL
	opAuthURL := config.AppConfig.OPAuthURL
	queryParams := url.Values{}
	queryParams.Add("response_type", "code")
	queryParams.Add("client_id", clientID)
	queryParams.Add("redirect_uri", redirectURI)
	queryParams.Add("scope", strings.Join(mappedScopes, ","))
	if state != "" {
		queryParams.Add("state", state)
	}
	if hasOpenID && nonce != "" {
		queryParams.Add("nonce", nonce)
	}

	// 构建完整 URL
	redirectURL, err := url.Parse(opAuthURL)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to parse OP auth URL for client: %s, error: %v", clientID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "error_description": "failed to parse OP auth URL"})
		return
	}
	redirectURL.RawQuery = queryParams.Encode()

	// 重定向到 OP
	utils.DebugLogger.Printf("Redirecting client: %s to OP", clientID)
	c.Redirect(http.StatusFound, redirectURL.String())
}
