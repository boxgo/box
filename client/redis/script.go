package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type (
	Script struct {
		script *redis.Script
		rdb    redis.UniversalClient
	}
)

func newScript(rdb redis.UniversalClient, src string) *Script {
	return &Script{
		rdb:    rdb,
		script: redis.NewScript(src),
	}
}

func (script Script) Run(ctx context.Context, keys []string, args ...interface{}) *redis.Cmd {
	return script.script.Run(ctx, script.rdb, keys, args...)
}
