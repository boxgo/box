package logger

var (
	// Default the default logger
	Default = New()
)

func Debug(args ...interface{}) {
	Default.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	Default.Debugf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	Default.Debugw(msg, keysAndValues...)
}

func Info(args ...interface{}) {
	Default.Info(args...)
}

func Infof(template string, args ...interface{}) {
	Default.Infof(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	Default.Infow(msg, keysAndValues...)
}

func Warn(args ...interface{}) {
	Default.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	Default.Warnf(template, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	Default.Warnw(msg, keysAndValues...)
}

func Error(args ...interface{}) {
	Default.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	Default.Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	Default.Errorw(msg, keysAndValues...)
}

func DPanic(args ...interface{}) {
	Default.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	Default.DPanicf(template, args...)
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	Default.DPanicw(msg, keysAndValues...)
}

func Panic(args ...interface{}) {
	Default.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	Default.Panicf(template, args...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	Default.Panicw(msg, keysAndValues...)
}

func Fatal(args ...interface{}) {
	Default.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	Default.Fatalf(template, args...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	Default.Fatalw(msg, keysAndValues...)
}
