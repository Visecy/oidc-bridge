package tests

import (
	"oidc-bridge/config"
	"oidc-bridge/service"
	"testing"
)

// TestNestedMappingConfig 测试嵌套映射配置是否正确加载
func TestNestedMappingConfig(t *testing.T) {
	// 准备测试配置文件路径
	configFile := "nested_mapping_example.yaml"

	// 加载配置
	if err := config.LoadConfig(configFile, "./private.key", "./public.key"); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 验证嵌套映射配置是否正确加载
	expectedMappings := map[string]string{
		"sub":              "sub",
		"data::email":       "email",
		"data::name":        "name",
		"data::avatar_url":  "picture",
	}

	// 检查映射是否正确
	for nestedKey, expectedValue := range expectedMappings {
		if actualValue, exists := config.AppConfig.AttrMapping[nestedKey]; !exists {
			t.Errorf("Expected nested mapping key '%s' not found", nestedKey)
		} else if actualValue != expectedValue {
			t.Errorf("Expected nested mapping for key '%s' to be '%s', got '%s'", nestedKey, expectedValue, actualValue)
		}
	}

	// 验证普通映射仍然有效
	if actualValue, exists := config.AppConfig.AttrMapping["sub"]; !exists {
		t.Error("Expected mapping key 'sub' not found")
	} else if actualValue != "sub" {
		t.Errorf("Expected mapping for key 'sub' to be 'sub', got '%s'", actualValue)
	}
}

// TestStandardMappingConfig 测试标准映射配置是否正确加载
func TestStandardMappingConfig(t *testing.T) {
	// 准备测试配置文件路径
	configFile := "config_test.yaml"

	// 加载配置
	if err := config.LoadConfig(configFile, "./private.key", "./public.key"); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 验证标准映射配置是否正确加载
	expectedMappings := map[string]string{
		"sub":   "sub",
		"name":  "name",
		"email": "email",
	}

	// 检查映射是否正确
	for key, expectedValue := range expectedMappings {
		if actualValue, exists := config.AppConfig.AttrMapping[key]; !exists {
			t.Errorf("Expected mapping key '%s' not found", key)
		} else if actualValue != expectedValue {
			t.Errorf("Expected mapping for key '%s' to be '%s', got '%s'", key, expectedValue, actualValue)
		}
	}
}

// TestNestedMappingFunctionality 测试嵌套映射功能是否正常工作
func TestNestedMappingFunctionality(t *testing.T) {
	// 准备测试配置文件路径
	configFile := "nested_mapping_example.yaml"

	// 加载配置
	if err := config.LoadConfig(configFile, "./private.key", "./public.key"); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 准备测试用的用户信息数据
	userInfo := map[string]interface{}{
		"sub":  "12345",
		"data": map[string]interface{}{
			"email":      "user@example.com",
			"name":       "John Doe",
			"avatar_url": "https://example.com/avatar.jpg",
		},
	}

	// 验证嵌套属性映射功能
	mappedUserInfo := make(map[string]interface{})
	for opAttr, oidcClaim := range config.AppConfig.AttrMapping {
		if value, ok := service.GetNestedValue(userInfo, opAttr); ok {
			mappedUserInfo[oidcClaim] = value
		}
	}

	// 检查映射结果
	expectedResults := map[string]interface{}{
		"sub":     "12345",
		"email":   "user@example.com",
		"name":    "John Doe",
		"picture": "https://example.com/avatar.jpg",
	}

	for claim, expectedValue := range expectedResults {
		if actualValue, exists := mappedUserInfo[claim]; !exists {
			t.Errorf("Expected mapped claim '%s' not found", claim)
		} else if actualValue != expectedValue {
			t.Errorf("Expected mapped value for claim '%s' to be '%s', got '%s'", claim, expectedValue, actualValue)
		}
	}
}

// TestStandardMappingFunctionality 测试标准映射功能是否正常工作
func TestStandardMappingFunctionality(t *testing.T) {
	// 准备测试配置文件路径
	configFile := "config_test.yaml"

	// 加载配置
	if err := config.LoadConfig(configFile, "./private.key", "./public.key"); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 准备测试用的用户信息数据
	userInfo := map[string]interface{}{
		"sub":   "12345",
		"name":  "John Doe",
		"email": "user@example.com",
	}

	// 验证标准属性映射功能
	mappedUserInfo := make(map[string]interface{})
	for opAttr, oidcClaim := range config.AppConfig.AttrMapping {
		if value, ok := service.GetNestedValue(userInfo, opAttr); ok {
			mappedUserInfo[oidcClaim] = value
		}
	}

	// 检查映射结果
	expectedResults := map[string]interface{}{
		"sub":   "12345",
		"name":  "John Doe",
		"email": "user@example.com",
	}

	for claim, expectedValue := range expectedResults {
		if actualValue, exists := mappedUserInfo[claim]; !exists {
			t.Errorf("Expected mapped claim '%s' not found", claim)
		} else if actualValue != expectedValue {
			t.Errorf("Expected mapped value for claim '%s' to be '%s', got '%s'", claim, expectedValue, actualValue)
		}
	}
}

// TestGetNestedValue 测试GetNestedValue函数
func TestGetNestedValue(t *testing.T) {
	// 准备测试数据
	data := map[string]interface{}{
		"level1": map[string]interface{}{
			"level2": map[string]interface{}{
				"value": "test",
			},
		},
		"simple": "value",
	}

	// 测试双冒号分隔符
	if value, ok := service.GetNestedValue(data, "level1::level2::value"); !ok {
		t.Error("Failed to get nested value with '::' separator")
	} else if value != "test" {
		t.Errorf("Expected 'test', got '%s'", value)
	}

	// 测试点号分隔符
	if value, ok := service.GetNestedValue(data, "level1.level2.value"); !ok {
		t.Error("Failed to get nested value with '.' separator")
	} else if value != "test" {
		t.Errorf("Expected 'test', got '%s'", value)
	}

	// 测试简单值
	if value, ok := service.GetNestedValue(data, "simple"); !ok {
		t.Error("Failed to get simple value")
	} else if value != "value" {
		t.Errorf("Expected 'value', got '%s'", value)
	}

	// 测试不存在的路径
	if _, ok := service.GetNestedValue(data, "level1::nonexistent"); ok {
		t.Error("Should not find value for nonexistent path")
	}
}