package rediscache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/boxgo/box/pkg/cache"
	"github.com/boxgo/box/pkg/client/redis"
)

type (
	// Locker redis lock
	redisCache struct {
		cfg    *Config
		client *redis.Redis
	}
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

	return l.client.Client().Set(ctx, cache.UnifiedKey(key), data, duration).Err()
}

// UnLock unlock key
func (l *redisCache) Get(ctx context.Context, key string, val interface{}) error {
	data, err := l.client.Client().Get(ctx, cache.UnifiedKey(key)).Bytes()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return nil
	}

	return json.Unmarshal(data, val)
}
