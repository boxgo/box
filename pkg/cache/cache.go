package cache

import (
	"context"
	"errors"
	"time"
)

type (
	Cache interface {
		Get(context.Context, string, interface{}) error
		Set(context.Context, string, interface{}, time.Duration) error
		Clear(context.Context, string) error
		Expire(context.Context, string, time.Duration) error
	}
)

var (
	ErrCacheMiss = errors.New("cache: key is missing")
)
