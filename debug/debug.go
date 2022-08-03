package debug

import (
	"context"
	"expvar"

	"github.com/boxgo/box/v2/insight"
)

type (
	Debug struct {
		cfg *Config
	}
)

var (
	Default = StdConfig().Build()
)

func newDebug(cfg *Config) *Debug {
	debug := &Debug{
		cfg: cfg,
	}

	insight.GetH("/debug/vars", expvar.Handler())

	return debug
}

func (pp *Debug) Name() string {
	return "debug"
}

func (pp *Debug) Serve(ctx context.Context) error {
	return nil
}

func (pp *Debug) Shutdown(ctx context.Context) error {
	return nil
}
