package handler

import (
	"net/http"
	"oidc-bridge/config"
	"oidc-bridge/model"

	"github.com/gin-gonic/gin"
)

func HandleDiscovery(c *gin.Context) {
	discovery := model.Discovery{
		Issuer:                           config.AppConfig.Issuer,
		AuthorizationEndpoint:            config.AppConfig.Issuer + "/authorize",
		TokenEndpoint:                    config.AppConfig.Issuer + "/token",
		UserInfoEndpoint:                 config.AppConfig.Issuer + "/userinfo",
		JwksURI:                          config.AppConfig.Issuer + "/.well-known/jwks.json",
		ScopesSupported:                  []string{"openid", "profile", "email"},
		ResponseTypesSupported:           []string{"code"},
		IDTokenSigningAlgValuesSupported: []string{config.AppConfig.SigningAlg},
	}
	c.JSON(http.StatusOK, discovery)
}