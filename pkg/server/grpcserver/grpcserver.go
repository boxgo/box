package grpcserver

import (
	"context"
	"net"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/dummybox"
	"google.golang.org/grpc"
)

type (
	Server struct {
		dummybox.DummyBox
		cfg           *Config
		lis           net.Listener
		server        *grpc.Server
		serverOptions []grpc.ServerOption
	}
)

func NewServer(opts ...grpc.ServerOption) *Server {
	return &Server{
		server:        nil,
		cfg:           nil,
		serverOptions: opts,
	}
}

func (s *Server) Name() string {
	return "grpcserver"
}

func (s *Server) Init(cfg config.SubConfigurator) (err error) {
	s.cfg = newConfig(s.Name(), cfg)
	s.lis, err = net.Listen(s.cfg.Network(), s.cfg.Address())
	s.server = grpc.NewServer(s.serverOptions...)
	s.cfg.Mount(s.cfg.Fields()...)

	return err
}

func (s *Server) Serve(ctx context.Context) error {
	err := s.server.Serve(s.lis)
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
