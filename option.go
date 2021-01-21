package box

import (
	"github.com/boxgo/box/pkg/component"
)

type (
	// Options new box options
	Options struct {
		StartupTimeout  int
		ShutdownTimeout int
		AutoMaxProcs    *bool
		Boxes           []component.Box
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
func WithBoxes(boxes ...component.Box) Option {
	return func(ops *Options) {
		ops.Boxes = boxes
	}
}
