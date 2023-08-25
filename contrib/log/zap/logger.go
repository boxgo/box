package zap

import (
	"context"
	"sort"
	"time"

	"github.com/boxgo/box/v2/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// Logger logger option
	Logger struct {
		opt   Options
		cfg   *zap.Config
		level *zap.AtomicLevel
		sugar *zap.SugaredLogger
	}

	Config = zap.Config
)

var (
	ProductionConfig = zap.NewProductionConfig()
)

func NewLogger(zapCfg *zap.Config, opts ...Option) (logger.Logger, error) {
	var (
		err              error
		opt              Options
		encoder          zapcore.Encoder
		outSink, errSink zapcore.WriteSyncer
	)

	for _, options := range opts {
		options(&opt)
	}

	if opt.maskRules == nil {
		opt.maskRules = DefaultMaskRules
	}

	if zapCfg.Encoding == "json" {
		encoder = zapcore.NewJSONEncoder(zapCfg.EncoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(zapCfg.EncoderConfig)
	}

	if outSink, errSink, err = openSinks(zapCfg); err != nil {
		return nil, err
	}

	zapLogger := zap.New(
		NewMaskCore(opt.maskRules, zapCfg.Level, encoder, outSink),
		buildOptions(zapCfg, errSink)...,
	)

	return &Logger{
		opt:   opt,
		cfg:   zapCfg,
		level: &zapCfg.Level,
		sugar: zapLogger.Sugar(),
	}, nil
}

func openSinks(cfg *zap.Config) (zapcore.WriteSyncer, zapcore.WriteSyncer, error) {
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

func (logger Logger) SetLevel(level logger.Level) {
	logger.level.SetLevel(zapcore.Level(level))
}

func (logger *Logger) With(keysAndValues ...interface{}) logger.Logger {
	return &Logger{
		opt:   logger.opt,
		cfg:   logger.cfg,
		level: logger.level,
		sugar: logger.sugar.With(keysAndValues...),
	}
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
func (logger *Logger) Trace(ctx context.Context) logger.Logger {
	var kv = make([]interface{}, len(logger.opt.traceFields)*2)

	for idx, field := range logger.opt.traceFields {
		if val, ok := ctx.Value(field).(string); ok {
			kv[idx*2] = field
			kv[idx*2+1] = val
		}
	}

	sugar := logger.sugar.With(kv...)

	return &Logger{
		opt:   logger.opt,
		cfg:   logger.cfg,
		level: logger.level,
		sugar: sugar,
	}
}
