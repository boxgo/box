package rediscache

import (
	"github.com/boxgo/box/pkg/cache"
	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger"
)

type (
	Config struct {
		path         string
		Prefix       string `config:"prefix" desc:"cache key prefix. if empty, auto prefix with format: ${serviceName}.cache.${key}."`
		Marshaller   string `config:"marshaller" desc:"support json only"`
		Config       string `config:"config" desc:"redis config path. eg: 'default' means use 'redis.default' config"`
		BigCacheSize int    `config:"bigCacheSize" desc:"big cache size, unit: byte"`
	}
)

func StdConfig(key string) *Config {
	cfg := DefaultConfig(key)

	if err := config.Scan(cfg); err != nil {
		logger.Panicf("redis cache load config error: %s", err)
	} else {
		logger.Debugf("redis cache load config %s", err)
	}

	return cfg
}

func DefaultConfig(key string) *Config {
	return &Config{
		path:         "cache." + key,
		Marshaller:   "json",
		Config:       key,
		BigCacheSize: 1000,
	}
}

func (c *Config) Path() string {
	return c.path
}

func (c *Config) Build() cache.Cache {
	return newCache(c)
}
