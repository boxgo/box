package logger_test

import (
	"testing"

	logger2 "github.com/boxgo/box/v2/logger"
	"github.com/boxgo/box/v2/util/jsonutil"
	"go.uber.org/zap"
)

func Benchmark_Zap_Infow(b *testing.B) {
	cfg := logger2.DefaultConfig("default")
	zapCfg := &zap.Config{}
	jsonutil.Copy(cfg, zapCfg)
	zapLogger, _ := zapCfg.Build()
	sugar := zapLogger.Sugar()

	for i := 0; i < b.N; i++ {
		sugar.Infow("123", "key", "value", "map", m)
	}
}

func Benchmark_Box_Infow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		logger2.Infow("123", "key", "value", "map", m)
	}
}

func Benchmark_Box_Trace_Infow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		logger2.Trace(ctx).Infow("123", "key", "value", "map", m)
	}
}
