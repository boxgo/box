package logger

import "context"

type (
	Logger interface {
		SetLevel(level Level)                     // set log level
		With(keysAndValues ...interface{}) Logger // log with fields
		Trace(ctx context.Context) Logger         // log with requestId, uid, bizId and spanId
		Debug(args ...interface{})
		Debugf(template string, args ...interface{})
		Debugw(msg string, keysAndValues ...interface{})
		Info(args ...interface{})
		Infof(template string, args ...interface{})
		Infow(msg string, keysAndValues ...interface{})
		Warn(args ...interface{})
		Warnf(template string, args ...interface{})
		Warnw(msg string, keysAndValues ...interface{})
		Error(args ...interface{})
		Errorf(template string, args ...interface{})
		Errorw(msg string, keysAndValues ...interface{})
		Fatal(args ...interface{})
		Fatalf(template string, args ...interface{})
		Fatalw(msg string, keysAndValues ...interface{})
	}

	// A Level is a logging priority. Higher levels are more important.
	Level int8
)

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel
)
