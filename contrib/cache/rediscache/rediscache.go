package rediscache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/boxgo/box/v2/cache"
	"github.com/boxgo/box/v2/client/redis"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type (
	// Cache redis cache
	redisCache struct {
		client  *redis.Redis
		tracer  trace.Tracer
		counter metric.Int64Counter
	}
)

func New(client *redis.Redis) (cache.Cache, error) {
	meter := otel.Meter("cache-redis")
	tracer := otel.Tracer("cache-redis")

	counter, err := meter.Int64Counter("cache_hit_total", metric.WithDescription("cache hit counter"))
	if err != nil {
		return nil, err
	}

	inst := &redisCache{
		client:  client,
		counter: counter,
		tracer:  tracer,
	}

	return inst, nil
}

// Set cache
func (inst *redisCache) Set(ctx context.Context, key string, val interface{}, duration time.Duration) error {
	ctx, span := inst.tracer.Start(ctx, "CacheSet")
	defer span.End()

	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return inst.client.Client().Set(ctx, cache.UnifiedKeyVer(key), data, duration).Err()
}

// Get cache
func (inst *redisCache) Get(ctx context.Context, key string, val interface{}) error {
	ctx, span := inst.tracer.Start(ctx, "CacheGet")
	defer span.End()

	data, err := inst.client.Client().Get(ctx, cache.UnifiedKeyVer(key)).Bytes()
	if err == redis.Nil {
		inst.counter.Add(ctx, 1, metric.WithAttributes(
			attribute.String("key", key),
			attribute.String("hit", "false"),
		))

		return cache.ErrCacheMiss
	} else if err != nil {
		return nil
	}

	inst.counter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("key", key),
		attribute.String("hit", "true"),
	))

	return json.Unmarshal(data, val)
}

func (inst *redisCache) Clear(ctx context.Context, key string) error {
	ctx, span := inst.tracer.Start(ctx, "CacheClear")
	defer span.End()

	return inst.client.Client().Del(ctx, cache.UnifiedKeyVer(key)).Err()
}
