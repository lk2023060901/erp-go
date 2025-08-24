package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, expiration time.Duration) error
	Del(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Keys(ctx context.Context, pattern string) ([]string, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	GetCacheStats(ctx context.Context) (map[string]interface{}, error)
}