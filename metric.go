package box

import (
	"context"
	"time"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/metric"
	"github.com/boxgo/box/pkg/system"
)

type (
	boxMetric struct{}
)

var (
	boxMetricGauge = metric.Default.NewGaugeVec(
		"box_info",
		"Information about the box config and environment.",
		[]string{"name", "version", "ip", "localhost", "start"})
)

func (boxMetric) Name() string {
	return "box-metric"
}

func (boxMetric) Serve(ctx context.Context) error {
	boxMetricGauge.WithLabelValues(
		config.ServiceName(),
		config.ServiceVersion(),
		system.IP(),
		system.Hostname(),
		system.StartAt().Format(time.RFC3339),
	).Inc()

	return nil
}

func (boxMetric) Shutdown(ctx context.Context) error {
	return nil
}
