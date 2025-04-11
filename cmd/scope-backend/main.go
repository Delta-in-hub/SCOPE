package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"scope/database/postgres"
	"scope/database/redis"
	_ "scope/docs/backend"
	"scope/internal/backend"
	"scope/internal/middleware"
	"scope/internal/utils"
	"sync"
	"time"

	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Scope Center Backend API
// @version 1.0
// @description Scope Center Backend API
// @termsOfService http://swagger.io/terms/

// @contact.name Delta
// @contact.url https://github.com/Delta-in-hub/ebpf-golang
// @contact.email DeltaMail@qq.com

// @host 127.0.0.1:18080
func main() {
	// 加载环境变量
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("加载环境变量失败: %v", err)
	}
	// 命令行参数
	port := flag.Int("port", 18080, "API服务端口")
	verbose := flag.Bool("verbose", false, "是否启用详细输出")
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

	// 初始化PostgreSQL连接
	dbConfig := postgres.Config{
		Host:     utils.GetEnvOrDefault("DB_HOST", "localhost"),
		Port:     utils.GetEnvAsIntOrDefault("DB_PORT", 5432),
		User:     utils.GetEnvOrDefault("DB_USER", "postgres"),
		Password: utils.GetEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName:   utils.GetEnvOrDefault("DB_NAME", "scope"),
		SSLMode:  utils.GetEnvOrDefault("DB_SSLMODE", "disable"),
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
		Addr:     utils.GetEnvOrDefault("REDIS_ADDR", "localhost:6379"),
		Password: utils.GetEnvOrDefault("REDIS_PASSWORD", ""),
		DB:       0, // 0 for user
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
	tokenService := middleware.NewTokenService(tokenConfig)
	authService := backend.NewAuthService(userStore, tokenService, tokenStore)

	// 创建认证处理器
	redisconfig4node := redisConfig
	redisconfig4node.DB = 2 // 2 for Node Stroe

	backendHandler := backend.NewHandler(authService, redisconfig4node)

	// 创建认证中间件
	middleware := middleware.NewAuthMiddleware(tokenService)

	// 设置路由
	router := backend.SetupRouter(backendHandler, middleware)

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:18080/swagger/doc.json"), //The url pointing to API definition
	))

	// 启动服务器
	serverAddr := fmt.Sprintf(":%d", *port)
	log.Printf("认证API服务启动在 http://localhost%s", serverAddr)

	// 接收Redis Stream 来自 agent

	timescaledb, err := postgres.NewDB(dbConfig)
	if err != nil {
		log.Fatalf("连接TimescaleDB失败: %v", err)
	}
	defer timescaledb.Close()

	initCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := postgres.InitializeTSDBSchema(initCtx, timescaledb); err != nil {
		log.Fatalf("初始化 TimescaleDB schema 失败: %v", err)
	}

	streamConfig := redis.Config{
		Addr:     utils.GetEnvOrDefault("REDIS_ADDR", "localhost:6379"),
		Password: utils.GetEnvOrDefault("REDIS_PASSWORD", ""),
		DB:       1, // 1 for stream messages queue
	}

	streamClient, err := redis.NewClient(streamConfig)
	if err != nil {
		log.Fatalf("连接Redis失败 For Stream: %v", err)
	}
	defer streamClient.Close()

	var wg sync.WaitGroup

	// 启动Node Ping Checker
	wg.Add(1)
	go backend.NodePingChecker(&wg, backendHandler)

	var cpunum = runtime.NumCPU() / 2
	if cpunum < 1 {
		cpunum = 1
	}

	for k := range cpunum {
		wg.Add(1)
		go backend.Receive(context.Background(), &wg, timescaledb, streamClient, *verbose, k)
	}

	wg.Add(1)
	go backend.XDelMessages(context.Background(), streamClient, *verbose)

	log.Fatal(http.ListenAndServe(serverAddr, router))
	wg.Wait()
}
