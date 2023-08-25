package logger

import (
	"context"
)

type (
	nopLogger struct{}
)

var (
	Nop = nopLogger{}
)

func (n nopLogger) SetLevel(level Level)                            {}
func (n nopLogger) With(keysAndValues ...interface{}) Logger        { return n }
func (n nopLogger) Trace(ctx context.Context) Logger                { return n }
func (n nopLogger) Debug(args ...interface{})                       {}
func (n nopLogger) Debugf(template string, args ...interface{})     {}
func (n nopLogger) Debugw(msg string, keysAndValues ...interface{}) {}
func (n nopLogger) Info(args ...interface{})                        {}
func (n nopLogger) Infof(template string, args ...interface{})      {}
func (n nopLogger) Infow(msg string, keysAndValues ...interface{})  {}
func (n nopLogger) Warn(args ...interface{})                        {}
func (n nopLogger) Warnf(template string, args ...interface{})      {}
func (n nopLogger) Warnw(msg string, keysAndValues ...interface{})  {}
func (n nopLogger) Error(args ...interface{})                       {}
func (n nopLogger) Errorf(template string, args ...interface{})     {}
func (n nopLogger) Errorw(msg string, keysAndValues ...interface{}) {}
func (n nopLogger) Fatal(args ...interface{})                       {}
func (n nopLogger) Fatalf(template string, args ...interface{})     {}
func (n nopLogger) Fatalw(msg string, keysAndValues ...interface{}) {}
