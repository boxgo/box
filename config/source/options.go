package source

import (
	"context"

	"github.com/boxgo/box/v2/codec"
	"github.com/boxgo/box/v2/codec/json"
)

type (
	Options struct {
		Encoder codec.Coder     // Encoder
		Context context.Context // for alternative data
	}

	Option func(o *Options)
)

func NewOptions(opts ...Option) Options {
	options := Options{
		Encoder: json.NewCoder(),
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// WithEncoder sets the source encoder
func WithEncoder(e codec.Coder) Option {
	return func(o *Options) {
		o.Encoder = e
	}
}
