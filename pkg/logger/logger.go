package logger

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/boxgo/box/pkg/component"
	"github.com/boxgo/box/pkg/config"
	"go.uber.org/zap"
)

type (
	// Logger logger option
	Logger struct {
		component.NoopBox
		name  string
		level *zap.AtomicLevel
		sugar *zap.SugaredLogger
		cfg   *Config
	}

	Options struct {
		name   string
		config config.SubConfigurator
	}

	OptionFunc func(*Options)
)

func WithName(name string) OptionFunc {
	return func(opts *Options) {
		opts.name = name
	}
}

func WithConfigurator(cfg config.Configurator) OptionFunc {
	return func(opts *Options) {
		opts.config = cfg
	}
}

func New(optFunc ...OptionFunc) (*Logger, error) {
	opts := &Options{}
	for _, fn := range optFunc {
		fn(opts)
	}

	if opts.config == nil {
		opts.config = config.Default
	}
	if opts.name == "" {
		opts.name = "logger"
	}

	cfg := newConfig(opts.name, opts.config)
	level := newAtomicLevelFromString(cfg.GetLevel())
	sugar := newLogger(level, cfg.GetEncoding()).Sugar()

	newLogger := &Logger{
		name:  opts.name,
		level: level,
		sugar: sugar,
		cfg:   cfg,
	}

	if err := newLogger.watch(); err != nil {
		return nil, err
	}

	return newLogger, nil
}

// Name logger config name
func (logger *Logger) Name() string {
	return logger.name
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
	if w, err := logger.cfg.Watch(logger.cfg.Level); err != nil {
		return err
	} else {
		go func() {
			for {
				time.Sleep(time.Second)

				v, _ := w.Next()
				newLv := v.String("info")
				oldLv := logger.level.String()

				if err := setAtomicLevel(logger.level, newLv); err != nil {
					log.Printf("logger.setAtomicLevel.error %s->%s\n", oldLv, newLv)
				} else {
					log.Printf("logger.setAtomicLevel.success %s->%s\n", oldLv, newLv)
				}
			}
		}()
	}

	return nil
}

func (logger *Logger) trace(ctx context.Context) *zap.SugaredLogger {
	var uid, requestID, spanId, bizId string

	traceUid := logger.cfg.GetTraceUid()
	traceRequestId := logger.cfg.GetTraceReqId()
	traceSpanId := logger.cfg.GetTraceSpanId()
	traceBizId := logger.cfg.GetTraceBizId()

	if uidStr, ok := ctx.Value(traceUid).(string); ok {
		uid = uidStr
	}
	if requestIDStr, ok := ctx.Value(traceRequestId).(string); ok {
		requestID = requestIDStr
	}
	if spanIdStr, ok := ctx.Value(traceSpanId).(string); ok {
		spanId = spanIdStr
	}
	if bizIdStr, ok := ctx.Value(traceBizId).(string); ok {
		bizId = bizIdStr
	}

	return logger.sugar.Named(fmt.Sprintf("[%s][%s][%s][%s]", uid, requestID, spanId, bizId))
}
