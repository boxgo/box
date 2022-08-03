package metric

import (
	"net/http"
)

// Handler metrics http
func Handler() http.Handler {
	return Default.Handler()
}

// NewCounterVec creates a new CounterVec based on the provided CounterOpts and partitioned by the given label names.
func NewCounterVec(name, help string, labels []string) *CounterVec {
	return Default.NewCounterVec(name, help, labels)
}

// NewSummaryVec creates a new SummaryVec based on the provided SummaryOpts and partitioned by the given label names.
func NewSummaryVec(name, help string, labels []string, objectives map[float64]float64) *SummaryVec {
	return Default.NewSummaryVec(name, help, labels, objectives)
}

// NewGaugeVec creates a new GaugeVec based on the provided GaugeOpts and partitioned by the given label names.
func NewGaugeVec(name, help string, labels []string) *GaugeVec {
	return Default.NewGaugeVec(name, help, labels)
}

// NewHistogramVec creates a new HistogramVec based on the provided HistogramOpts and partitioned by the given label names.
func NewHistogramVec(name, help string, labels []string, buckets []float64) *HistogramVec {
	return Default.NewHistogramVec(name, help, labels, buckets)
}
