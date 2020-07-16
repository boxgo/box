package server

import "context"

type (
	// Server
	Server interface {
		Name() string                       // server name
		Serve(ctx context.Context) error    // start server
		Shutdown(ctx context.Context) error // gracefully shutdown server
	}
)
