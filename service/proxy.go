package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

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

// GetNestedValue 从嵌套的 map 中获取值
// 支持两种分隔符：点号(.)和双冒号(::)
func GetNestedValue(data map[string]interface{}, path string) (interface{}, bool) {
	// 首先尝试使用双冒号分隔符分割路径
	parts := strings.Split(path, "::")

	// 如果没有双冒号分隔符，则使用点号分隔符
	if len(parts) <= 1 {
		parts = strings.Split(path, ".")
	}

	current := data

	// 遍历路径的每个部分
	for i, part := range parts {
		// 如果是最后一部分，直接返回值
		if i == len(parts)-1 {
			if value, ok := current[part]; ok {
				return value, true
			}
			return nil, false
		}

		// 如果不是最后一部分，确保当前值是一个 map
		if next, ok := current[part].(map[string]interface{}); ok {
			current = next
		} else {
			return nil, false
		}
	}

	return nil, false
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

	// 映射用户属性
	mappedUserInfo := make(map[string]interface{})
	for opAttr, oidcClaim := range config.AppConfig.AttrMapping {
		if value, ok := GetNestedValue(userInfo, opAttr); ok {
			mappedUserInfo[oidcClaim] = value
		}
	}

	// 保留未映射的属性
	for key, value := range userInfo {
		if _, mapped := config.AppConfig.AttrMapping[key]; !mapped {
			mappedUserInfo[key] = value
		}
	}

	return mappedUserInfo, nil
}
