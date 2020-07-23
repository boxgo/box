package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newLogger(level *zap.AtomicLevel, encoding string) *zap.Logger {
	cfg := &zap.Config{
		Development: false,
		Level:       *level,
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

	return log
}

func newAtomicLevelFromString(level string) *zap.AtomicLevel {
	var atomicLevel zap.AtomicLevel
	if lv, err := newLevelFromString(level); err != nil {
		atomicLevel = zap.NewAtomicLevel()
	} else {
		atomicLevel = zap.NewAtomicLevelAt(lv)
	}

	return &atomicLevel
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
