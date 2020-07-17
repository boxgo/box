package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newLogger(level, encoding string) (*zap.Logger, *zap.AtomicLevel) {
	lv, _ := newLevelFromString(level)
	atomicLevel := zap.NewAtomicLevelAt(lv)

	cfg := &zap.Config{
		Development: false,
		Level:       atomicLevel,
		Encoding:    encoding,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	log, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return log, &atomicLevel
}

func newLevelFromString(str string) (lv zapcore.Level, err error) {
	lv = zapcore.Level(0)
	err = lv.Set(str)

	return
}

func setAtomicLevel(atomicLevel *zap.AtomicLevel, str string) error {
	if lv, err := newLevelFromString(str); err == nil {
		atomicLevel.SetLevel(lv)
		return nil
	} else {
		return err
	}
}
