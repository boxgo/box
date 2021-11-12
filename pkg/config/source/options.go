package source

import (
	"context"

	"github.com/boxgo/box/pkg/codec"
	"github.com/boxgo/box/pkg/codec/json"
)

type Options struct {
	// Encoder
	Encoder codec.Marshaler

	// for alternative data
	Context context.Context
}

type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		Encoder: json.NewMarshaler(),
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// WithEncoder sets the source encoder
func WithEncoder(e codec.Marshaler) Option {
	return func(o *Options) {
		o.Encoder = e
	}
}
