package metric

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/boxgo/box/internal/basicauth"
	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/insight"
	"github.com/boxgo/box/pkg/logger"
	"github.com/boxgo/box/pkg/system"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
)

type (
	// Metric config
	Metric struct {
		cfg  *Config
		mu   sync.Mutex
		stop chan struct{}
	}
)

var (
	Default = StdConfig("").Build()
)

func init() {
	insight.GetH("/metrics", Default.Handler())
}

// Name config name
func (m *Metric) Name() string {
	return "metric"
}

// Serve start serve
func (m *Metric) Serve(context.Context) error {
	if !m.cfg.PushEnabled {
		return nil
	}
	if m.cfg.PushInterval <= 0 {
		return errors.New("metric pushInterval must be positive")
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if m.stop != nil {
		return nil
	}
	stop := make(chan struct{})
	m.stop = stop

	pusher := push.
		New(m.cfg.PushTargetURL, config.ServiceName()).
		Gatherer(prometheus.DefaultGatherer).
		Grouping("namespace", config.ServiceNamespace()).
		Grouping("instance", system.Hostname())
	if auth := m.cfg.PushBasicAuth; auth.Username != "" || auth.Password != "" {
		pusher.BasicAuth(auth.Username, auth.Password)
	}

	go func() {
		ticker := time.NewTicker(m.cfg.PushInterval)
		defer ticker.Stop()

		for {
			select {
			case <-stop:
				return
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
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.stop != nil {
		close(m.stop)
		m.stop = nil
	}
	return nil
}

// Handler metrics http
func (m *Metric) Handler() http.Handler {
	handler := promhttp.Handler()
	if len(m.cfg.BasicAuth) == 0 {
		return handler
	}

	verify := basicauth.NewVerifier(m.cfg.BasicAuth)
	challenge := basicauth.Challenge("metrics")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := verify(r); ok {
			handler.ServeHTTP(w, r)
			return
		}

		w.Header().Set("WWW-Authenticate", challenge)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	})
}

// NewCounterVec creates a new CounterVec based on the provided CounterOpts and partitioned by the given label names.
func (m *Metric) NewCounterVec(name, help string, labels []string) *CounterVec {
	vec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: m.cfg.Namespace,
			Subsystem: m.cfg.Subsystem,
			Name:      name,
			Help:      help,
		},
		labels,
	)

	register(vec)

	return vec
}

// NewSummaryVec creates a new SummaryVec based on the provided SummaryOpts and partitioned by the given label names.
func (m *Metric) NewSummaryVec(name, help string, labels []string, objectives map[float64]float64) *SummaryVec {
	vec := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  m.cfg.Namespace,
			Subsystem:  m.cfg.Subsystem,
			Name:       name,
			Help:       help,
			Objectives: objectives,
		},
		labels,
	)

	register(vec)

	return vec
}

// NewGaugeVec creates a new GaugeVec based on the provided GaugeOpts and partitioned by the given label names.
func (m *Metric) NewGaugeVec(name, help string, labels []string) *GaugeVec {
	vec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.cfg.Namespace,
			Subsystem: m.cfg.Subsystem,
			Name:      name,
			Help:      help,
		},
		labels,
	)

	register(vec)

	return vec
}

// NewHistogramVec creates a new HistogramVec based on the provided HistogramOpts and partitioned by the given label names.
func (m *Metric) NewHistogramVec(name, help string, labels []string, buckets []float64) *HistogramVec {
	vec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: m.cfg.Namespace,
			Subsystem: m.cfg.Subsystem,
			Name:      name,
			Help:      help,
			Buckets:   buckets,
		},
		labels,
	)

	register(vec)

	return vec
}

func register(cs prometheus.Collector) {
	if err := prometheus.Register(cs); err != nil {
		logger.Warnw("metric register error", "err", err)
	}
}
