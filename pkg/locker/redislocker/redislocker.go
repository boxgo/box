package redislocker

import (
	"context"
	"time"

	"github.com/boxgo/box/pkg/client/redis"
	"github.com/boxgo/box/pkg/locker"
)

type (
	// Locker redis lock
	Locker struct {
		cfg    *Config
		client *redis.Redis
	}
)

func newLocker(cfg *Config) locker.MutexLocker {
	lock := &Locker{
		cfg:    cfg,
		client: redis.StdConfig(cfg.Config).Build(),
	}

	return lock
}

// Lock key
func (l *Locker) Lock(key string, duration time.Duration) (bool, error) {
	return l.client.Client().SetNX(context.TODO(), locker.UnifiedKey(key), time.Now().Unix(), duration).Result()
}

// IsLocked return is is locked
func (l *Locker) IsLocked(key string) (bool, error) {
	result, err := l.client.Client().Exists(context.TODO(), locker.UnifiedKey(key)).Result()

	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return result == 1, nil
}

// UnLock unlock key
func (l *Locker) UnLock(key string) error {
	_, err := l.client.Client().Del(context.TODO(), locker.UnifiedKey(key)).Result()

	return err
}
