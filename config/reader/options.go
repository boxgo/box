package reader

import (
	"github.com/boxgo/box/v2/codec"
	"github.com/boxgo/box/v2/codec/json"
	"github.com/boxgo/box/v2/codec/toml"
	"github.com/boxgo/box/v2/codec/xml"
	"github.com/boxgo/box/v2/codec/yaml"
)

type Options struct {
	Encoding map[string]codec.Coder
}

type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		Encoding: map[string]codec.Coder{
			"json": json.NewCoder(),
			"yaml": yaml.NewCoder(),
			"toml": toml.NewCoder(),
			"xml":  xml.NewCoder(),
			"yml":  yaml.NewCoder(),
		},
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

func WithEncoder(e codec.Coder) Option {
	return func(o *Options) {
		if o.Encoding == nil {
			o.Encoding = make(map[string]codec.Coder)
		}
		o.Encoding[e.String()] = e
	}
}
