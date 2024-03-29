package metric

import (
	"context"
	"net/http"
	"time"

	"github.com/boxgo/box/v2/config"
	"github.com/boxgo/box/v2/insight"
	"github.com/boxgo/box/v2/logger"
	"github.com/boxgo/box/v2/system"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
)

type (
	// Metric config
	Metric struct {
		cfg  *Config
		stop chan bool
	}
)

var (
	Default = StdConfig("").Build()
)

func newMetric(cfg *Config) *Metric {
	m := &Metric{
		cfg: cfg,
	}

	insight.GetH("/metrics", m.Handler())

	return m
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

	go func() {
		ticker := time.NewTicker(m.cfg.PushInterval)
		defer ticker.Stop()

		pusher := push.
			New(m.cfg.PushTargetURL, config.ServiceName()+"-"+config.ServiceVersion()).
			Gatherer(prometheus.DefaultRegisterer.(prometheus.Gatherer)).
			Grouping("instance", system.Hostname())

		for {
			select {
			case <-m.stop:
				break
			case <-ticker.C:
				if err := pusher.Add(); err != nil {
					logger.Error("metrics.pusher.add.error", err)
				} else {
					logger.Debug("metrics.pusher.add.success")
				}
			}
		}
	}()

	return nil
}

// Shutdown close clients when Shutdown
func (m *Metric) Shutdown(context.Context) error {
	if !m.cfg.PushEnabled {
		return nil
	}

	go func() {
		m.stop <- true
	}()

	return nil
}

// Handler metrics http
func (m *Metric) Handler() http.Handler {
	return promhttp.Handler()
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
