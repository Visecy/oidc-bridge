package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"oidc-bridge/config"
	"oidc-bridge/utils"
	"time"
)

var (
	RedisClient *redis.Client
	useRedis     bool
)

func InitRedis() {
	// 如果配置了Redis地址，则初始化Redis客户端
	if config.AppConfig.RedisAddr != "" {
		RedisClient = redis.NewClient(&redis.Options{
			Addr: config.AppConfig.RedisAddr,
		})
		
		// 测试Redis连接
		_, err := RedisClient.Ping(context.Background()).Result()
		if err != nil {
			utils.ErrorLogger.Printf("Failed to connect to Redis: %v", err)
			utils.InfoLogger.Println("Falling back to memory cache")
			useRedis = false
			InitMemoryCache()
		} else {
			utils.InfoLogger.Println("Successfully connected to Redis")
			useRedis = true
		}
	} else {
		fmt.Println("Redis not configured, using memory cache")
		useRedis = false
		InitMemoryCache()
	}
}

func SetNonce(clientID, redirectURI, nonce string) error {
	cacheKey := "nonce:" + clientID + ":" + redirectURI
	
	if useRedis {
		return RedisClient.Set(context.Background(), cacheKey, nonce, time.Duration(config.AppConfig.NonceCacheTTL)*time.Second).Err()
	} else {
		// 使用本地内存缓存
		ttl := time.Duration(config.AppConfig.NonceCacheTTL) * time.Second
		GlobalMemoryCache.Set(cacheKey, nonce, ttl)
		utils.DebugLogger.Printf("Set nonce in memory cache: %s", cacheKey)
		return nil
	}
}

func GetNonce(clientID, redirectURI string) (string, error) {
	cacheKey := "nonce:" + clientID + ":" + redirectURI
	
	if useRedis {
		return RedisClient.Get(context.Background(), cacheKey).Result()
	} else {
		// 使用本地内存缓存
		if value, exists := GlobalMemoryCache.Get(cacheKey); exists {
			utils.DebugLogger.Printf("Get nonce from memory cache: %s", cacheKey)
			return value, nil
		}
		utils.DebugLogger.Printf("Nonce not found in memory cache: %s", cacheKey)
		return "", errors.New("nonce not found")
	}
}