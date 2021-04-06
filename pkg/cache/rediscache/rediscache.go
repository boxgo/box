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
	// Locker redis lock
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

// Lock key
func (l *redisCache) Set(ctx context.Context, key string, val interface{}, duration time.Duration) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return l.client.Client().Set(ctx, l.cacheKey(key), data, duration).Err()
}

// UnLock unlock key
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

func (l redisCache) cacheKey(key string) string {
	cacheKey := l.cfg.Prefix
	if cacheKey == "" {
		cacheKey = cache.UnifiedKey(key)
	}

	return cacheKey
}
