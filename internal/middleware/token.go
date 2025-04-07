package middleware

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenConfig 包含令牌生成的配置
type TokenConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

// TokenClaims 表示JWT令牌中的声明
type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// TokenService 处理JWT令牌的生成和验证
type TokenService struct {
	config TokenConfig
}

// NewTokenService 创建一个新的TokenService实例
func NewTokenService(config TokenConfig) *TokenService {
	return &TokenService{
		config: config,
	}
}

// GenerateAccessToken 生成访问令牌
func (s *TokenService) GenerateAccessToken(userID, email string) (string, time.Time, error) {
	expiryTime := time.Now().Add(s.config.AccessTokenExpiry)

	claims := TokenClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.config.AccessTokenSecret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("生成访问令牌失败: %w", err)
	}

	return signedToken, expiryTime, nil
}

// GenerateRefreshToken 生成刷新令牌
func (s *TokenService) GenerateRefreshToken(userID, email string) (string, error) {
	expiryTime := time.Now().Add(s.config.RefreshTokenExpiry)

	claims := TokenClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.config.RefreshTokenSecret))
	if err != nil {
		return "", fmt.Errorf("生成刷新令牌失败: %w", err)
	}

	return signedToken, nil
}

// ValidateAccessToken 验证访问令牌
func (s *TokenService) ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	return s.validateToken(tokenString, s.config.AccessTokenSecret)
}

// ValidateRefreshToken 验证刷新令牌
func (s *TokenService) ValidateRefreshToken(tokenString string) (*TokenClaims, error) {
	return s.validateToken(tokenString, s.config.RefreshTokenSecret)
}

// GetRefreshTokenExpiry 获取刷新令牌的过期时间
func (s *TokenService) GetRefreshTokenExpiry() time.Duration {
	return s.config.RefreshTokenExpiry
}

// validateToken 验证JWT令牌
func (s *TokenService) validateToken(tokenString, secret string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return []byte(secret), nil // hex.DecodeString(secret) ???
	})

	if err != nil {
		return nil, fmt.Errorf("解析令牌失败: %w", err)
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}
