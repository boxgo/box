package grpcserver

import (
	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger"
	"google.golang.org/grpc"
)

type (
	Config struct {
		path                    string
		wrap                    Wrap
		serverOptions           []grpc.ServerOption
		unaryServerInterceptor  []grpc.UnaryServerInterceptor
		streamServerInterceptor []grpc.StreamServerInterceptor
		Network                 string `config:"network" desc:"The network must be \"tcp\", \"tcp4\", \"tcp6\", \"unix\" or \"unixpacket\""`
		Addr                    string `config:"addr" desc:"format: host:port"`
		Reflection              bool   `config:"reflection" desc:"Enable server reflection service"`
	}

	Wrap func(*grpc.Server)

	OptionFunc func(*Config)
)

func WithWrap(wrap Wrap) OptionFunc {
	return func(c *Config) {
		c.wrap = wrap
	}
}

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
		logger.Panicf("gRPC server load config error: %s", err)
	} else {
		logger.Debugw("gRPC server load config", "config", cfg)
	}

	return cfg
}

func DefaultConfig(key string) *Config {
	return &Config{
		path:       "grpcServer." + key,
		wrap:       func(server *grpc.Server) {},
		Network:    "tcp4",
		Addr:       ":9001",
		Reflection: false,
	}
}

func (c *Config) Path() string {
	return c.path
}

func (c *Config) Build() *Server {
	return newGRpcServer(c)
}
