package metric

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
)

type (
	// Metric config
	Metric struct {
		stop chan bool
		cfg  config.SubConfigurator
	}

	Options struct {
		cfg config.SubConfigurator
	}

	OptionFunc func(*Options)
)

const (
	name = "metric"
)

var (
	Default       = New()
	Namespace     = config.NewField(name, "namespace", "metric namespace", "")
	Subsystem     = config.NewField(name, "subsystem", "metric subsystem", "")
	PushEnabled   = config.NewField(name, "pushEnabled", "enable push", false)
	PushTargetURL = config.NewField(name, "pushTargetURL", "pushgateway url", "")
	PushInterval  = config.NewField(name, "pushInterval", "push to a pushgateway interval, millisecond", time.Second*3)
)

// New a metrics
func New(optFunc ...OptionFunc) *Metric {
	opts := &Options{}
	for _, fn := range optFunc {
		fn(opts)
	}

	if opts.cfg == nil {
		opts.cfg = config.Default
	}

	m := &Metric{
		cfg: opts.cfg,
	}

	opts.cfg.Mount(Namespace, Subsystem, PushEnabled, PushTargetURL, PushInterval)

	return m
}

// Name config name
func (m *Metric) Name() string {
	return name
}

// Serve start serve
func (m *Metric) Serve(context.Context) error {
	if !m.cfg.GetBool(PushEnabled) {
		return nil
	}

	hostname, _ := os.Hostname()

	go func() {
		ticker := time.NewTicker(m.cfg.GetDuration(PushInterval) * time.Millisecond)
		defer ticker.Stop()

		pusher := push.
			New(m.cfg.GetString(PushTargetURL), m.cfg.GetBoxName()).
			Gatherer(prometheus.DefaultRegisterer.(prometheus.Gatherer)).
			Grouping("instance", hostname)

		for {
			select {
			case <-m.stop:
				break
			case <-ticker.C:
				if err := pusher.Add(); err != nil {
					logger.Error("metrics.pusher.add.error", err)
				}
			}
		}
	}()

	return nil
}

// Shutdown close clients when Shutdown
func (m *Metric) Shutdown(context.Context) error {
	if !m.cfg.GetBool(PushEnabled) {
		return nil
	}

	go func() {
		m.stop <- true
	}()

	return nil
}

// Metric metrics http
func (m *Metric) Metrics() http.Handler {
	return promhttp.Handler()
}
