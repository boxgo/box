package gormx

import (
	"context"
	"time"

	"github.com/boxgo/box/v2/logger"
	gormlog "gorm.io/gorm/logger"
)

type (
	Logger struct{}
)

func (l *Logger) LogMode(gormlog.LogLevel) gormlog.Interface {
	return l
}

func (*Logger) Info(ctx context.Context, tmpl string, args ...interface{}) {
	logger.Trace(ctx).Infof(tmpl, args...)
}

func (*Logger) Warn(ctx context.Context, tmpl string, args ...interface{}) {
	logger.Trace(ctx).Warnf(tmpl, args...)
}

func (*Logger) Error(ctx context.Context, tmpl string, args ...interface{}) {
	logger.Trace(ctx).Errorf(tmpl, args...)
}

func (*Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, n := fc()
	elapsed := time.Since(begin)

	if err != nil {
		logger.Trace(ctx).Warnw("GormExecuteError", "sql", sql, "elapsed", elapsed, "error", err)
	} else {
		logger.Trace(ctx).Debugw("GormExecuteSuccess", "sql", sql, "elapsed", elapsed, "effect", n)
	}
}
