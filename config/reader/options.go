package reader

import (
	"github.com/boxgo/box/v2/codec"
)

type (
	Options struct {
		TargetEncoder  codec.Coder
		SourceEncoders map[string]codec.Coder
	}

	Option func(o *Options)
)

func WithTargetEncoder(e codec.Coder) Option {
	return func(o *Options) {
		o.TargetEncoder = e
	}
}

func WithSourceEncoder(e codec.Coder) Option {
	return func(o *Options) {
		if o.SourceEncoders == nil {
			o.SourceEncoders = make(map[string]codec.Coder)
		}
		o.SourceEncoders[e.String()] = e
	}
}
