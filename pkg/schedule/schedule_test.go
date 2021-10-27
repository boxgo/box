package schedule_test

import (
	"github.com/boxgo/box"
	"github.com/boxgo/box/pkg/logger"
	"github.com/boxgo/box/pkg/schedule"
)

func Example() {
	onceHandler := func(args map[string]interface{}) error {
		logger.Info("once handler executing...", args)
		return nil
	}

	timingHandler := func(args map[string]interface{}) error {
		logger.Info("timing handler executing...", args)
		return nil
	}

	sch := schedule.StdConfig("default").Build(schedule.WithHandler(onceHandler, timingHandler))

	app := box.New(
		box.WithBoxes(sch),
	)

	if err := app.Run(); err != nil {
		logger.Fatal(err)
	}
}
