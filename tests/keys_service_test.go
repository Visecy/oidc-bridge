package tests

import (
	"oidc-bridge/config"
	"oidc-bridge/model"
	"oidc-bridge/service"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func init() {
	// 加载配置
	viper.SetConfigName("config_test")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	config.AppConfig = &model.Config{}
	if err := viper.Unmarshal(config.AppConfig); err != nil {
		panic(err)
	}
}

func TestLoadPrivateKey(t *testing.T) {
	// 确保密钥文件存在
	if _, err := os.Stat(config.AppConfig.PrivateKeyPath); os.IsNotExist(err) {
		t.Skip("Private key file not found, skipping private key loading tests")
	}

	// 调用函数
	key, err := service.LoadPrivateKey()
	if err != nil {
		t.Errorf("Failed to load private key: %v", err)
	}

	if key == nil {
		t.Error("Expected private key, got nil")
	}
}

func TestLoadPublicKey(t *testing.T) {
	// 确保密钥文件存在
	if _, err := os.Stat(config.AppConfig.PublicKeyPath); os.IsNotExist(err) {
		t.Skip("Public key file not found, skipping public key loading tests")
	}

	// 调用函数
	key, err := service.LoadPublicKey()
	if err != nil {
		t.Errorf("Failed to load public key: %v", err)
	}

	if key == nil {
		t.Error("Expected public key, got nil")
	}
}
