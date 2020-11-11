package pprof

import (
	"context"
	"net/http"
	"net/http/pprof"
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

	return pp
}

func (pp *PProf) Name() string {
	return "pprof"
}

func (pp *PProf) Serve(ctx context.Context) error {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/debug/pprof/", pprof.Index)
	serveMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	serveMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	serveMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	serveMux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return http.ListenAndServe(pp.cfg.Addr, serveMux)
}

func (pp *PProf) Shutdown(ctx context.Context) error {
	return nil
}
