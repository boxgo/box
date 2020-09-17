package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
)

type (
	Metric struct {
		addr   []string
		master string
		db     int
	}
)

const (
	start = "start"
)

var (
	cmdTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redis_client_command_total",
			Help: "redis command counter",
		},
		[]string{"address", "db", "masterName", "pipe", "cmd", "error"},
	)
	cmdElapsedSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "redis_client_command_duration_seconds",
			Help: "redis command duration seconds",
			Objectives: map[float64]float64{
				0.25: 0.05,
				0.5:  0.05,
				0.75: 0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		},
		[]string{"address", "db", "masterName", "pipe", "cmd", "error"},
	)
)

func init() {
	prometheus.MustRegister(cmdElapsedSummary, cmdTotal)
}

func (m *Metric) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, start, time.Now()), nil
}

func (m *Metric) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	start := ctx.Value(start).(time.Time)
	elapsed := time.Since(start)

	m.report(false, elapsed, cmd)

	return nil
}

func (m *Metric) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, start, time.Now()), nil
}

func (m *Metric) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	start := ctx.Value(start).(time.Time)
	elapsed := time.Since(start)

	m.report(true, elapsed, cmds...)

	return nil
}

func (m *Metric) report(pipe bool, elapsed time.Duration, cmds ...redis.Cmder) {
	addressStr := strings.Join(m.addr, ",")
	dbStr := fmt.Sprintf("%d", m.db)
	masterNameStr := m.master
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

	values := []string{
		addressStr,
		dbStr,
		masterNameStr,
		pipeStr,
		cmdStr,
		errStr,
	}

	cmdElapsedSummary.WithLabelValues(values...).Observe(elapsed.Seconds())
	cmdTotal.WithLabelValues(values...).Inc()
}
