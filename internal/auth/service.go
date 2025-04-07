package auth

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"scope/internal/middleware"
	"scope/internal/models"
)

// 错误定义
var (
	ErrUserNotFound       = errors.New("用户不存在")
	ErrEmailAlreadyExists = errors.New("邮箱已被注册")
	ErrInvalidCredentials = errors.New("无效的凭证")
)

// TokenStore 定义令牌存储接口
type TokenStore interface {
	// AddToBlacklist 将令牌添加到黑名单
	AddToBlacklist(ctx context.Context, token string, expiry time.Time) error
	// IsBlacklisted 检查令牌是否在黑名单中
	IsBlacklisted(ctx context.Context, token string) (bool, error)
	// StoreRefreshToken 存储刷新令牌与用户ID的关联
	StoreRefreshToken(ctx context.Context, userID, token string, expiry time.Time) error
	// GetUserIDByRefreshToken 通过刷新令牌获取用户ID
	GetUserIDByRefreshToken(ctx context.Context, token string) (string, error)
	// RemoveRefreshToken 删除刷新令牌
	RemoveRefreshToken(ctx context.Context, token string) error
}

// AuthService 提供认证相关的功能
type AuthService struct {
	userStore    models.UserStore
	tokenService *middleware.TokenService
	tokenStore   TokenStore
}

// NewAuthService 创建一个新的认证服务
func NewAuthService(userStore models.UserStore, tokenService *middleware.TokenService, tokenStore TokenStore) *AuthService {
	return &AuthService{
		userStore:    userStore,
		tokenService: tokenService,
		tokenStore:   tokenStore,
	}
}

// RegisterUser 注册新用户
func (s *AuthService) RegisterUser(email, password, displayName string) (*models.User, error) {
	// 检查邮箱是否已存在
	_, err := s.userStore.FindByEmail(email)
	if err == nil {
		return nil, ErrEmailAlreadyExists
	} else if !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}

	// 对密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建新用户
	user := &models.User{
		Email:       email,
		Password:    string(hashedPassword),
		DisplayName: displayName,
	}

	if err := s.userStore.Create(user); err != nil {
		return nil, err
	}

	// 返回用户信息（不包含密码）
	return user, nil
}

// LoginUser 用户登录
func (s *AuthService) LoginUser(email, password string) (string, string, time.Time, error) {
	// 查找用户
	user, err := s.userStore.FindByEmail(email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return "", "", time.Time{}, ErrInvalidCredentials
		}
		return "", "", time.Time{}, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", time.Time{}, ErrInvalidCredentials
	}

	// 生成访问令牌
	accessToken, expiryTime, err := s.tokenService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return "", "", time.Time{}, err
	}

	// 生成刷新令牌
	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return "", "", time.Time{}, err
	}

	// 存储刷新令牌
	refreshExpiry := time.Now().Add(s.tokenService.GetRefreshTokenExpiry())
	ctx := context.Background()
	err = s.tokenStore.StoreRefreshToken(ctx, user.ID, refreshToken, refreshExpiry)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return accessToken, refreshToken, expiryTime, nil
}

// RefreshToken 刷新访问令牌
func (s *AuthService) RefreshToken(refreshToken string) (string, time.Time, error) {
	ctx := context.Background()

	// 检查令牌是否在黑名单中
	blacklisted, err := s.tokenStore.IsBlacklisted(ctx, refreshToken)
	if err != nil {
		return "", time.Time{}, err
	}
	if blacklisted {
		return "", time.Time{}, errors.New("刷新令牌已失效")
	}

	// 验证刷新令牌
	claims, err := s.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", time.Time{}, err
	}

	// 检查用户是否存在
	user, err := s.userStore.FindByID(claims.UserID)
	if err != nil {
		return "", time.Time{}, err
	}

	// 生成新的访问令牌
	accessToken, expiryTime, err := s.tokenService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return "", time.Time{}, err
	}

	return accessToken, expiryTime, nil
}

// LogoutUser 用户登出
func (s *AuthService) LogoutUser(refreshToken string) error {
	ctx := context.Background()

	// 验证刷新令牌
	claims, err := s.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	// 从存储中删除刷新令牌
	err = s.tokenStore.RemoveRefreshToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	// 将刷新令牌添加到黑名单
	expiryTime := time.Unix(claims.ExpiresAt.Unix(), 0)
	err = s.tokenStore.AddToBlacklist(ctx, refreshToken, expiryTime)
	if err != nil {
		return err
	}

	return nil
}
