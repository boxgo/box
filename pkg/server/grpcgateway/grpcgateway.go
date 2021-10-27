package grpcgateway

import (
	"context"
	"net/http"
	"time"

	"github.com/boxgo/box/pkg/logger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type (
	Server struct {
		cfg    *Config
		server *http.Server
		client *grpc.ClientConn
	}
)

func newGRpcGateway(cfg *Config) *Server {
	var (
		mux    = runtime.NewServeMux()
		server = &http.Server{
			Addr:    cfg.Addr,
			Handler: mux,
		}
		ctx, cancel = context.WithTimeout(context.TODO(), time.Second*10)
		conn, err   = grpc.DialContext(
			ctx,
			cfg.Target,
			grpc.WithInsecure(),
		)
	)

	defer cancel()

	if err != nil {
		logger.Panicf("%s", err)
	}

	if err = cfg.wrap(context.TODO(), mux, conn); err != nil {
		logger.Panicf("%s", err)
	}

	for _, handler := range cfg.handlers {
		if err = handler(context.TODO(), mux, conn); err != nil {
			logger.Panicf("%s", err)
		}
	}

	return &Server{
		cfg:    cfg,
		server: server,
		client: conn,
	}
}

func (server *Server) Name() string {
	return "gRPC-gateway"
}

func (server *Server) Serve(ctx context.Context) error {
	if err := server.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (server *Server) Shutdown(ctx context.Context) error {
	if err := server.client.Close(); err != nil {
		logger.Errorf("gRPC-Gateway shutdown client error: %s", err)
	}

	return server.server.Shutdown(ctx)
}
