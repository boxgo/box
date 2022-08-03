package redis

import (
	"context"
	"encoding/json"
	"log"

	"github.com/boxgo/box/v2/config/source"
)

type (
	redisConfigKey struct{}
	prefixKey      struct{}

	redisConfig struct {
		MasterName     string   `config:"masterName" desc:"The sentinel master name. Only failover clients."`
		Address        []string `config:"address" desc:"Either a single address or a seed list of host:port addresses of cluster/sentinel nodes."`
		Password       string   `config:"password" desc:"Redis password"`
		DB             int      `config:"db" desc:"Database to be selected after connecting to the server. Only single-node and failover clients."`
		PoolSize       int      `config:"poolSize" desc:"Connection pool size"`
		MinIdleConnCnt int      `config:"minIdleConnCnt" desc:"Min idle connections."`
	}
)

func WithConfig(data []byte) []source.Option {
	v := redisConfig{}
	if err := json.Unmarshal(data, &v); err != nil {
		log.Panicf("config redis build error: %#v", err)
	}

	return []source.Option{
		WithRedisConfig(v),
	}
}

func WithRedisConfig(cfg redisConfig) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, redisConfigKey{}, cfg)
	}
}

func WithPrefix(prefix string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, prefixKey{}, prefix)
	}
}
