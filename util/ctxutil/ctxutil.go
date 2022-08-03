package ctxutil

import (
	"context"
	"time"
)

func WithCancel() (context.Context, func()) {
	return context.WithCancel(context.Background())
}

func WithDeadline(d time.Time) (context.Context, func()) {
	return context.WithDeadline(context.Background(), d)
}

func WithTimeout(timeout time.Duration) (context.Context, func()) {
	return context.WithTimeout(context.Background(), timeout)
}
