package logger

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/system"
	"github.com/boxgo/box/pkg/util/jsonutil"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// Logger logger option
	Logger struct {
		level *zap.AtomicLevel
		sugar *zap.SugaredLogger
		cfg   *Config
	}
)

func newLogger(cfg *Config) (*Logger, error) {
	var zapCfg zap.Config
	if err := jsonutil.Copy(cfg, &zapCfg); err != nil {
		return nil, err
	}

	if zapLogger, err := zapCfg.Build(zap.AddCallerSkip(2)); err == nil {
		logger := &Logger{
			level: &zapCfg.Level,
			sugar: zapLogger.Sugar(),
			cfg:   cfg,
		}

		go func() {
			if err := logger.watch(); err != nil {
				logger.Errorf("logger config watch error: %s", err)
			}
		}()

		return logger, nil
	} else {
		return nil, err
	}
}

func (logger *Logger) String() string {
	return "logger"
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.sugar.Debug(args...)
}

func (logger *Logger) Debugf(template string, args ...interface{}) {
	logger.sugar.Debugf(template, args...)
}

func (logger *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	logger.sugar.Debugw(msg, keysAndValues...)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.sugar.Info(args...)
}

func (logger *Logger) Infof(template string, args ...interface{}) {
	logger.sugar.Infof(template, args...)
}

func (logger *Logger) Infow(msg string, keysAndValues ...interface{}) {
	logger.sugar.Infow(msg, keysAndValues...)
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.sugar.Warn(args...)
}

func (logger *Logger) Warnf(template string, args ...interface{}) {
	logger.sugar.Warnf(template, args...)
}

func (logger *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	logger.sugar.Warnw(msg, keysAndValues...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.sugar.Error(args...)
}

func (logger *Logger) Errorf(template string, args ...interface{}) {
	logger.sugar.Errorf(template, args...)
}

func (logger *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	logger.sugar.Errorw(msg, keysAndValues...)
}

func (logger *Logger) DPanic(args ...interface{}) {
	logger.sugar.DPanic(args...)
}

func (logger *Logger) DPanicf(template string, args ...interface{}) {
	logger.sugar.DPanicf(template, args...)
}

func (logger *Logger) DPanicw(msg string, keysAndValues ...interface{}) {
	logger.sugar.DPanicw(msg, keysAndValues...)
}

func (logger *Logger) Panic(args ...interface{}) {
	logger.sugar.Panic(args...)
}

func (logger *Logger) Panicf(template string, args ...interface{}) {
	logger.sugar.Panicf(template, args...)
}

func (logger *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	logger.sugar.Panicw(msg, keysAndValues...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	logger.sugar.Fatal(args...)
}

func (logger *Logger) Fatalf(template string, args ...interface{}) {
	logger.sugar.Fatalf(template, args...)
}

func (logger *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	logger.sugar.Fatalw(msg, keysAndValues...)
}

// Trace logger with requestId and uid
func (logger *Logger) Trace(ctx context.Context) *zap.SugaredLogger {
	return logger.trace(ctx)
}

func (logger *Logger) TraceRaw(ctx context.Context) *zap.Logger {
	return logger.trace(ctx).Desugar()
}

func (logger *Logger) Named(name string) *zap.SugaredLogger {
	return logger.sugar.Named(name)
}

func (logger *Logger) Desugar() *zap.Logger {
	return logger.sugar.Desugar()
}

func (logger *Logger) watch() error {
	if w, err := config.Watch(logger.cfg.Path(), "level"); err != nil {
		return err
	} else {
		go func() {
			for {
				time.Sleep(time.Second)

				v, _ := w.Next()
				newLv := v.String("info")
				oldLv := logger.level.String()

				zapLevel := zapcore.InfoLevel
				if err := zapLevel.Set(newLv); err != nil {
					log.Printf("logger.setAtomicLevel.error %s->%s\n", oldLv, newLv)
				} else {
					logger.level.SetLevel(zapLevel)
					log.Printf("logger.setAtomicLevel.success %s->%s\n", oldLv, newLv)
				}
			}
		}()
	}

	return nil
}

func (logger *Logger) trace(ctx context.Context) *zap.SugaredLogger {
	var uid, requestID, spanID, bizID string

	traceUID := system.TraceUID()
	traceRequestID := system.TraceReqID()
	traceSpanID := system.TraceSpanID()
	traceBizID := system.TraceBizID()

	if uidStr, ok := ctx.Value(traceUID).(string); ok {
		uid = uidStr
	}
	if requestIDStr, ok := ctx.Value(traceRequestID).(string); ok {
		requestID = requestIDStr
	}
	if spanIDStr, ok := ctx.Value(traceSpanID).(string); ok {
		spanID = spanIDStr
	}
	if bizIDStr, ok := ctx.Value(traceBizID).(string); ok {
		bizID = bizIDStr
	}

	return logger.sugar.Named(fmt.Sprintf("[%s][%s][%s][%s]", uid, requestID, spanID, bizID))
}
