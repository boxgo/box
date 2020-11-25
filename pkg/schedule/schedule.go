package schedule

import (
	"context"

	"github.com/boxgo/box/pkg/locker"
	"github.com/boxgo/box/pkg/locker/redislocker"
	"github.com/boxgo/box/pkg/logger"
	"github.com/robfig/cron"
)

type (
	Schedule struct {
		cfg           *Config
		cron          *cron.Cron
		locker        locker.MutexLocker
		onceHandler   Handler
		timingHandler Handler
	}

	Handler func() error
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

	switch sch.cfg.Type {
	case Once:
		sch.execOnce()
	case Timing:
		sch.execTiming()
	case OnceAndTiming:
		sch.execOnce()
		sch.execTiming()
	}

	return sch
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
				if err := sch.locker.UnLock(context.Background(), sch.cfg.key); err != nil {
					logger.Errorf("Schedule [%s] unlock error: [%s]", sch.cfg.key, err)
				}
			}

			if err := recover(); err != nil {
				logger.Errorf("Schedule [%s] crash: %s", sch.cfg.key, err)
				return
			}
		}()

		if sch.cfg.Compete {
			locked, err := sch.locker.Lock(context.Background(), sch.cfg.key, sch.cfg.lockDuration)
			if err != nil {
				logger.Errorf("Schedule [%s] compete error: [%s]", sch.cfg.key, err)
				return
			} else if !locked {
				logger.Warnf("Schedule [%s] compete fail", sch.cfg.key)
				return
			} else {
				logger.Infof("Schedule [%s] compete win", sch.cfg.key)
			}
		}

		logger.Infof("Schedule [%s] run start", sch.cfg.key)

		if err := handler(); err != nil {
			logger.Errorf("Schedule [%s] run error: [%s]", sch.cfg.key, err)
		} else {
			logger.Infof("Schedule [%s] run success", sch.cfg.key)
		}
	}()
}
