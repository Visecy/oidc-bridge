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
		issuer = c.Request.URL.Scheme + "://" + c.Request.URL.Host
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
