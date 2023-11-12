package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type redisCache struct {
	client *redis.Client
}

// Set value in cache
func (r *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *redisCache) HMSet(ctx context.Context, key string, values ...interface{}) error {
	return r.client.HMSet(ctx, key, values).Err()
}

func (r *redisCache) HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
	return r.client.HMGet(ctx, key, fields...).Result()
}

func (r *redisCache) HSet(ctx context.Context, key string, values ...interface{}) error {
	return r.client.HSet(ctx, key, values).Err()
}

func (r *redisCache) HGet(ctx context.Context, key string, field string) (string, error) {
	return r.client.HGet(ctx, key, field).Result()
}

// Get value from cache
func (r *redisCache) Get(ctx context.Context, key string) (interface{}, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *redisCache) GetWithTTL(ctx context.Context, key string) (interface{}, time.Duration, error) {
	s := r.client.Get(ctx, key)
	if err := s.Err(); err != nil {
		return nil, 0, err
	}
	ttl := r.client.TTL(ctx, key)
	if err := ttl.Err(); err != nil {
		return nil, 0, err
	}
	return s.Val(), ttl.Val(), nil
}

// Delete a key from cache
func (r *redisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// NewRedisCache create new redis cache
func NewRedisCache(redisClient *redis.Client) Cache {
	return &redisCache{redisClient}
}
