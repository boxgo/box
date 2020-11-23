package grpcserver

import (
	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger"
	"google.golang.org/grpc"
)

type (
	Config struct {
		path                    string
		serverOptions           []grpc.ServerOption
		unaryServerInterceptor  []grpc.UnaryServerInterceptor
		streamServerInterceptor []grpc.StreamServerInterceptor
		Network                 string `config:"network" desc:"The network must be \"tcp\", \"tcp4\", \"tcp6\", \"unix\" or \"unixpacket\""`
		Addr                    string `config:"addr" desc:"format: host:port"`
	}

	OptionFunc func(*Config)
)

func WithServerOption(opt ...grpc.ServerOption) OptionFunc {
	return func(c *Config) {
		c.serverOptions = append(c.serverOptions, opt...)
	}
}

func WithUnaryServerInterceptor(interceptor ...grpc.UnaryServerInterceptor) OptionFunc {
	return func(c *Config) {
		c.unaryServerInterceptor = append(c.unaryServerInterceptor, interceptor...)
	}
}

func WithStreamServerInterceptor(interceptor ...grpc.StreamServerInterceptor) OptionFunc {
	return func(c *Config) {
		c.streamServerInterceptor = append(c.streamServerInterceptor, interceptor...)
	}
}

func StdConfig(key string, optionFunc ...OptionFunc) *Config {
	cfg := DefaultConfig(key)

	for _, fn := range optionFunc {
		fn(cfg)
	}

	if err := config.Scan(cfg); err != nil {
		logger.Panicf("gRPC server build error: %s", err)
	}

	return cfg
}

func DefaultConfig(key string) *Config {
	return &Config{
		path:    "grpc_server." + key,
		Network: "tcp4",
		Addr:    ":9001",
	}
}

func (c *Config) Path() string {
	return c.path
}

func (c *Config) Build() *Server {
	return newGRpcServer(c)
}
