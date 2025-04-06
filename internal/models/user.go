package models

import (
	"time"
)

// User 表示系统中的用户
type User struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Password    string    `json:"-"` // 不在JSON响应中显示密码
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserStore 定义用户存储接口
type UserStore interface {
	// FindByID 根据ID查找用户
	FindByID(id string) (*User, error)

	// FindByEmail 根据邮箱查找用户
	FindByEmail(email string) (*User, error)

	// Create 创建新用户
	Create(user *User) error

	// Update 更新用户信息
	Update(user *User) error
}
