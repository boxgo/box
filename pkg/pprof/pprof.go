package pprof

import (
	"context"
	"net/http/pprof"

	"github.com/boxgo/box/pkg/insight"
)

type (
	PProf struct {
		cfg *Config
	}
)

var (
	Default = StdConfig().Build()
)

func newPProf(cfg *Config) *PProf {
	pp := &PProf{
		cfg: cfg,
	}

	insight.GetF("/debug/pprof/", pprof.Index)
	insight.GetF("/debug/pprof/allocs", pprof.Index)
	insight.GetF("/debug/pprof/block", pprof.Index)
	insight.GetF("/debug/pprof/cmdline", pprof.Cmdline)
	insight.GetF("/debug/pprof/goroutine", pprof.Index)
	insight.GetF("/debug/pprof/heap", pprof.Index)
	insight.GetF("/debug/pprof/mutex", pprof.Index)
	insight.GetF("/debug/pprof/profile", pprof.Profile)
	insight.GetF("/debug/pprof/threadcreate", pprof.Index)
	insight.GetF("/debug/pprof/symbol", pprof.Symbol)
	insight.GetF("/debug/pprof/trace", pprof.Trace)

	return pp
}

func (pp *PProf) Name() string {
	return "pprof"
}

func (pp *PProf) Serve(ctx context.Context) error {
	return nil
}

func (pp *PProf) Shutdown(ctx context.Context) error {
	return nil
}
