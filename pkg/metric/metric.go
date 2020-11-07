package metric

import (
	"context"
	"net/http"
	"time"

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
		stop chan bool
	}
)

var (
	Default = DefaultConfig().Build()
)

func newMetric(cfg *Config) *Metric {
	m := &Metric{
		cfg: cfg,
	}

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
			New(m.cfg.PushTargetURL, system.ServiceName()).
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
func (m *Metric) NewCounterVec(name, help string, labels []string) *prometheus.CounterVec {
	vec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: m.cfg.Namespace,
			Subsystem: m.cfg.Subsystem,
			Name:      name,
			Help:      help,
		},
		labels,
	)

	prometheus.MustRegister(vec)

	return vec
}

// NewSummaryVec creates a new SummaryVec based on the provided SummaryOpts and partitioned by the given label names.
func (m *Metric) NewSummaryVec(name, help string, labels []string, objectives map[float64]float64) *prometheus.SummaryVec {
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

	prometheus.MustRegister(vec)

	return vec
}

// NewGaugeVec creates a new GaugeVec based on the provided GaugeOpts and partitioned by the given label names.
func (m *Metric) NewGaugeVec(name, help string, labels []string) *prometheus.GaugeVec {
	vec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.cfg.Namespace,
			Subsystem: m.cfg.Subsystem,
			Name:      name,
			Help:      help,
		},
		labels,
	)

	prometheus.MustRegister(vec)

	return vec
}

// NewHistogramVec creates a new HistogramVec based on the provided HistogramOpts and partitioned by the given label names.
func (m *Metric) NewHistogramVec(name, help string, labels []string, buckets []float64) *prometheus.HistogramVec {
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

	prometheus.MustRegister(vec)

	return vec
}
