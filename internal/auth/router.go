package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	authMiddleware "scope/internal/middleware"
)

// SetupRouter 配置认证API的路由
func SetupRouter(handler *Handler, authMiddleware *authMiddleware.AuthMiddleware) *chi.Mux {
	r := chi.NewRouter()

	// 全局中间件
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// 公共路由（无需认证）
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/login", handler.Login)
		r.Post("/register", handler.Register)
		r.Post("/refreshToken", handler.RefreshToken)

		// 需要认证的路由
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)
			r.Post("/logout", handler.Logout)
		})
	})

	return r
}
