package backend

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"scope/database/redis"
	"scope/internal/models"
	"time"

	"github.com/go-playground/validator/v10"
)

// 请求和响应结构体
type (
	LoginRequest struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	LoginResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}

	RegisterRequest struct {
		Email       string `json:"email" validate:"required"`
		Password    string `json:"password" validate:"required"`
		DisplayName string `json:"display_name" validate:"required"`
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
//
// @Summary      User login
// @Description  Authenticates a user and returns access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login credentials"
// @Router       /api/v1/auth/login [post]
// @Success      200 {object} LoginResponse
// @Failure      400 {object} string "Invalid request body"
// @Failure      401 {object} string "Invalid credentials"
// @Failure      500 {object} string "Login failed"
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	// 解析请求体
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}

	// 处理登录
	accessToken, refreshToken, expiryTime, err := h.authService.LoginUser(req.Email, req.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("登录失败: %v", err), http.StatusUnauthorized)
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
//
// @Summary      User registration
// @Description  Registers a new user and returns user information
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Registration information"
// @Router       /api/v1/auth/register [post]
// @Success      201 {object} RegisterResponse
// @Failure      400 {object} string "Invalid request body"
// @Failure      409 {object} string "Email already exists"
// @Failure      500 {object} string "Registration failed"
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {

	// 解析请求体
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}

	// 处理注册
	user, err := h.authService.RegisterUser(req.Email, req.Password, req.DisplayName)
	if err != nil {
		http.Error(w, fmt.Sprintf("注册失败: %v", err), http.StatusInternalServerError)
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
//
// @Summary      Refresh access token
// @Description  Uses a refresh token to generate a new access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RefreshTokenRequest true "Refresh token"
// @Router       /api/v1/auth/refreshToken [post]
// @Success      200 {object} RefreshTokenResponse
// @Failure      400 {object} string "Invalid request body or empty refresh token"
// @Failure      401 {object} string "Refresh token failed"
func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
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
//
// @Summary      User logout
// @Description  Invalidates the user's refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LogoutRequest true "Refresh token"
// @Router       /api/v1/auth/logout [post]
// @Security     ApiKeyAuth
// @Success      204 "No Content"
// @Failure      400 {object} string "Invalid request body or empty refresh token"
// @Failure      500 {object} string "Logout failed"
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
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

// NodeUp registers a node as online
//
// @Summary      Register node as online
// @Description  Updates a node's status to online and returns a token
// @Tags         node
// @Accept       json
// @Produce      json
// @Param        node body models.NodeInfo true "Node information"
// @Router       /api/v1/node/up [post]
// @Success      200 {object} map[string]string "Returns token"
// @Failure      400 {object} string "Invalid request body or incomplete node information"
// @Failure      500 {object} string "Failed to update node"
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

// NodeDown registers a node as offline
//
// @Summary      Register node as offline
// @Description  Updates a node's status to offline
// @Tags         node
// @Accept       json
// @Produce      json
// @Param        node body models.NodeInfo true "Node information"
// @Router       /api/v1/node/down [post]
// @Success      200 "OK"
// @Failure      400 {object} string "Invalid request body, incomplete node information, node doesn't exist, or token mismatch"
// @Failure      500 {object} string "Failed to update node"
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

// NodeList returns a list of all nodes
//
// @Summary      Get all nodes
// @Description  Returns a list of all registered nodes
// @Tags         node
// @Accept       json
// @Produce      json
// @Router       /api/v1/node/list [get]
// @Security     ApiKeyAuth
// @Success      200 {array} models.NodeInfo
// @Failure      500 {object} string "Failed to get node list"
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
