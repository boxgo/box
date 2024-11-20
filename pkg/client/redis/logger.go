package redis

import (
	"context"
	"strings"
	"time"

	"github.com/boxgo/box/pkg/logger"
	"github.com/go-redis/redis/v8"
)

type (
	Logger struct {
		cfg *Config
	}
)

func (inst *Logger) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (inst *Logger) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	return nil
}

func (inst *Logger) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (inst *Logger) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}

func (inst *Logger) log(ctx context.Context, pipe bool, elapsed time.Duration, cmds ...redis.Cmder) {
	var (
		cmdArr = make([]string, len(cmds))
		errArr = make([]string, len(cmds))
	)

	for idx, cmd := range cmds {
		cmdArr[idx] = cmd.Name()

		if err := cmd.Err(); err != nil && err != redis.Nil {
			errArr[idx] = cmd.Err().Error()
		}
	}

	if len(errArr) > 0 {
		logger.Trace(ctx).Errorw("Redis.Error",
			"address", strings.Join(inst.cfg.Address, ","),
			"db", inst.cfg.DB,
			"err", strings.Join(errArr, ";"),
			"cmd", strings.Join(cmdArr, ";"),
		)
	}

}
