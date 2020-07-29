package mongodb

import (
	"context"
	"sync"
	"time"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Monitor struct {
		once               sync.Once
		stopWatch          chan bool
		cfg                config.SubConfigurator
		cmdTotal           *prometheus.CounterVec
		cmdDurationSummary *prometheus.SummaryVec
		workingSession     *prometheus.GaugeVec
	}
)

func newMonitor(cfg config.SubConfigurator) *Monitor {
	m := &Monitor{
		cfg:       cfg,
		stopWatch: make(chan bool),
	}

	return m
}

func (mon *Monitor) CommandMonitor() *event.CommandMonitor {
	return &event.CommandMonitor{
		Started:   mon.started,
		Succeeded: mon.succeeded,
		Failed:    mon.failed,
	}
}

func (mon *Monitor) PoolEventMonitor() *event.PoolMonitor {
	return &event.PoolMonitor{
		Event: mon.event,
	}
}

func (mon *Monitor) init() {
	mon.once.Do(func() {
		mon.cmdTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "",
				Subsystem: "",
				Name:      "mongo_client_command_total",
				Help:      "mongodb client command cmdTotal",
			},
			[]string{"command", "error"},
		)
		mon.cmdDurationSummary = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace: "",
				Subsystem: "",
				Name:      "mongo_client_command_duration_summary",
				Help:      "mongodb client command duration cmdDurationSummary",
			},
			[]string{"command", "error"},
		)
		mon.workingSession = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "",
				Subsystem: "",
				Name:      "mongo_client_session_in_progress",
				Help:      "mongo client session in progress gauge",
			},
			[]string{},
		)

		prometheus.MustRegister(mon.cmdDurationSummary, mon.cmdTotal, mon.workingSession)
	})
}

func (mon *Monitor) started(ctx context.Context, ev *event.CommandStartedEvent) {
	mon.init()

	logger.Trace(ctx).Debugf("mongo_command_start cmd: %s, reqId: %d, connId: %s, db: %s", ev.CommandName, ev.RequestID, ev.ConnectionID, ev.DatabaseName)
}

func (mon *Monitor) succeeded(ctx context.Context, ev *event.CommandSucceededEvent) {
	mon.init()

	labels := []string{ev.CommandName, ""}
	mon.cmdTotal.WithLabelValues(labels...).Inc()
	mon.cmdDurationSummary.WithLabelValues(labels...).Observe(float64(ev.DurationNanos))

	logger.Trace(ctx).Debugf("mongo_command_success cmd: %s, reqId: %d, connId: %s, duration: %s", ev.CommandName, ev.RequestID, ev.ConnectionID, time.Duration(ev.DurationNanos))
}

func (mon *Monitor) failed(ctx context.Context, ev *event.CommandFailedEvent) {
	mon.init()

	labels := []string{ev.CommandName, ev.Failure}
	mon.cmdTotal.WithLabelValues(labels...).Inc()
	mon.cmdDurationSummary.WithLabelValues(labels...).Observe(float64(ev.DurationNanos))

	logger.Trace(ctx).Debugf("mongo_command_error cmd: %s, reqId: %d, connId: %s, duration: %d, error: %s", ev.CommandName, ev.RequestID, ev.ConnectionID, ev.DurationNanos, ev.Failure)
}

func (mon *Monitor) event(ev *event.PoolEvent) {
	mon.init()

	logger.Debugf("mongo_pool_event type: %s, address: %s, connId: %d, reason: %s", ev.Type, ev.Address, ev.ConnectionID, ev.Reason)
}

func (mon *Monitor) watch(db *mongo.Client) {
	go func() {
		for {
			time.Sleep(time.Second)

			select {
			case <-mon.stopWatch:
				logger.Debugf("mongo monitor watch exit")
				close(mon.stopWatch)
				break
			default:
				mon.workingSession.WithLabelValues().Set(float64(db.NumberSessionsInProgress()))
			}
		}
	}()
}

func (mon *Monitor) shutdown() {
	mon.stopWatch <- true
}
