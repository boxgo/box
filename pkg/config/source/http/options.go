package http

import (
	"context"
	"encoding/json"
	"log"

	"github.com/boxgo/box/pkg/config/source"
)

type (
	Config struct {
		Namespace     string        `config:"namespace" desc:"default is global namespace"`
		Service       string        `config:"service" desc:"default is global service"`
		Version       string        `config:"version" desc:"default is global version"`
		Url           string        `config:"url"`
		Authorization authorization `json:"authorization"`
	}

	httpConfig struct {
		Url           string
		Authorization authorization
	}

	authorization struct {
		Type        string `config:"type" desc:"auth type: Basic/Bearer"`
		Credentials string `config:"credentials" desc:"Basic(username:password),Bearer(jwt)"`
	}

	namespaceKey  struct{}
	serviceKey    struct{}
	versionKey    struct{}
	httpConfigKey struct{}
)

func WithConfig(data []byte) []source.Option {
	v := Config{}
	if err := json.Unmarshal(data, &v); err != nil {
		log.Panicf("config http build error: %#v", err)
	}

	return []source.Option{
		WithHTTPConfig(httpConfig{
			Url:           v.Url,
			Authorization: authorization{Type: v.Authorization.Type, Credentials: v.Authorization.Credentials},
		}),
		WithNamespace(v.Namespace),
		WithService(v.Service),
		WithVersion(v.Version),
	}
}

func WithNamespace(namespace string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, namespaceKey{}, namespace)
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

func WithHTTPConfig(httpCfg httpConfig) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, httpConfigKey{}, httpCfg)
	}
}
