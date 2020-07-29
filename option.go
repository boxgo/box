package box

import (
	"github.com/boxgo/box/pkg/config"
)

type (
	// Options new box options
	Options struct {
		Silent          bool
		StartupTimeout  int
		ShutdownTimeout int
		Boxes           []Box
		Configurator    config.Configurator
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

// WithConfig set configurator
func WithConfig(cfg config.Configurator) Option {
	return func(ops *Options) {
		ops.Configurator = cfg
	}
}

// WithBoxes set boxes
func WithBoxes(boxes ...Box) Option {
	return func(ops *Options) {
		ops.Boxes = boxes
	}
}
