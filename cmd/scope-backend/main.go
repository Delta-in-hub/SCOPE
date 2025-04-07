package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"scope/database/postgres"
	"scope/database/redis"
	"scope/internal/auth"
	"scope/internal/middleware"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsIntOrDefault 获取环境变量并转换为整数，如果不存在或转换失败则返回默认值
func getEnvAsIntOrDefault(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

func main() {
	// 加载环境变量
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("加载环境变量失败: %v", err)
	}
	// 命令行参数
	port := flag.Int("port", 18080, "API服务端口")
	flag.Parse()

	// 从环境变量获取密钥
	accessTokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	if accessTokenSecret == "" {
		log.Fatalf("访问令牌密钥未设置")
	}

	refreshTokenSecret := os.Getenv("REFRESH_TOKEN_SECRET")
	if refreshTokenSecret == "" {
		log.Fatalf("刷新令牌密钥未设置")
	}

	// 创建令牌服务
	tokenConfig := middleware.TokenConfig{
		AccessTokenSecret:  accessTokenSecret,
		RefreshTokenSecret: refreshTokenSecret,
		AccessTokenExpiry:  time.Hour,          // 访问令牌有效期1小时
		RefreshTokenExpiry: time.Hour * 24 * 7, // 刷新令牌有效期7天
	}
	tokenService := middleware.NewTokenService(tokenConfig)

	// 初始化PostgreSQL连接
	dbConfig := postgres.Config{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvAsIntOrDefault("DB_PORT", 5432),
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName:   getEnvOrDefault("DB_NAME", "scope"),
		SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
	}

	db, err := postgres.NewDB(dbConfig)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	// 初始化数据库用户表
	if err := postgres.InitDB4User(db); err != nil {
		log.Fatalf("初始化数据库用户表失败: %v", err)
	}

	// 初始化Redis连接
	redisConfig := redis.Config{
		Addr:     getEnvOrDefault("REDIS_ADDR", "localhost:6379"),
		Password: getEnvOrDefault("REDIS_PASSWORD", ""),
		DB:       getEnvAsIntOrDefault("REDIS_DB", 0),
	}

	redisClient, err := redis.NewClient(redisConfig)
	if err != nil {
		log.Fatalf("连接Redis失败: %v", err)
	}
	defer redisClient.Close()

	// 创建令牌存储
	tokenStore := redis.NewTokenStore(redisClient)

	// 创建用户存储
	userStore := postgres.NewUserStore(db)

	// 创建认证服务
	authService := auth.NewAuthService(userStore, tokenService, tokenStore)

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
