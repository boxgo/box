package box

import (
	"context"
	"strings"
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
		[]string{"tags", "ip", "localhost", "start"})
)

func (boxMetric) Name() string {
	return "box-metric"
}

func (boxMetric) Serve(ctx context.Context) error {
	boxMetricGauge.WithLabelValues(
		strings.Join(config.ServiceTag(), ","),
		system.IP(),
		system.Hostname(),
		system.StartAt().Format(time.RFC3339),
	).Inc()

	return metric.Default.Serve(ctx)
}

func (boxMetric) Shutdown(ctx context.Context) error {
	return metric.Default.Shutdown(ctx)
}
