package redis

import (
	"context"
	"strings"

	"github.com/boxgo/box/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type (
	Logger struct {
		cfg  *Config
		addr string
	}
)

func newLogger(cfg *Config) *Logger {
	return &Logger{
		cfg:  cfg,
		addr: strings.Join(cfg.Address, ","),
	}
}

func (inst *Logger) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

func (inst *Logger) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		err := next(ctx, cmd)

		if err != nil {
			inst.log(ctx, false, cmd)
		}

		return err
	}
}

func (inst *Logger) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		err := next(ctx, cmds)

		if err != nil {
			inst.log(ctx, false, cmds...)
		}

		return err
	}
}

func (inst *Logger) log(ctx context.Context, pipe bool, cmds ...redis.Cmder) {
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
			"address", inst.addr,
			"db", inst.cfg.DB,
			"err", strings.Join(errArr, ";"),
			"cmd", strings.Join(cmdArr, ";"),
		)
	}

}
