package logger

import "github.com/boxgo/box/pkg/config"

type (
	Options struct {
		name   string
		config config.SubConfigurator
	}

	OptionFunc func(*Options)
)

func WithName(name string) OptionFunc {
	return func(opts *Options) {
		opts.name = name
	}
}

func WithConfigurator(cfg config.Configurator) OptionFunc {
	return func(opts *Options) {
		opts.config = cfg
	}
}
