package dummybox

import (
	"context"

	"github.com/boxgo/box/pkg/config"
)

type (
	DummyBox struct{}
)

func (*DummyBox) Name() string {
	return "dummyBox"
}

func (*DummyBox) Configure(config.WatchMountGetter) {
}

func (*DummyBox) Serve(context.Context) error {
	return nil
}

func (*DummyBox) Shutdown(context.Context) error {
	return nil
}
