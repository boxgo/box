package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/boxgo/box/pkg/metric"
	"github.com/boxgo/box/pkg/trace"
	"github.com/go-redis/redis/v8"
)

type (
	Metric struct {
		cfg *Config
	}

	startKey struct{}
)

var (
	cmdTotal = metric.NewCounterVec(
		"redis_client_command_total",
		"redis command counter",
		[]string{"bid", "address", "db", "masterName", "pipe", "cmd", "error"},
	)
	cmdDuration = metric.NewSummaryVec(
		"redis_client_command_duration_seconds",
		"redis command duration seconds",
		[]string{"bid", "address", "db", "masterName", "pipe", "cmd", "error"},
		map[float64]float64{
			0.5:  0.05,
			0.75: 0.05,
			0.9:  0.01,
			0.99: 0.001,
			1:    0.001,
		},
	)
)

func (m *Metric) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, startKey{}, time.Now()), nil
}

func (m *Metric) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	start := ctx.Value(startKey{}).(time.Time)
	elapsed := time.Since(start)

	m.report(ctx, false, elapsed, cmd)

	return nil
}

func (m *Metric) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, startKey{}, time.Now()), nil
}

func (m *Metric) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	start := ctx.Value(startKey{}).(time.Time)
	elapsed := time.Since(start)

	m.report(ctx, true, elapsed, cmds...)

	return nil
}

func (m *Metric) report(ctx context.Context, pipe bool, elapsed time.Duration, cmds ...redis.Cmder) {
	addressStr := strings.Join(m.cfg.Address, ",")
	dbStr := fmt.Sprintf("%d", m.cfg.DB)
	masterNameStr := m.cfg.MasterName
	errStr := ""
	cmdStr := ""
	pipeStr := fmt.Sprintf("%t", pipe)

	for _, cmd := range cmds {
		cmdStr += cmd.Name() + ";"

		if err := cmd.Err(); err != nil && err != redis.Nil {
			errStr += err.Error() + ";"
		}
	}
	cmdStr = strings.TrimSuffix(cmdStr, ";")

	var (
		bizID string
	)

	if bizIDStr, ok := ctx.Value(trace.BizID()).(string); ok {
		bizID = bizIDStr
	}

	values := []string{
		bizID,
		addressStr,
		dbStr,
		masterNameStr,
		pipeStr,
		cmdStr,
		errStr,
	}

	cmdDuration.WithLabelValues(values...).Observe(elapsed.Seconds())
	cmdTotal.WithLabelValues(values...).Inc()
}
