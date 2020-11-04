package logger

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

var (
	// defaultLogger the default logger.
	defaultLogger *Logger
)

func init() {
	if logger, err := StdConfig("default").Build(); err != nil {
		panic(fmt.Errorf("logger init error: %w", err))
	} else {
		defaultLogger = logger
	}
}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	defaultLogger.Debugf(template, args...)
}

// Debugw logs a message with some additional context. The variadic key-value pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
// 	s.With(keysAndValues).Debug(msg)
func Debugw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Debugw(msg, keysAndValues...)
}

// Info uses fmt.Sprint to construct and log a message.
func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	defaultLogger.Infof(template, args...)
}

// Infow logs a message with some additional context. The variadic key-value pairs are treated as they are in With.
func Infow(msg string, keysAndValues ...interface{}) {
	defaultLogger.Infow(msg, keysAndValues...)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	defaultLogger.Warnf(template, args...)
}

// Warnw logs a message with some additional context. The variadic key-value pairs are treated as they are in With.
func Warnw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Warnw(msg, keysAndValues...)
}

// Error uses fmt.Sprint to construct and log a message.
func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	defaultLogger.Errorf(template, args...)
}

// Errorw logs a message with some additional context. The variadic key-value pairs are treated as they are in With.
func Errorw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Errorw(msg, keysAndValues...)
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the logger then panics. (See DPanicLevel for details.)
func DPanic(args ...interface{}) {
	defaultLogger.DPanic(args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the logger then panics. (See DPanicLevel for details.)
func DPanicf(template string, args ...interface{}) {
	defaultLogger.DPanicf(template, args...)
}

// DPanicw logs a message with some additional context. In development, the logger then panics. (See DPanicLevel for details.) The variadic key-value pairs are treated as they are in With.
func DPanicw(msg string, keysAndValues ...interface{}) {
	defaultLogger.DPanicw(msg, keysAndValues...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(args ...interface{}) {
	defaultLogger.Panic(args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(template string, args ...interface{}) {
	defaultLogger.Panicf(template, args...)
}

// Panicw logs a message with some additional context, then panics. The variadic key-value pairs are treated as they are in With.
func Panicw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Panicw(msg, keysAndValues...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(template string, args ...interface{}) {
	defaultLogger.Fatalf(template, args...)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The variadic key-value pairs are treated as they are in With.
func Fatalw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Fatalw(msg, keysAndValues...)
}

// Trace logs a message with trace prefix and return *zap.SugaredLogger.
func Trace(ctx context.Context) *zap.SugaredLogger {
	return defaultLogger.Trace(ctx)
}

// TraceRaw logs a message with trace prefix and return *zap.Logger.
func TraceRaw(ctx context.Context) *zap.Logger {
	return defaultLogger.TraceRaw(ctx)
}

// Named adds a sub-scope to the logger's name. See Logger.Named for details.
func Named(name string) *zap.SugaredLogger {
	return defaultLogger.Named(name)
}

// Desugar unwraps a SugaredLogger, exposing the original Logger. Desugaring is quite inexpensive, so it's reasonable for a single application to use both Loggers and SugaredLoggers, converting between them on the boundaries of performance-sensitive code.
func Desugar() *zap.Logger {
	return defaultLogger.Desugar()
}
