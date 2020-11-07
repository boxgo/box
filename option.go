package box

import (
	"github.com/boxgo/box/pkg/component"
)

type (
	// Options new box options
	Options struct {
		Silent          bool
		StartupTimeout  int
		ShutdownTimeout int
		Boxes           []component.Box
	}

	// Option setter
	Option func(ops *Options)
)

func WithSilent(silent bool) Option {
	return func(ops *Options) {
		ops.Silent = silent
	}
}

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

// WithBoxes set boxes
func WithBoxes(boxes ...component.Box) Option {
	return func(ops *Options) {
		ops.Boxes = boxes
	}
}
