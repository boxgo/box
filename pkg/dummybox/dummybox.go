package dummybox

import (
	"context"
)

type (
	DummyBox struct{}
)

func (*DummyBox) Name() string {
	return "dummyBox"
}

func (*DummyBox) Serve(context.Context) error {
	return nil
}

func (*DummyBox) Shutdown(context.Context) error {
	return nil
}
