package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisRepository provides methods to interact with Redis
type RedisRepository struct {
	client *redis.Client
}

// NewRedisRepository creates a new instance
func NewRedisRepository(addr, password string, db int) *RedisRepository {
	client := redis.NewClient(&redis.Options{
		Addr:     addr, // e.g., "redis:6379"
		Password: password,
		DB:       db,
	})

	return &RedisRepository{
		client: client,
	}
}

// GetAssignment retrieves the assigned variant (A/B) for a user in a test
func (r *RedisRepository) GetAssignment(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		// Key not found
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

// SetAssignment saves the assigned variant (A/B) for a user in a test with TTL
func (r *RedisRepository) SetAssignment(ctx context.Context, key, variant string, ttl time.Duration) error {
	err := r.client.Set(ctx, key, variant, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set assignment in redis: %w", err)
	}
	return nil
}

// Close closes the connection to Redis
func (r *RedisRepository) Close() error {
	return r.client.Close()
}
