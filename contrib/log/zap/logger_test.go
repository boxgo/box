package zap_test

import (
	"context"
	"os"
	"testing"

	"github.com/boxgo/box/v2/contrib/log/zap/v2"
)

var (
	ctx = context.TODO()
	m   = map[string]string{"password": "abc", "b": "bcd"}
)

func init() {
	ctx = context.WithValue(ctx, "traceUID", "traceUID")
	ctx = context.WithValue(ctx, "traceRequestID", "traceRequestID")
	ctx = context.WithValue(ctx, "traceSpanID", "traceSpanID")
	ctx = context.WithValue(ctx, "traceBizID", "traceBizID")

}

func Test_Box_Infow(t *testing.T) {
	hostname, _ := os.Hostname()

	log, err := zap.NewLogger(&zap.ProductionConfig,
		zap.WithTraceFields([]string{"traceUID", "traceRequestID", "traceSpanID", "traceBizID"}),
	)
	if err != nil {
		t.Fatal(err)
	}

	log = log.With(
		"service.id", hostname,
		"service.name", "zap_test",
		"service.version", "v1.0.0",
		"service.namespace", "wechat")

	log.Trace(ctx).Infow("123", "key", "value", "map", m)
	log.Trace(ctx).Infow("123", "password", "123456")
}
