package redislocker

import (
	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/locker"
	"github.com/boxgo/box/pkg/logger"
)

type (
	Config struct {
		path   string
		Config string `config:"config" desc:"redis config path. eg: 'default' means use 'redis.default' config"`
	}
)

func StdConfig(key string) *Config {
	cfg := DefaultConfig(key)

	if err := config.Scan(cfg); err != nil {
		logger.Panicf("redis locker build error: %s", err)
	}

	return cfg
}

func DefaultConfig(key string) *Config {
	return &Config{
		path:   "locker." + key,
		Config: key,
	}
}

func (c *Config) Path() string {
	return c.path
}

func (c *Config) Build() locker.MutexLocker {
	return newLocker(c)
}
