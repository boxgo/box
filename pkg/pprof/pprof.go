package pprof

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/pprof"
	"net/url"
	"strconv"

	"github.com/boxgo/box/pkg/client/wukong"
	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/insight"
	"github.com/boxgo/box/pkg/server/ginserver"
	"github.com/boxgo/box/pkg/system"
	"golang.org/x/sync/errgroup"
)

type (
	PProf struct {
		cfg *Config
	}

	StartReq struct {
		// post pprof file to the target
		Target string `json:"target" binding:"required"`

		// pprof command
		Profiles []string `form:"profiles" binding:"required"`

		// debug=N (all profiles): response format: N = 0: binary (default), N > 0: plaintext
		Debug *int `form:"debug"`

		// gc=N (heap profile): N > 0: run a garbage collection cycle before profiling
		GC *int `form:"gc"`

		// seconds=N (allocs, block, goroutine, heap, mutex, threadcreate profiles): return a delta profile
		// seconds=N (cpu (profile), trace profiles): profile for the given duration
		Seconds *int `form:"seconds"`
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
	insight.GetF("/debug/pprof/trace", pprof.Trace)
	insight.GetF("/debug/pprof/symbol", pprof.Symbol)
	insight.PostF("/debug/pprof/symbol", pprof.Symbol)

	insight.Post("/debug/pprof/start", pp.Start)

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

func (pp *PProf) Start(ctx *ginserver.Context) {
	req := &StartReq{
		Profiles: []string{"profile", "allocs", "heap", "block", "mutex", "goroutine", "threadcreate"},
	}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.String(400, "param error: %s", err.Error())
		return
	}

	var (
		eg, newCtx = errgroup.WithContext(ctx)
	)

	query := url.Values{}
	if req.Debug != nil {
		query.Add("debug", strconv.Itoa(*req.Debug))
	}
	if req.GC != nil {
		query.Add("gc", strconv.Itoa(*req.GC))
	}
	if req.Seconds != nil {
		query.Add("seconds", strconv.Itoa(*req.Seconds))
	}

	for _, profile := range req.Profiles {
		profile := profile

		eg.Go(func() error {
			if data, err := pp.fetchProfileData(newCtx, profile, query.Encode()); err != nil {
				return err
			} else if err = pp.uploadProfileData(newCtx, req.Target, profile, data); err != nil {
				return err
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		ctx.String(500, "profiling error: %s", err.Error())
	} else {
		ctx.String(200, "success")
	}
}

func (pp *PProf) fetchProfileData(ctx context.Context, profile, query string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/debug/pprof/%s?%s", profile, query), nil)
	if err != nil {
		return nil, err
	} else {
		req = req.WithContext(ctx)
	}

	rw := httptest.NewRecorder()
	insight.ServeHTTP(rw, req)

	return rw.Body.Bytes(), nil
}

func (pp *PProf) uploadProfileData(ctx context.Context, target, profile string, data []byte) error {
	return wukong.
		New("").
		Post(target).
		WithCTX(ctx).
		SendFileReader("profile", profile, bytes.NewBuffer(data)).
		Form(map[string]string{
			"namespace": config.ServiceNamespace(),
			"service":   config.ServiceName(),
			"version":   config.ServiceVersion(),
			"ip":        system.IP(),
			"hostname":  system.Hostname(),
		}).
		End().
		Error()
}
