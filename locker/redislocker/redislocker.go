package redislocker

import (
	"context"
	"time"

	redis2 "github.com/boxgo/box/v2/client/redis"
	"github.com/boxgo/box/v2/locker"
)

type (
	// Locker redis lock
	Locker struct {
		cfg    *Config
		client *redis2.Redis
	}
)

func newLocker(cfg *Config) locker.MutexLocker {
	lock := &Locker{
		cfg:    cfg,
		client: redis2.StdConfig(cfg.Config).Build(),
	}

	return lock
}

// Lock key
func (l *Locker) Lock(ctx context.Context, key string, duration time.Duration) (bool, error) {
	return l.client.Client().SetNX(ctx, l.cacheKey(key), time.Now().Unix(), duration).Result()
}

// IsLocked return is is locked
func (l *Locker) IsLocked(ctx context.Context, key string) (bool, error) {
	result, err := l.client.Client().Exists(ctx, l.cacheKey(key)).Result()

	if err == redis2.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return result == 1, nil
}

// UnLock unlock key
func (l *Locker) UnLock(ctx context.Context, key string) error {
	_, err := l.client.Client().Del(ctx, l.cacheKey(key)).Result()

	return err
}

func (l *Locker) cacheKey(key string) string {
	cacheKey := l.cfg.Prefix
	if cacheKey == "" {
		cacheKey = locker.UnifiedKey(key)
	}

	return cacheKey
}
