package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"oidc-bridge/config"
	"oidc-bridge/model"
)

func ProxyToOPTokenEndpoint(req model.TokenRequest) (*model.OPTokenResponse, error) {
	// 构建请求参数
	form := url.Values{}
	form.Add("grant_type", req.GrantType)
	form.Add("code", req.Code)
	form.Add("redirect_uri", req.RedirectURI)
	form.Add("client_id", req.ClientID)
	form.Add("client_secret", req.ClientSecret)

	// 发送 POST 请求
	resp, err := http.PostForm(config.AppConfig.OPTokenURL, form)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to OP token endpoint: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var opResp model.OPTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&opResp); err != nil {
		return nil, fmt.Errorf("failed to decode OP token response: %v", err)
	}

	return &opResp, nil
}

func GetUserInfoFromOP(accessToken string) (map[string]interface{}, error) {
	// 创建请求
	req, err := http.NewRequest("GET", config.AppConfig.OPUserInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create userinfo request: %v", err)
	}

	// 添加 Authorization 头
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send userinfo request: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode userinfo response: %v", err)
	}

	return userInfo, nil
}