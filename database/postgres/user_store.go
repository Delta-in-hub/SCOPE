package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"scope/internal/models"
)

var (
	ErrUserNotFound       = errors.New("用户不存在")
	ErrEmailAlreadyExists = errors.New("邮箱已被注册")
)

// UserStore 实现了基于PostgreSQL的用户存储
type UserStore struct {
	db *sqlx.DB
}

// NewUserStore 创建一个新的PostgreSQL用户存储
func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

// FindByID 通过ID查找用户
func (s *UserStore) FindByID(id string) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, password, display_name, created_at, updated_at FROM users WHERE id = $1`
	err := s.db.Get(&user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}

// FindByEmail 通过邮箱查找用户
func (s *UserStore) FindByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, password, display_name, created_at, updated_at FROM users WHERE email = $1`
	err := s.db.Get(&user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}

// Create 创建新用户
func (s *UserStore) Create(user *models.User) error {
	// 检查邮箱是否已存在
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = $1`
	err := s.db.Get(&count, query, user.Email)
	if err != nil {
		return fmt.Errorf("检查邮箱是否存在失败: %w", err)
	}
	if count > 0 {
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

	// 插入用户
	insertQuery := `
		INSERT INTO users (id, email, password, display_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = s.db.Exec(insertQuery, user.ID, user.Email, user.Password, user.DisplayName, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}

	return nil
}

// Update 更新用户信息
func (s *UserStore) Update(user *models.User) error {
	// 检查用户是否存在
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	err := s.db.Get(&exists, checkQuery, user.ID)
	if err != nil {
		return fmt.Errorf("检查用户是否存在失败: %w", err)
	}
	if !exists {
		return ErrUserNotFound
	}

	// 如果邮箱已更改，检查新邮箱是否已被其他用户使用
	if user.Email != "" {
		var count int
		emailQuery := `SELECT COUNT(*) FROM users WHERE email = $1 AND id != $2`
		err := s.db.Get(&count, emailQuery, user.Email, user.ID)
		if err != nil {
			return fmt.Errorf("检查邮箱是否可用失败: %w", err)
		}
		if count > 0 {
			return ErrEmailAlreadyExists
		}
	}

	// 更新时间戳
	user.UpdatedAt = time.Now()

	// 更新用户
	updateQuery := `
		UPDATE users
		SET email = $1, password = $2, display_name = $3, updated_at = $4
		WHERE id = $5
	`
	_, err = s.db.Exec(updateQuery, user.Email, user.Password, user.DisplayName, user.UpdatedAt, user.ID)
	if err != nil {
		return fmt.Errorf("更新用户失败: %w", err)
	}

	return nil
}

// InitDB4User 初始化数据库表
func InitDB4User(db *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(36) PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		display_name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, schema)
	if err != nil {
		return fmt.Errorf("初始化用户表失败: %w", err)
	}

	return nil
}
