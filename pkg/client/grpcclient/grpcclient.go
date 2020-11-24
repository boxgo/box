package grpcclient

import (
	"context"

	"github.com/boxgo/box/pkg/logger"
	"google.golang.org/grpc"
)

type (
	Client struct {
		cfg      *Config
		conn     *grpc.ClientConn
		dialOpts []grpc.DialOption
	}
)

func newGRpcClient(cfg *Config) *Client {
	dialOpts := append(
		cfg.dialOptions,
		grpc.WithChainUnaryInterceptor(cfg.unaryClientInterceptor...),
		grpc.WithChainStreamInterceptor(cfg.streamClientInterceptor...),
	)

	if cfg.Dial {
		if conn, err := grpc.DialContext(context.Background(), cfg.Target, cfg.dialOptions...); err != nil {
			logger.Panicf("new grpc client error: %s", err)
		} else {
			return &Client{
				cfg:      cfg,
				conn:     conn,
				dialOpts: dialOpts,
			}
		}
	}

	return &Client{
		cfg:      cfg,
		conn:     nil,
		dialOpts: dialOpts,
	}
}

func (client *Client) Name() string {
	return "gRPC-client"
}

func (client *Client) Serve(ctx context.Context) error {
	if conn, err := grpc.DialContext(ctx, client.cfg.Target, client.dialOpts...); err != nil {
		return err
	} else {
		client.conn = conn
	}

	return nil
}

func (client *Client) Shutdown(ctx context.Context) error {
	if client.conn != nil {
		return client.conn.Close()
	}

	return nil
}

func (client *Client) Conn() *grpc.ClientConn {
	return client.conn
}
