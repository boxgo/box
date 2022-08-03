package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	CounterVec   = prometheus.CounterVec
	SummaryVec   = prometheus.SummaryVec
	GaugeVec     = prometheus.GaugeVec
	HistogramVec = prometheus.HistogramVec
)
