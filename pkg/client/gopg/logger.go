package gopg

import (
	"context"

	"github.com/boxgo/box/pkg/logger"
	"github.com/go-pg/pg/v10"
)

type (
	internalLogger struct{}
)

func init() {
	pg.SetLogger(&internalLogger{})
}

func (log internalLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	logger.Trace(ctx).Infof(format, v...)
}
