package grpcserver

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type (
	Server struct {
		cfg    *Config
		server *grpc.Server
	}
)

func newGRpcServer(cfg *Config) *Server {
	serverOpts := append(
		cfg.serverOptions,
		grpc.UnaryInterceptor(ChainUnaryServer(cfg.unaryServerInterceptor...)),
		grpc.StreamInterceptor(ChainStreamServer(cfg.streamServerInterceptor...)),
	)

	server := grpc.NewServer(serverOpts...)

	if cfg.Reflection {
		reflection.Register(server)
	}

	return &Server{
		cfg:    cfg,
		server: server,
	}
}

func (server *Server) Name() string {
	return "gRPC-server"
}

func (server *Server) Serve(ctx context.Context) error {
	lis, err := net.Listen(server.cfg.Network, server.cfg.Addr)
	if err != nil {
		return err
	}

	err = server.server.Serve(lis)
	if err != grpc.ErrServerStopped {
		return nil
	}

	return err
}

func (server *Server) Shutdown(ctx context.Context) error {
	server.server.GracefulStop()

	return nil
}

func (server *Server) RawServer() *grpc.Server {
	return server.server
}
