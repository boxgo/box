package procsutil

import (
	"github.com/boxgo/box/pkg/logger"
	"go.uber.org/automaxprocs/maxprocs"
)

var (
	err  error
	undo func()
)

// EnableAutoMaxProcs enable
func EnableAutoMaxProcs() {
	if undo, err = maxprocs.Set(); err != nil {
		logger.Errorw("AutoMaxProcs enable error:", "err", err)
	}
}

// DisableAutoMaxProcs disable
func DisableAutoMaxProcs() {
	if undo != nil {
		undo()
	}
}
