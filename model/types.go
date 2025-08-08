package model

type Config struct {
	OPAuthURL        string            `mapstructure:"op_authorize_url"`
	OPTokenURL       string            `mapstructure:"op_token_url"`
	OPUserInfoURL    string            `mapstructure:"op_userinfo_url"`
	Issuer           string            `mapstructure:"issuer"`
	IDTokenLifetime  int               `mapstructure:"id_token_lifetime"`
	NonceCacheTTL    int               `mapstructure:"nonce_cache_ttl"`
	SigningAlg       string            `mapstructure:"id_token_signing_alg"`
	ScopeMapping     map[string]string `mapstructure:"scope_mapping"`
	AttrMapping      map[string]string `mapstructure:"user_attribute_mapping"`
	RedisAddr        string            `mapstructure:"redis_addr"`
	PrivateKeyPath   string            `mapstructure:"private_key_path"`
	PublicKeyPath    string            `mapstructure:"public_key_path"`
}

type Discovery struct {
	Issuer                           string   `json:"issuer"`
	AuthorizationEndpoint            string   `json:"authorization_endpoint"`
	TokenEndpoint                    string   `json:"token_endpoint"`
	UserInfoEndpoint                 string   `json:"userinfo_endpoint"`
	JwksURI                          string   `json:"jwks_uri"`
	ScopesSupported                  []string `json:"scopes_supported"`
	ResponseTypesSupported           []string `json:"response_types_supported"`
	IDTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported"`
}

type TokenRequest struct {
	GrantType    string `form:"grant_type" json:"grant_type"`
	Code         string `form:"code" json:"code"`
	RedirectURI  string `form:"redirect_uri" json:"redirect_uri"`
	ClientID     string `form:"client_id" json:"client_id"`
	ClientSecret string `form:"client_secret" json:"client_secret"`
}

type OPTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Error        string `json:"error,omitempty"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
}

type JWK struct {
	KTY string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}