package handler

import (
	"fmt"
	"net/http"
	"oidc-bridge/config"
	"oidc-bridge/service"

	"github.com/gin-gonic/gin"
)

func HandleUserInfo(c *gin.Context) {
	// 1. 提取 access_token
	accessToken := c.GetHeader("Authorization")
	if accessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "error_description": "Authorization header is missing"})
		return
	}

	// 移除 "Bearer " 前缀
	if len(accessToken) > 7 && accessToken[:7] == "Bearer " {
		accessToken = accessToken[7:]
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "error_description": "Invalid Authorization header format"})
		return
	}

	// 2. 调用 OP 获取用户信息
	userInfo, err := service.GetUserInfoFromOP(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "error_description": fmt.Sprintf("failed to get user info from OP: %v", err)})
		return
	}

	// 3. 映射用户属性
	mappedUserInfo := make(map[string]interface{})
	for opAttr, oidcClaim := range config.AppConfig.AttrMapping {
		if value, ok := userInfo[opAttr]; ok {
			mappedUserInfo[oidcClaim] = value
		}
	}

	// 4. 保留未映射的属性
	for key, value := range userInfo {
		if _, mapped := config.AppConfig.AttrMapping[key]; !mapped {
			mappedUserInfo[key] = value
		}
	}

	c.JSON(http.StatusOK, mappedUserInfo)
}