package logger

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger/core"
	"github.com/boxgo/box/pkg/trace"
	"github.com/boxgo/box/pkg/util/jsonutil"
	"github.com/boxgo/box/pkg/util/strutil"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// Logger logger option
	Logger struct {
		cfg   *Config
		level *zap.AtomicLevel
		sugar *zap.SugaredLogger
	}
)

const (
	traceSplitterL = '['
	traceSplitterR = ']'
)

func newLogger(cfg *Config) (*Logger, error) {
	var (
		zapCfg           zap.Config
		err              error
		encoder          zapcore.Encoder
		outSink, errSink zapcore.WriteSyncer
	)
	if err = jsonutil.Copy(cfg, &zapCfg); err != nil {
		return nil, err
	}

	if zapCfg.Encoding == "json" {
		encoder = zapcore.NewJSONEncoder(zapCfg.EncoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(zapCfg.EncoderConfig)
	}

	if outSink, errSink, err = openSinks(cfg); err != nil {
		return nil, err
	}

	rule := cfg.MaskRules
	if !cfg.Mask {
		rule = nil
	}

	zapLogger := zap.New(
		core.NewMaskCore(rule, cfg.Level, encoder, outSink),
		buildOptions(&zapCfg, errSink)...,
	)

	logger := &Logger{
		cfg:   cfg,
		level: &zapCfg.Level,
		sugar: zapLogger.Sugar(),
	}

	if err = logger.watch(); err != nil {
		logger.Errorf("logger config watch error: %s", err)
	}

	return logger, nil
}

func openSinks(cfg *Config) (zapcore.WriteSyncer, zapcore.WriteSyncer, error) {
	sink, closeOut, err := zap.Open(cfg.OutputPaths...)
	if err != nil {
		return nil, nil, err
	}
	errSink, _, err := zap.Open(cfg.ErrorOutputPaths...)
	if err != nil {
		closeOut()
		return nil, nil, err
	}
	return sink, errSink, nil
}

func buildOptions(cfg *zap.Config, errSink zapcore.WriteSyncer) []zap.Option {
	opts := []zap.Option{zap.ErrorOutput(errSink)}

	if cfg.Development {
		opts = append(opts, zap.Development())
	}

	if !cfg.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}

	stackLevel := zap.ErrorLevel
	if cfg.Development {
		stackLevel = zap.WarnLevel
	}
	if !cfg.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}

	if scfg := cfg.Sampling; scfg != nil {
		opts = append(opts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			var samplerOpts []zapcore.SamplerOption
			if scfg.Hook != nil {
				samplerOpts = append(samplerOpts, zapcore.SamplerHook(scfg.Hook))
			}
			return zapcore.NewSamplerWithOptions(
				core,
				time.Second,
				cfg.Sampling.Initial,
				cfg.Sampling.Thereafter,
				samplerOpts...,
			)
		}))
	}

	if len(cfg.InitialFields) > 0 {
		fs := make([]zap.Field, 0, len(cfg.InitialFields))
		keys := make([]string, 0, len(cfg.InitialFields))
		for k := range cfg.InitialFields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fs = append(fs, zap.Any(k, cfg.InitialFields[k]))
		}
		opts = append(opts, zap.Fields(fs...))
	}

	return opts
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

func (logger *Logger) With(args ...interface{}) *Logger {
	return &Logger{
		level: logger.level,
		sugar: logger.sugar.With(args...),
		cfg:   logger.cfg,
	}
}

// Trace logger with requestId and uid
func (logger *Logger) Trace(ctx context.Context) *Logger {
	var (
		uid, requestID, spanID, bizID string
		sugar                         *zap.SugaredLogger
	)

	if uidStr, ok := ctx.Value(trace.ID()).(string); ok {
		uid = uidStr
	}
	if requestIDStr, ok := ctx.Value(trace.ReqID()).(string); ok {
		requestID = requestIDStr
	}
	if spanIDStr, ok := ctx.Value(trace.SpanID()).(string); ok {
		spanID = spanIDStr
	}
	if bizIDStr, ok := ctx.Value(trace.BizID()).(string); ok {
		bizID = bizIDStr
	}

	switch logger.cfg.Encoding {
	case "console":
		prefixBuilder := strings.Builder{}
		prefixBuilder.WriteByte(traceSplitterL)
		prefixBuilder.Write(strutil.String2Bytes(uid))
		prefixBuilder.WriteByte(traceSplitterR)
		prefixBuilder.WriteByte(traceSplitterL)
		prefixBuilder.Write(strutil.String2Bytes(requestID))
		prefixBuilder.WriteByte(traceSplitterR)
		prefixBuilder.WriteByte(traceSplitterL)
		prefixBuilder.Write(strutil.String2Bytes(spanID))
		prefixBuilder.WriteByte(traceSplitterR)
		prefixBuilder.WriteByte(traceSplitterL)
		prefixBuilder.Write(strutil.String2Bytes(bizID))
		prefixBuilder.WriteByte(traceSplitterR)
		sugar = logger.sugar.Named(prefixBuilder.String())
	default:
		sugar = logger.sugar.With(
			"uid", uid,
			"reqId", requestID,
			"bizId", bizID,
			"spanId", spanID,
		)
	}

	return &Logger{
		level: logger.level,
		sugar: sugar,
		cfg:   logger.cfg,
	}
}

func (logger *Logger) Internal() interface{} {
	return logger.sugar.Desugar()
}

func (logger *Logger) watch() error {
	w, err := config.Watch(logger.cfg.Path(), "level")
	if err != nil {
		return err
	}

	go func() {
		for {
			time.Sleep(logger.cfg.WatchInterval)

			v, _ := w.Next()
			newLv := v.String("info")
			oldLv := logger.level.String()

			zapLevel := zapcore.InfoLevel
			if err := zapLevel.Set(newLv); err != nil {
				Infof("Logger.UpdateLevel.Error: %s -> %s\n", oldLv, newLv)
			} else {
				logger.level.SetLevel(zapLevel)
				Infof("Logger.UpdateLevel.Success: %s -> %s\n", oldLv, newLv)
			}
		}
	}()

	return nil
}
