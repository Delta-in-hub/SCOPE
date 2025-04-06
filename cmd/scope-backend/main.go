package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"scope/internal/auth"
	"scope/internal/middleware"
)

// 生成随机密钥
func generateRandomKey(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func main() {
	// 命令行参数
	port := flag.Int("port", 8080, "API服务端口")
	flag.Parse()

	// 生成随机密钥（在生产环境中应从配置或环境变量中获取）
	accessTokenSecret, err := generateRandomKey(32)
	if err != nil {
		log.Fatalf("生成访问令牌密钥失败: %v", err)
	}

	refreshTokenSecret, err := generateRandomKey(32)
	if err != nil {
		log.Fatalf("生成刷新令牌密钥失败: %v", err)
	}

	encryptionKey, err := generateRandomKey(32)
	if err != nil {
		log.Fatalf("生成加密密钥失败: %v", err)
	}

	// 创建令牌服务
	tokenConfig := middleware.TokenConfig{
		AccessTokenSecret:  accessTokenSecret,
		RefreshTokenSecret: refreshTokenSecret,
		AccessTokenExpiry:  time.Hour,          // 访问令牌有效期1小时
		RefreshTokenExpiry: time.Hour * 24 * 7, // 刷新令牌有效期7天
	}
	tokenService := middleware.NewTokenService(tokenConfig)

	// 创建用户存储
	userStore := auth.NewMemoryUserStore()

	// 创建认证服务
	authService := auth.NewAuthService(userStore, tokenService, encryptionKey)

	// 创建认证处理器
	authHandler := auth.NewHandler(authService)

	// 创建认证中间件
	authMiddleware := middleware.NewAuthMiddleware(tokenService)

	// 设置路由
	router := auth.SetupRouter(authHandler, authMiddleware)

	// 启动服务器
	serverAddr := fmt.Sprintf(":%d", *port)
	log.Printf("认证API服务启动在 http://localhost%s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}
