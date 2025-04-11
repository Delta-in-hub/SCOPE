package backend

import (
	"encoding/json"
	"log"
	"net/http"
	"scope/database/redis"
	"scope/internal/models"
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

type NodeHandler struct {
	nodeService *NodeService
}

// Handler 处理认证相关的请求
type Handler struct {
	authService *AuthService
	nodeHandler *NodeHandler
}

// NewHandler 创建一个新的认证处理器
func NewHandler(authService *AuthService, redisconf4node redis.Config) *Handler {
	handler := Handler{
		authService: authService,
	}
	if redisconf4node.DB != 2 {
		redisconf4node.DB = 2 // 2 for Node Stroe
	}
	client, _ := redis.NewClient(redisconf4node)
	nodestore := redis.NewNodeStore(client)
	nodeservice := NodeService{
		nodeStore: nodestore,
	}
	handler.nodeHandler = &NodeHandler{
		nodeService: &nodeservice,
	}
	return &handler
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

func (h *NodeHandler) NodeUp(w http.ResponseWriter, r *http.Request) {

	var node models.NodeInfo
	if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}
	log.Printf("NodeInfo: %v", node)
	if node.Status != "online" || node.ID == "" || node.IPs == nil || len(node.IPs) == 0 || node.LastSeen.IsZero() {
		http.Error(w, "节点信息不完整", http.StatusBadRequest)
		return
	}

	token, err := h.nodeService.NodeUp(r.Context(), node)
	if err != nil {
		http.Error(w, "更新节点失败", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *NodeHandler) NodeDown(w http.ResponseWriter, r *http.Request) {
	var node models.NodeInfo
	if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}
	log.Printf("NodeInfo: %v", node)
	if node.ID == "" || node.Status != "offline" {
		http.Error(w, "节点信息不完整", http.StatusBadRequest)
		return
	}

	noderedis, err := h.nodeService.GetNode(r.Context(), node.ID)
	if err != nil {
		http.Error(w, "节点不存在", http.StatusBadRequest)
		return
	}

	if noderedis.Token != node.Token {
		http.Error(w, "token不匹配", http.StatusBadRequest)
		return
	}

	if err := h.nodeService.NodeDown(r.Context(), node); err != nil {
		http.Error(w, "更新节点失败", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *NodeHandler) NodeList(w http.ResponseWriter, r *http.Request) {
	nodes, err := h.nodeService.ListNodes(r.Context())
	if err != nil {
		http.Error(w, "获取节点列表失败", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nodes)
}
