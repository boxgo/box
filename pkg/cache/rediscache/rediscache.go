package rediscache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/boxgo/box/pkg/cache"
	"github.com/boxgo/box/pkg/client/redis"
	"github.com/boxgo/box/pkg/metric"
)

type (
	// Cache redis cache
	redisCache struct {
		cfg    *Config
		client *redis.Redis
	}
)

var (
	cacheHitCounter = metric.NewCounterVec(
		"cache_hit_total",
		"cache hit counter",
		[]string{"key", "hit"},
	)
)

func newCache(cfg *Config) cache.Cache {
	lock := &redisCache{
		cfg:    cfg,
		client: redis.StdConfig(cfg.Config).Build(),
	}

	return lock
}

// Set cache
func (l *redisCache) Set(ctx context.Context, key string, val interface{}, duration time.Duration) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return l.client.Client().Set(ctx, l.cacheKey(key), data, duration).Err()
}

// Get cache
func (l *redisCache) Get(ctx context.Context, key string, val interface{}) error {
	data, err := l.client.Client().Get(ctx, l.cacheKey(key)).Bytes()
	if err == redis.Nil {
		cacheHitCounter.WithLabelValues(key, "false").Inc()
		return cache.ErrCacheMiss
	} else if err != nil {
		return nil
	}

	cacheHitCounter.WithLabelValues(key, "true").Inc()
	return json.Unmarshal(data, val)
}

func (l *redisCache) Clear(ctx context.Context, key string) error {
	return l.client.Client().Del(ctx, l.cacheKey(key)).Err()
}

func (l *redisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return l.client.Client().Expire(ctx, l.cacheKey(key), expiration).Err()
}

func (l redisCache) cacheKey(key string) string {
	cacheKey := l.cfg.Prefix
	if cacheKey == "" {
		cacheKey = cache.UnifiedKey(key)
	}

	return cacheKey
}
