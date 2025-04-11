package backend

import (
	"encoding/json"
	"net/http"

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
	r.Use(middleware.Heartbeat("/health"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// Welcome message
		w.Write([]byte("Welcome to Scope Cetner Backend\nSee More at https://github.com/Delta-in-hub/ebpf-golang\n"))
	})

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

	r.Route("/api/v1/node", func(r chi.Router) {
		r.Post("/up", handler.nodeHandler.NodeUp)
		r.Post("/down", handler.nodeHandler.NodeDown)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)
			r.Get("/list", handler.nodeHandler.NodeList)
		})
	})

	// 新增的/apis路由，返回所有路由信息
	r.Get("/apis", func(w http.ResponseWriter, req *http.Request) {
		type RouteInfo struct {
			Method string `json:"method"`
			Path   string `json:"path"`
		}
		var routes []RouteInfo
		chi.Walk(r, func(method string, path string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
			route := RouteInfo{
				Method: method,
				Path:   path,
			}
			routes = append(routes, route)
			return nil
		})
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(routes); err != nil {
			http.Error(w, "Failed to encode routes", http.StatusInternalServerError)
		}
	})

	return r
}
