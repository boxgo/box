package logger_test

import (
	"context"
	"testing"

	"github.com/boxgo/box/pkg/logger"
	"github.com/boxgo/box/pkg/trace"
)

var (
	traceUID       = trace.ID()
	traceRequestID = trace.ReqID()
	traceSpanID    = trace.SpanID()
	traceBizID     = trace.BizID()
	ctx            = context.TODO()
	m              = map[string]string{"password": "abc", "b": "bcd"}
)

func init() {
	ctx = context.WithValue(ctx, traceUID, "traceUID")
	ctx = context.WithValue(ctx, traceRequestID, "traceRequestID")
	ctx = context.WithValue(ctx, traceSpanID, "traceSpanID")
	ctx = context.WithValue(ctx, traceBizID, "traceBizID")

}

func Test_Box_Infow(t *testing.T) {
	logger.Infow("123", "key", "value", "map", m)
}
