package redis

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/boxgo/box/pkg/metric"
	"github.com/redis/go-redis/v9"
)

type (
	Metric struct {
		cfg  *Config
		addr string
	}

	startKey struct{}
)

func newMetric(cfg *Config) *Metric {
	return &Metric{
		cfg:  cfg,
		addr: strings.Join(cfg.Address, ","),
	}
}

var (
	cmdTotal = metric.NewCounterVec(
		"redis_client_requests_total",
		"The total number of Redis commands executed.",
		[]string{"address", "db", "masterName", "pipe", "cmd", "result"},
	)
	cmdDuration = metric.NewHistogramVec(
		"redis_client_request_duration_seconds",
		"The Redis command latencies in seconds.",
		[]string{"address", "db", "masterName", "pipe", "cmd", "result"},
		// 100us, 250us, 500us, 1ms, 2.5ms, 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms
		[]float64{0.0001, 0.00025, 0.0005, 0.001, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5},
	)
)

func (m *Metric) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

func (m *Metric) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()

		err := next(ctx, cmd)

		m.report(ctx, false, time.Since(start), cmd)

		return err
	}
}

func (m *Metric) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		start := time.Now()

		err := next(ctx, cmds)

		m.report(ctx, false, time.Since(start), cmds...)

		return err
	}
}

func (m *Metric) report(ctx context.Context, pipe bool, elapsed time.Duration, cmds ...redis.Cmder) {
	cmdStr := ""
	result := "success"
	masterNameStr := m.cfg.MasterName
	addressStr := m.addr
	dbStr := strconv.Itoa(m.cfg.DB)
	pipeStr := strconv.FormatBool(pipe)

	if pipe {
		cmdStr = "pipeline"
	} else if len(cmds) > 0 {
		cmdStr = cmds[0].Name()
	}

	for _, cmd := range cmds {
		if err := cmd.Err(); err != nil && err != redis.Nil {
			result = "error"
			break
		}
	}

	values := []string{
		addressStr,
		dbStr,
		masterNameStr,
		pipeStr,
		cmdStr,
		result,
	}

	cmdDuration.WithLabelValues(values...).Observe(elapsed.Seconds())
	cmdTotal.WithLabelValues(values...).Inc()
}
