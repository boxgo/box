package redis

import (
	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger"
)

type (
	Config struct {
		path           string
		MasterName     string   `config:"masterName" desc:"The sentinel master name. Only failover clients."`
		Address        []string `config:"address" desc:"Either a single address or a seed list of host:port addresses of cluster/sentinel nodes."`
		Password       string   `config:"password" desc:"Redis password"`
		DB             int      `config:"db" desc:"Database to be selected after connecting to the server. Only single-node and failover clients."`
		PoolSize       int      `config:"poolSize" desc:"Connection pool size"`
		MinIdleConnCnt int      `config:"minIdleConnCnt" desc:"Min idle connections."`
	}
)

func StdConfig(key string) *Config {
	cfg := DefaultConfig(key)
	if err := config.Scan(cfg); err != nil {
		logger.Panicf("redis load config error: %s", err)
	} else {
		logger.Debugw("redis load config", "config", cfg)
	}

	return cfg
}

func DefaultConfig(key string) *Config {
	return &Config{
		path:           "redis." + key,
		MasterName:     "",
		Address:        []string{"127.0.0.1:6379"},
		Password:       "",
		DB:             0,
		PoolSize:       100,
		MinIdleConnCnt: 5,
	}
}

func (c *Config) Path() string {
	return c.path
}

func (c *Config) Build() *Redis {
	return newRedis(c)
}
