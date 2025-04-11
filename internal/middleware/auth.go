package middleware

import (
	"context"
	"net/http"
	"strings"
	// "scope/internal/middleware"
)

// 用于在上下文中存储用户信息的键
type contextKey string

const (
	UserIDKey contextKey = "user_id"
	EmailKey  contextKey = "email"
)

// AuthMiddleware 用于验证JWT令牌的中间件
type AuthMiddleware struct {
	tokenService *TokenService
}

// NewAuthMiddleware 创建一个新的认证中间件
func NewAuthMiddleware(tokenService *TokenService) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: tokenService,
	}
}

// Authenticate 验证请求中的JWT令牌
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从Authorization头中获取令牌
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "未提供认证令牌", http.StatusUnauthorized)
			return
		}

		// 检查Authorization头的格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "认证头格式无效", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// 验证令牌
		claims, err := m.tokenService.ValidateAccessToken(tokenString)
		if err != nil {
			http.Error(w, "无效的认证令牌", http.StatusUnauthorized)
			return
		}

		// 将用户信息添加到请求上下文中
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, EmailKey, claims.Email)

		// 使用更新后的上下文继续处理请求
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID 从请求上下文中获取用户ID
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

// GetEmail 从请求上下文中获取用户邮箱
func GetEmail(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(EmailKey).(string)
	return email, ok
}
