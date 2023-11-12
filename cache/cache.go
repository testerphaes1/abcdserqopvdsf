package cache

import (
	"context"
	"time"
)

type Cache interface {
	HMSet(ctx context.Context, key string, values ...interface{}) error
	HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error)
	HSet(ctx context.Context, key string, values ...interface{}) error
	HGet(ctx context.Context, key string, field string) (string, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (interface{}, error)
	GetWithTTL(ctx context.Context, key string) (interface{}, time.Duration, error)
	Delete(ctx context.Context, key string) error
}
