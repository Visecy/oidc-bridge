package config

import (
	"oidc-bridge/model"
	"oidc-bridge/utils"
	"os"

	"github.com/spf13/viper"
)

var AppConfig *model.Config

func LoadConfig(configFile, privateKeyPath, publicKeyPath string) error {
	// 检查环境变量是否指定了配置文件路径
	if envConfigFile := os.Getenv("CONFIG_FILE"); envConfigFile != "" {
		configFile = envConfigFile
		utils.DebugLogger.Printf("Config file path overridden by environment variable: %s", configFile)
	} else {
		utils.DebugLogger.Printf("Loading config from file: %s", configFile)
	}

	// 使用 viper 加载配置文件
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		utils.ErrorLogger.Printf("Failed to read config file: %v", err)
		return err
	}

	utils.InfoLogger.Println("Config file loaded successfully")

	AppConfig = &model.Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		utils.ErrorLogger.Printf("Failed to unmarshal config: %v", err)
		return err
	}

	// 优先级：命令行参数 > 环境变量 > 配置文件
	// 先处理环境变量，若有值则覆盖配置文件中的设置
	if envPrivateKeyPath := os.Getenv("PRIVATE_KEY_PATH"); envPrivateKeyPath != "" {
		AppConfig.PrivateKeyPath = envPrivateKeyPath
		utils.DebugLogger.Printf("Private key path overridden by environment variable: %s", envPrivateKeyPath)
	}
	if envPublicKeyPath := os.Getenv("PUBLIC_KEY_PATH"); envPublicKeyPath != "" {
		AppConfig.PublicKeyPath = envPublicKeyPath
		utils.DebugLogger.Printf("Public key path overridden by environment variable: %s", envPublicKeyPath)
	}
	if envRedisAddr := os.Getenv("REDIS_ADDR"); envRedisAddr != "" {
		AppConfig.RedisAddr = envRedisAddr
		utils.DebugLogger.Printf("Redis address overridden by environment variable: %s", envRedisAddr)
	}

	// 再处理命令行参数，若有值则覆盖环境变量和配置文件中的设置
	if privateKeyPath != "" {
		AppConfig.PrivateKeyPath = privateKeyPath
		utils.DebugLogger.Printf("Private key path overridden by command-line argument: %s", privateKeyPath)
	}
	if publicKeyPath != "" {
		AppConfig.PublicKeyPath = publicKeyPath
		utils.DebugLogger.Printf("Public key path overridden by command-line argument: %s", publicKeyPath)
	}

	utils.InfoLogger.Println("Configuration loaded successfully")
	return nil
}
