package zap_test

import (
	"testing"

	"github.com/boxgo/box/v2/contrib/log/zap/v2"
)

func Benchmark_Zap_Infow(b *testing.B) {
	zapLogger, _ := zap.ProductionConfig.Build()
	sugar := zapLogger.Sugar()

	for i := 0; i < b.N; i++ {
		sugar.Infow("123", "key", "value", "map", m)
	}
}

func Benchmark_Box_Infow(b *testing.B) {
	log, err := zap.NewLogger(&zap.ProductionConfig)

	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		log.Infow("123", "key", "value", "map", m)
	}
}

func Benchmark_Box_Trace_Infow(b *testing.B) {
	log, err := zap.NewLogger(&zap.ProductionConfig)

	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		log.Trace(ctx).Infow("123", "key", "value", "map", m)
	}
}
