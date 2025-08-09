package main

import (
	"flag"
	"oidc-bridge/config"
	"oidc-bridge/handler"
	"oidc-bridge/service"
	"oidc-bridge/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.InfoLogger.Println("Starting OIDC Bridge service")

	// 定义命令行参数
	configFile := flag.String("config", "config.yaml", "Path to config file")
	privateKeyPath := flag.String("private-key", "", "Path to private key file")
	publicKeyPath := flag.String("public-key", "", "Path to public key file")
	port := flag.String("port", "8080", "Port to run the server on")
	flag.Parse()

	// 1. 加载配置（支持通过命令行参数或环境变量指定配置文件和密钥路径）
	if err := config.LoadConfig(*configFile, *privateKeyPath, *publicKeyPath); err != nil {
		utils.ErrorLogger.Fatalf("Failed to load config: %v", err)
	}

	// 2. 初始化 Redis 或内存缓存
	service.InitRedis()

	// 3. 初始化 Gin
	r := gin.Default()

	// 4. 注册路由
	r.GET("/.well-known/openid-configuration", handler.HandleDiscovery)
	r.GET("/authorize", handler.HandleAuthorize)
	r.POST("/token", handler.HandleToken)
	r.GET("/userinfo", handler.HandleUserInfo)
	r.GET("/.well-known/jwks.json", handler.HandleJWKS)

	// 5. 启动服务
	serverAddr := ":" + *port
	utils.InfoLogger.Printf("Server starting on port %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		utils.ErrorLogger.Fatalf("Failed to start server: %v", err)
	}
}
