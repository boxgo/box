package redis

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
)

type (
	Redis struct {
		cfg    *Config
		client redis.UniversalClient
	}
)

func newRedis(cfg *Config) *Redis {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		MasterName:   cfg.MasterName,
		Addrs:        cfg.Address,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConnCnt,
	})

	client.AddHook(&Metric{cfg: cfg})

	r := &Redis{
		cfg:    cfg,
		client: client,
	}

	return r
}

func (r *Redis) Name() string {
	return "redis"
}

func (r *Redis) Serve(ctx context.Context) error {
	if r.client != nil {
		return r.client.Ping(ctx).Err()
	}

	return errors.New("redis client not init")
}

func (r *Redis) Shutdown(ctx context.Context) error {
	if r.client != nil {
		return r.client.Close()
	}

	return errors.New("redis client not init")
}

func (r *Redis) Client() redis.UniversalClient {
	return r.client
}
