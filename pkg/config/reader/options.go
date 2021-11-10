package reader

import (
	"github.com/boxgo/box/pkg/codec"
	"github.com/boxgo/box/pkg/codec/json"
	"github.com/boxgo/box/pkg/codec/toml"
	"github.com/boxgo/box/pkg/codec/xml"
	"github.com/boxgo/box/pkg/codec/yaml"
)

type Options struct {
	Encoding map[string]codec.Marshaler
}

type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		Encoding: map[string]codec.Marshaler{
			"json": json.NewMarshaler(),
			"yaml": yaml.NewMarshaler(),
			"toml": toml.NewMarshaler(),
			"xml":  xml.NewMarshaler(),
			"yml":  yaml.NewMarshaler(),
		},
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

func WithEncoder(e codec.Marshaler) Option {
	return func(o *Options) {
		if o.Encoding == nil {
			o.Encoding = make(map[string]codec.Marshaler)
		}
		o.Encoding[e.String()] = e
	}
}
