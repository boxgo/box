package grpcgateway

import (
	"context"

	"github.com/boxgo/box/v2/config"
	"github.com/boxgo/box/v2/logger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type (
	Config struct {
		path     string
		wrap     Wrap
		handlers []Handler
		muxOpts  []runtime.ServeMuxOption
		Target   string `config:"target" desc:"target gprc server addr"`
		Addr     string `config:"addr" desc:"format: host:port"`
	}

	Wrap       func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
	Handler    Wrap
	OptionFunc func(*Config)
)

func WithWrap(wrap Wrap) OptionFunc {
	return func(c *Config) {
		c.wrap = wrap
	}
}

func WithRegisterHandler(handlers ...Handler) OptionFunc {
	return func(c *Config) {
		c.handlers = append(c.handlers, handlers...)
	}
}

func WithMuxOption(opts ...runtime.ServeMuxOption) OptionFunc {
	return func(c *Config) {
		c.muxOpts = append(c.muxOpts, opts...)
	}
}

func StdConfig(key string, optionFunc ...OptionFunc) *Config {
	cfg := DefaultConfig(key)

	for _, fn := range optionFunc {
		fn(cfg)
	}

	if err := config.Scan(cfg); err != nil {
		logger.Panicf("gRPC gateway build error: %s", err)
	}

	return cfg
}

func DefaultConfig(key string) *Config {
	return &Config{
		path:   "grpc_gateway." + key,
		Target: ":9001",
		Addr:   ":9002",
		wrap: func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
			return nil
		},
	}
}

func (c *Config) Path() string {
	return c.path
}

func (c *Config) Build() *Server {
	return newGRpcGateway(c)
}
