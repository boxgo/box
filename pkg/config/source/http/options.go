package http

import (
	"context"
	"encoding/json"
	"log"

	"github.com/boxgo/box/pkg/config/source"
)

type (
	httpConfig struct {
		Url           string `config:"url"`
		Config        string `config:"config" desc:"fetch config name"`
		Authorization *struct {
			Type        string `config:"type" desc:"auth type: Basic/Bearer"`
			Credentials string `config:"credentials" desc:"Basic(username:password),Bearer(jwt)"`
		} `json:"authorization"`
	}

	serviceKey    struct{}
	versionKey    struct{}
	httpConfigKey struct{}
)

func WithConfig(data []byte) []source.Option {
	v := httpConfig{}
	if err := json.Unmarshal(data, &v); err != nil {
		log.Panicf("config http build error: %#v", err)
	}

	return []source.Option{
		WithHTTPConfig(v),
	}
}

func WithService(service string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, serviceKey{}, service)
	}
}

func WithVersion(version string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, versionKey{}, version)
	}
}

func WithHTTPConfig(cfg httpConfig) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, httpConfigKey{}, cfg)
	}
}
