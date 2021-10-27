package server

import "context"

type (
	// Server
	Server interface {
		Name() string                       // server name
		Serve(ctx context.Context) error    // start server
		Shutdown(ctx context.Context) error // gracefully shutdown server
	}

	// Hooks previous and post hooks
	Hook interface {
		BeforeServe(ctx context.Context) error
		AfterServe(ctx context.Context) error
		BeforeShutdown(ctx context.Context) error
		AfterShutdown(ctx context.Context) error
	}
)
