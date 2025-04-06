package auth

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"scope/internal/middleware"
	"scope/internal/models"
)

// 黑名单存储已注销的刷新令牌
type RefreshTokenBlacklist struct {
	tokens map[string]time.Time
}

// NewRefreshTokenBlacklist 创建一个新的刷新令牌黑名单
func NewRefreshTokenBlacklist() *RefreshTokenBlacklist {
	return &RefreshTokenBlacklist{
		tokens: make(map[string]time.Time),
	}
}

// Add 将令牌添加到黑名单
func (b *RefreshTokenBlacklist) Add(token string, expiry time.Time) {
	b.tokens[token] = expiry
}

// IsBlacklisted 检查令牌是否在黑名单中
func (b *RefreshTokenBlacklist) IsBlacklisted(token string) bool {
	expiry, exists := b.tokens[token]
	if !exists {
		return false
	}

	// 如果令牌已过期，从黑名单中删除
	if time.Now().After(expiry) {
		delete(b.tokens, token)
		return false
	}

	return true
}

// AuthService 提供认证相关的功能
type AuthService struct {
	userStore      models.UserStore
	tokenService   *middleware.TokenService
	tokenBlacklist *RefreshTokenBlacklist
	encryptionKey  string // 用于加密敏感数据
}

// NewAuthService 创建一个新的认证服务
func NewAuthService(userStore models.UserStore, tokenService *middleware.TokenService, encryptionKey string) *AuthService {
	return &AuthService{
		userStore:      userStore,
		tokenService:   tokenService,
		tokenBlacklist: NewRefreshTokenBlacklist(),
		encryptionKey:  encryptionKey,
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

	return accessToken, refreshToken, expiryTime, nil
}

// RefreshToken 刷新访问令牌
func (s *AuthService) RefreshToken(refreshToken string) (string, time.Time, error) {
	// 检查令牌是否在黑名单中
	if s.tokenBlacklist.IsBlacklisted(refreshToken) {
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
	// 验证刷新令牌
	claims, err := s.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	// 将刷新令牌添加到黑名单
	expiryTime := time.Unix(claims.ExpiresAt.Unix(), 0)
	s.tokenBlacklist.Add(refreshToken, expiryTime)

	return nil
}
