package pprof

import (
	"context"
	"net/http"
	"net/http/pprof"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/dummybox"
)

type (
	PProf struct {
		dummybox.DummyBox
		cfg      config.SubConfigurator
		serveMux *http.ServeMux
		addr     *config.Field
	}

	Options struct {
		cfg config.SubConfigurator
	}

	OptionFunc func(*Options)
)

var (
	Default = New()
)

func WithConfigurator(cfg config.SubConfigurator) OptionFunc {
	return func(opts *Options) {
		opts.cfg = cfg
	}
}

func New(optionFunc ...OptionFunc) *PProf {
	opts := &Options{}
	for _, fn := range optionFunc {
		fn(opts)
	}

	if opts.cfg == nil {
		opts.cfg = config.Default
	}

	pp := &PProf{
		cfg:  opts.cfg,
		addr: config.NewField("pprof", "addr", "pprof http server listen addr", ":9091"),
	}

	pp.cfg.Mount(pp.addr)

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

	return http.ListenAndServe(pp.cfg.GetString(pp.addr), serveMux)
}

func (pp *PProf) Shutdown(ctx context.Context) error {
	return nil
}
