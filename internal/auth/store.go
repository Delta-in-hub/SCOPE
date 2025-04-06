package auth

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"

	"scope/internal/models"
)

var (
	ErrUserNotFound       = errors.New("用户不存在")
	ErrEmailAlreadyExists = errors.New("邮箱已被注册")
	ErrInvalidCredentials = errors.New("无效的凭证")
)

// MemoryUserStore 实现了基于内存的用户存储
type MemoryUserStore struct {
	users map[string]*models.User
	// 用于快速通过邮箱查找用户
	emailIndex map[string]string
	mu         sync.RWMutex
}

// NewMemoryUserStore 创建一个新的内存用户存储
func NewMemoryUserStore() *MemoryUserStore {
	return &MemoryUserStore{
		users:      make(map[string]*models.User),
		emailIndex: make(map[string]string),
	}
}

// FindByID 通过ID查找用户
func (s *MemoryUserStore) FindByID(id string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}

	// 返回用户的副本以避免并发修改
	userCopy := *user
	return &userCopy, nil
}

// FindByEmail 通过邮箱查找用户
func (s *MemoryUserStore) FindByEmail(email string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userID, exists := s.emailIndex[email]
	if !exists {
		return nil, ErrUserNotFound
	}

	user := s.users[userID]
	userCopy := *user
	return &userCopy, nil
}

// Create 创建新用户
func (s *MemoryUserStore) Create(user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查邮箱是否已存在
	if _, exists := s.emailIndex[user.Email]; exists {
		return ErrEmailAlreadyExists
	}

	// 如果没有提供ID，则生成一个
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// 设置创建和更新时间
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// 存储用户并更新索引
	userCopy := *user
	s.users[user.ID] = &userCopy
	s.emailIndex[user.Email] = user.ID

	return nil
}

// Update 更新用户信息
func (s *MemoryUserStore) Update(user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existingUser, exists := s.users[user.ID]
	if !exists {
		return ErrUserNotFound
	}

	// 如果邮箱已更改，更新索引
	if existingUser.Email != user.Email {
		// 检查新邮箱是否已被其他用户使用
		if id, exists := s.emailIndex[user.Email]; exists && id != user.ID {
			return ErrEmailAlreadyExists
		}

		// 删除旧邮箱索引
		delete(s.emailIndex, existingUser.Email)
		// 添加新邮箱索引
		s.emailIndex[user.Email] = user.ID
	}

	// 更新时间戳
	user.UpdatedAt = time.Now()
	user.CreatedAt = existingUser.CreatedAt

	// 更新用户
	userCopy := *user
	s.users[user.ID] = &userCopy

	return nil
}
