package redis

import (
	"context"
	"errors"

	"github.com/boxgo/box/pkg/logger"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
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

	client.AddHook(newMetric(cfg))
	client.AddHook(newLogger(cfg))

	if err := redisotel.InstrumentTracing(client); err != nil {
		logger.Panicf("Redis.InstrumentTracing.Error: %s", err)
	}

	if err := redisotel.InstrumentTracing(client); err != nil {
		logger.Panicf("Redis.InstrumentTracing.Error: %s", err)
	}

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

func (r *Redis) NewScript(script string) *Script {
	return newScript(r.client, script)
}
