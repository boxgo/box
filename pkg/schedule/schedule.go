// Package schedule is to help you manage schedule tasks.
package schedule

import (
	"context"
	"fmt"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/locker"
	"github.com/boxgo/box/pkg/locker/redislocker"
	"github.com/boxgo/box/pkg/logger"
	"github.com/boxgo/box/pkg/metric"
	"github.com/robfig/cron"
)

type (
	// Schedule instance is a server, you should mount bo box application to manage lifecycle.
	Schedule struct {
		cfg           *Config
		cron          *cron.Cron
		locker        locker.MutexLocker
		onceHandler   Handler
		timingHandler Handler
	}

	Handler func(args map[string]interface{}) error
)

var (
	scheduleCounter = metric.NewCounterVec(
		"schedule_total",
		"schedule counter",
		[]string{"task", "error", "panic"},
	)
)

func newSchedule(cfg *Config) *Schedule {
	sch := &Schedule{
		cfg:           cfg,
		cron:          cron.New(),
		locker:        cfg.locker,
		onceHandler:   cfg.onceHandler,
		timingHandler: cfg.timingHandler,
	}

	if sch.locker == nil {
		sch.locker = redislocker.Default
	}

	if sch.onceHandler == nil && sch.timingHandler == nil {
		logger.Panic("schedule handler is not set")
	}

	if sch.cfg.Type != Once {
		if err := sch.cron.AddFunc(sch.cfg.Spec, func() {
			sch.exec(sch.timingHandler)
		}); err != nil {
			logger.Panicf("schedule build error: %s", err)
		}
	}

	return sch
}

func (sch *Schedule) Name() string {
	return "schedule"
}

// Serve schedule
func (sch *Schedule) Serve(context.Context) error {
	switch sch.cfg.Type {
	case Once:
		sch.execOnce()
	case Timing:
		sch.execTiming()
	case OnceAndTiming:
		sch.execOnce()
		sch.execTiming()
	}

	return nil
}

// Shutdown stop cron
func (sch *Schedule) Shutdown(context.Context) error {
	if sch.cron != nil {
		sch.cron.Stop()
	}

	return nil
}

// ExecOnce exec once handler immediately
func (sch *Schedule) ExecOnce() {
	sch.execOnce()
}

// ExecTiming exec timing handler immediately
func (sch *Schedule) ExecTiming() {
	if sch.timingHandler == nil {
		return
	}

	sch.exec(sch.timingHandler)
}

func (sch *Schedule) execOnce() {
	if sch.onceHandler == nil {
		return
	}

	sch.exec(sch.onceHandler)
}

func (sch *Schedule) execTiming() {
	if sch.timingHandler == nil {
		return
	}

	sch.cron.Start()
}

func (sch *Schedule) exec(handler Handler) {
	go func() {
		defer func() {
			if sch.cfg.Compete && sch.cfg.AutoUnlock {
				if err := sch.locker.UnLock(context.Background(), sch.key()); err != nil {
					logger.Errorf("Schedule [%s] unlock error: [%s]", sch.key(), err)
				}
			}

			if err := recover(); err != nil {
				scheduleCounter.WithLabelValues(sch.key(), "", fmt.Sprintf("%s", err)).Inc()
				logger.Errorf("Schedule [%s] crash: %s", sch.key(), err)
				return
			}
		}()

		if sch.cfg.Compete {
			locked, err := sch.locker.Lock(context.Background(), sch.key(), sch.cfg.lockDuration)
			if err != nil {
				logger.Errorf("Schedule [%s] compete error: [%s]", sch.key(), err)
				return
			} else if !locked {
				logger.Warnf("Schedule [%s] compete fail", sch.key())
				return
			} else {
				logger.Infof("Schedule [%s] compete win", sch.key())
			}
		}

		logger.Infof("Schedule [%s] run start", sch.key())

		if err := handler(sch.cfg.Args); err != nil {
			scheduleCounter.WithLabelValues(sch.key(), err.Error(), "").Inc()
			logger.Errorf("Schedule [%s] run error: [%s]", sch.key(), err)
		} else {
			scheduleCounter.WithLabelValues(sch.key(), "", "").Inc()
			logger.Infof("Schedule [%s] run success", sch.key())
		}
	}()
}

func (sch *Schedule) key() string {
	return fmt.Sprintf("%s.%s", config.ServiceName(), sch.cfg.path)
}
