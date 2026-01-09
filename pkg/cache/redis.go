package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// TokenCache provides methods for caching authentication tokens
type TokenCache interface {
	SetToken(ctx context.Context, token string, userID uint, expiration time.Duration) error
	GetToken(ctx context.Context, token string) (uint, error)
	DeleteToken(ctx context.Context, token string) error
	DeleteUserTokens(ctx context.Context, userID uint) error
	Close() error
}

type redisCache struct {
	client *redis.Client
	prefix string
}

// NewRedisCache creates a new Redis cache instance for token caching
func NewRedisCache(client *redis.Client) TokenCache {
	return &redisCache{
		client: client,
		prefix: "auth:token:",
	}
}

// SetToken stores a token in Redis with the associated user ID
func (r *redisCache) SetToken(ctx context.Context, token string, userID uint, expiration time.Duration) error {
	key := r.getKey(token)
	userIDStr := fmt.Sprintf("%d", userID)
	
	// Store token -> userID mapping
	if err := r.client.Set(ctx, key, userIDStr, expiration).Err(); err != nil {
		return fmt.Errorf("failed to cache token: %w", err)
	}
	
	// Also store userID -> token mapping for easy invalidation
	userKey := r.getUserKey(userID)
	if err := r.client.SAdd(ctx, userKey, token).Err(); err != nil {
		return fmt.Errorf("failed to cache user token: %w", err)
	}
	
	// Set expiration on user key as well
	if err := r.client.Expire(ctx, userKey, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set expiration on user key: %w", err)
	}
	
	return nil
}

// GetToken retrieves the user ID associated with a token
func (r *redisCache) GetToken(ctx context.Context, token string) (uint, error) {
	key := r.getKey(token)
	userIDStr, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, fmt.Errorf("token not found in cache")
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get token from cache: %w", err)
	}
	
	var userID uint
	if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err != nil {
		return 0, fmt.Errorf("failed to parse user ID: %w", err)
	}
	
	return userID, nil
}

// DeleteToken removes a specific token from the cache
func (r *redisCache) DeleteToken(ctx context.Context, token string) error {
	key := r.getKey(token)
	
	// Get userID before deleting to clean up user key
	userIDStr, err := r.client.Get(ctx, key).Result()
	if err == nil {
		var userID uint
		if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err == nil {
			userKey := r.getUserKey(userID)
			r.client.SRem(ctx, userKey, token)
		}
	}
	
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}
	
	return nil
}

// DeleteUserTokens removes all tokens for a specific user
func (r *redisCache) DeleteUserTokens(ctx context.Context, userID uint) error {
	userKey := r.getUserKey(userID)
	
	// Get all tokens for this user
	tokens, err := r.client.SMembers(ctx, userKey).Result()
	if err != nil {
		return fmt.Errorf("failed to get user tokens: %w", err)
	}
	
	// Delete all tokens
	for _, token := range tokens {
		key := r.getKey(token)
		r.client.Del(ctx, key)
	}
	
	// Delete the user key
	if err := r.client.Del(ctx, userKey).Err(); err != nil {
		return fmt.Errorf("failed to delete user key: %w", err)
	}
	
	return nil
}

// Close closes the Redis connection
func (r *redisCache) Close() error {
	return r.client.Close()
}

// getKey returns the Redis key for a token
func (r *redisCache) getKey(token string) string {
	return r.prefix + token
}

// getUserKey returns the Redis key for a user's token set
func (r *redisCache) getUserKey(userID uint) string {
	return fmt.Sprintf("auth:user:%d:tokens", userID)
}
