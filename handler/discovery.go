package handler

import (
	"net/http"
	"oidc-bridge/config"
	"oidc-bridge/model"

	"github.com/gin-gonic/gin"
)

func HandleDiscovery(c *gin.Context) {
	// 如果配置中没有提供 Issuer，则从请求的 URL 中获取
	issuer := config.AppConfig.Issuer
	if issuer == "" {
		// 获取请求的协议
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		} else if c.Request.Header.Get("X-Forwarded-Proto") != "" {
			scheme = c.Request.Header.Get("X-Forwarded-Proto")
		}

		// 获取请求的主机名
		host := c.Request.Host
		// 如果 host 为空（如在测试环境中），使用默认值
		if host == "" {
			host = "example.com"
		}

		issuer = scheme + "://" + host
	}

	discovery := model.Discovery{
		Issuer:                           issuer,
		AuthorizationEndpoint:            issuer + "/authorize",
		TokenEndpoint:                    issuer + "/token",
		UserInfoEndpoint:                 issuer + "/userinfo",
		JwksURI:                          issuer + "/.well-known/jwks.json",
		ScopesSupported:                  []string{"openid", "profile", "email"},
		ResponseTypesSupported:           []string{"code"},
		IDTokenSigningAlgValuesSupported: []string{config.AppConfig.SigningAlg},
	}
	c.JSON(http.StatusOK, discovery)
}
