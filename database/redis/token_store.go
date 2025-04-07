package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// TokenStore 实现基于Redis的令牌存储
type TokenStore struct {
	client *redis.Client
}

// NewTokenStore 创建一个新的Redis令牌存储
func NewTokenStore(client *redis.Client) *TokenStore {
	return &TokenStore{
		client: client,
	}
}

// AddToBlacklist 将令牌添加到黑名单
func (s *TokenStore) AddToBlacklist(ctx context.Context, token string, expiry time.Time) error {
	// 计算过期时间（从现在到过期时间的持续时间）
	duration := time.Until(expiry)
	if duration <= 0 {
		// 如果令牌已过期，不需要添加到黑名单
		return nil
	}

	// 将令牌添加到黑名单，并设置相同的过期时间
	key := fmt.Sprintf("blacklist:%s", token)
	err := s.client.Set(ctx, key, "1", duration).Err()
	if err != nil {
		return fmt.Errorf("添加令牌到黑名单失败: %w", err)
	}

	return nil
}

// IsBlacklisted 检查令牌是否在黑名单中
func (s *TokenStore) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", token)
	result, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("检查令牌是否在黑名单中失败: %w", err)
	}

	return result > 0, nil
}

// StoreRefreshToken 存储刷新令牌与用户ID的关联
func (s *TokenStore) StoreRefreshToken(ctx context.Context, userID, token string, expiry time.Time) error {
	// 存储令牌到用户的映射（用于跟踪用户的活跃会话）
	userKey := fmt.Sprintf("user:%s:tokens", userID)
	tokenKey := fmt.Sprintf("token:%s", token)

	// 使用管道批量执行命令
	pipe := s.client.Pipeline()

	// 存储令牌到用户的映射
	pipe.HSet(ctx, userKey, token, time.Now().Unix())
	// 设置用户令牌映射的过期时间（可选，根据需求设置）
	pipe.ExpireAt(ctx, userKey, expiry)

	// 存储令牌到用户ID的映射（用于查找令牌所属的用户）
	pipe.Set(ctx, tokenKey, userID, time.Until(expiry))

	// 执行管道
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("存储刷新令牌失败: %w", err)
	}

	return nil
}

// GetUserIDByRefreshToken 通过刷新令牌获取用户ID
func (s *TokenStore) GetUserIDByRefreshToken(ctx context.Context, token string) (string, error) {
	tokenKey := fmt.Sprintf("token:%s", token)
	userID, err := s.client.Get(ctx, tokenKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", errors.New("刷新令牌不存在或已过期")
		}
		return "", fmt.Errorf("获取用户ID失败: %w", err)
	}

	return userID, nil
}

// RemoveRefreshToken 删除刷新令牌
func (s *TokenStore) RemoveRefreshToken(ctx context.Context, token string) error {
	// 先获取令牌对应的用户ID
	tokenKey := fmt.Sprintf("token:%s", token)
	userID, err := s.client.Get(ctx, tokenKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 令牌不存在，可能已过期
			return nil
		}
		return fmt.Errorf("获取令牌对应的用户ID失败: %w", err)
	}

	// 使用管道批量执行命令
	pipe := s.client.Pipeline()

	// 从用户的令牌映射中删除该令牌
	userKey := fmt.Sprintf("user:%s:tokens", userID)
	pipe.HDel(ctx, userKey, token)

	// 删除令牌到用户ID的映射
	pipe.Del(ctx, tokenKey)

	// 执行管道
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("删除刷新令牌失败: %w", err)
	}

	return nil
}

// RemoveAllUserRefreshTokens 删除用户的所有刷新令牌
func (s *TokenStore) RemoveAllUserRefreshTokens(ctx context.Context, userID string) error {
	userKey := fmt.Sprintf("user:%s:tokens", userID)
	
	// 获取用户的所有令牌
	tokens, err := s.client.HGetAll(ctx, userKey).Result()
	if err != nil {
		return fmt.Errorf("获取用户的所有令牌失败: %w", err)
	}

	if len(tokens) == 0 {
		return nil
	}

	// 使用管道批量执行命令
	pipe := s.client.Pipeline()

	// 删除每个令牌到用户ID的映射
	for token := range tokens {
		tokenKey := fmt.Sprintf("token:%s", token)
		pipe.Del(ctx, tokenKey)
	}

	// 删除用户的令牌映射
	pipe.Del(ctx, userKey)

	// 执行管道
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("删除用户的所有刷新令牌失败: %w", err)
	}

	return nil
}
