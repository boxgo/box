package grpcserver

import (
	"context"
	"net"

	"github.com/boxgo/box/pkg/config"
	"google.golang.org/grpc"
)

type (
	Server struct {
		name   string
		cfg    *Config
		server *grpc.Server
	}

	Options struct {
		name          string
		cfg           config.SubConfigurator
		serverOptions []grpc.ServerOption
	}

	OptionFunc func(*Options)
)

func WithName(name string) OptionFunc {
	return func(opts *Options) {
		opts.name = name
	}
}

func WithConfigurator(cfg config.SubConfigurator) OptionFunc {
	return func(opts *Options) {
		opts.cfg = cfg
	}
}

func WithServerOption(serverOptions ...grpc.ServerOption) OptionFunc {
	return func(opts *Options) {
		opts.serverOptions = append(opts.serverOptions, serverOptions...)
	}
}

func New(optFunc ...OptionFunc) *Server {
	opts := &Options{}
	for _, fn := range optFunc {
		fn(opts)
	}

	if opts.name == "" {
		opts.name = "grpc.server"
	}
	if opts.cfg == nil {
		opts.cfg = config.Default
	}

	return &Server{
		name:   opts.name,
		cfg:    newConfig(opts.name, opts.cfg),
		server: grpc.NewServer(opts.serverOptions...),
	}
}

func (s *Server) Name() string {
	return s.name
}

func (s *Server) Serve(ctx context.Context) error {
	lis, err := net.Listen(s.cfg.Network(), s.cfg.Address())
	if err != nil {
		return err
	}

	err = s.server.Serve(lis)
	if err != grpc.ErrServerStopped {
		return nil
	}

	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.server.GracefulStop()

	return nil
}

func (s *Server) RawServer() *grpc.Server {
	return s.server
}
