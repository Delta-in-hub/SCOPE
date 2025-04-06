package auth

import (
	"encoding/json"
	"net/http"
	"time"
)

// 请求和响应结构体
type (
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	
	LoginResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	
	RegisterRequest struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		DisplayName string `json:"display_name"`
	}
	
	RegisterResponse struct {
		UserID      string `json:"user_id"`
		Email       string `json:"email"`
		DisplayName string `json:"display_name"`
	}
	
	RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token"`
	}
	
	RefreshTokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	
	LogoutRequest struct {
		RefreshToken string `json:"refresh_token"`
	}
)

// Handler 处理认证相关的请求
type Handler struct {
	authService *AuthService
}

// NewHandler 创建一个新的认证处理器
func NewHandler(authService *AuthService) *Handler {
	return &Handler{
		authService: authService,
	}
}

// Login 处理用户登录请求
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	// 只接受POST请求
	if r.Method != http.MethodPost {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}
	
	// 解析请求体
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}
	
	// 验证请求参数
	if req.Email == "" || req.Password == "" {
		http.Error(w, "邮箱和密码不能为空", http.StatusBadRequest)
		return
	}
	
	// 处理登录
	accessToken, refreshToken, expiryTime, err := h.authService.LoginUser(req.Email, req.Password)
	if err != nil {
		if err == ErrInvalidCredentials {
			http.Error(w, "无效的凭证", http.StatusUnauthorized)
		} else {
			http.Error(w, "登录失败", http.StatusInternalServerError)
		}
		return
	}
	
	// 计算令牌过期时间（秒）
	expiresIn := int(time.Until(expiryTime).Seconds())
	
	// 构建响应
	resp := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}
	
	// 返回JSON响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Register 处理用户注册请求
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	// 只接受POST请求
	if r.Method != http.MethodPost {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}
	
	// 解析请求体
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}
	
	// 验证请求参数
	if req.Email == "" || req.Password == "" || req.DisplayName == "" {
		http.Error(w, "邮箱、密码和显示名称不能为空", http.StatusBadRequest)
		return
	}
	
	// 处理注册
	user, err := h.authService.RegisterUser(req.Email, req.Password, req.DisplayName)
	if err != nil {
		if err == ErrEmailAlreadyExists {
			http.Error(w, "邮箱已被注册", http.StatusConflict)
		} else {
			http.Error(w, "注册失败", http.StatusInternalServerError)
		}
		return
	}
	
	// 构建响应
	resp := RegisterResponse{
		UserID:      user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
	}
	
	// 返回JSON响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// RefreshToken 处理刷新令牌请求
func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// 只接受POST请求
	if r.Method != http.MethodPost {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}
	
	// 解析请求体
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}
	
	// 验证请求参数
	if req.RefreshToken == "" {
		http.Error(w, "刷新令牌不能为空", http.StatusBadRequest)
		return
	}
	
	// 处理刷新令牌
	accessToken, expiryTime, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "刷新令牌失败", http.StatusUnauthorized)
		return
	}
	
	// 计算令牌过期时间（秒）
	expiresIn := int(time.Until(expiryTime).Seconds())
	
	// 构建响应
	resp := RefreshTokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   expiresIn,
	}
	
	// 返回JSON响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Logout 处理用户登出请求
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// 只接受POST请求
	if r.Method != http.MethodPost {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}
	
	// 解析请求体
	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}
	
	// 验证请求参数
	if req.RefreshToken == "" {
		http.Error(w, "刷新令牌不能为空", http.StatusBadRequest)
		return
	}
	
	// 处理登出
	if err := h.authService.LogoutUser(req.RefreshToken); err != nil {
		http.Error(w, "登出失败", http.StatusInternalServerError)
		return
	}
	
	// 返回成功响应（无内容）
	w.WriteHeader(http.StatusNoContent)
}
