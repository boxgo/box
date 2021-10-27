package box

type (
	// Options new box options
	Options struct {
		StartupTimeout  int
		ShutdownTimeout int
		AutoMaxProcs    *bool
		Boxes           []Box
		Tags            []string
	}

	// Option setter
	Option func(ops *Options)
)

func WithStartupTimeout(timeout int) Option {
	return func(ops *Options) {
		ops.StartupTimeout = timeout
	}
}

// WithShutdownTimeout
func WithShutdownTimeout(timeout int) Option {
	return func(ops *Options) {
		ops.ShutdownTimeout = timeout
	}
}

// WithAutoMaxProcs
func WithAutoMaxProcs(autoMaxProcs bool) Option {
	return func(ops *Options) {
		ops.AutoMaxProcs = &autoMaxProcs
	}
}

// WithBoxes set boxes
func WithBoxes(boxes ...Box) Option {
	return func(ops *Options) {
		ops.Boxes = boxes
	}
}

// WithTag function
func WithTag(fn func() []string) Option {
	return func(ops *Options) {
		ops.Tags = fn()
	}
}
