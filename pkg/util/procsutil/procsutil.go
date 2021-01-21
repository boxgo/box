package procsutil

import (
	"log"

	"github.com/boxgo/box/pkg/logger"
	"go.uber.org/automaxprocs/maxprocs"
)

var (
	err  error
	undo func()
)

// EnableAutoMaxProcs enable
func EnableAutoMaxProcs() {
	if undo, err = maxprocs.Set(maxprocs.Logger(log.Printf)); err != nil {
		logger.Errorw("AutoMaxProcs enable error:", "err", err)
	}
}

// DisableAutoMaxProcs disable
func DisableAutoMaxProcs() {
	if undo != nil {
		undo()
	}
}
