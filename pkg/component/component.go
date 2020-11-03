package component

import (
	"context"

	"github.com/boxgo/box/pkg/server"
)

type (
	// Box component interface
	Box interface {
		server.Server
	}

	NoopBox struct{}
)

func (*NoopBox) Name() string {
	return "noop-box"
}

func (*NoopBox) Serve(context.Context) error {
	return nil
}

func (*NoopBox) Shutdown(context.Context) error {
	return nil
}
