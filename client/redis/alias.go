package redis

import (
	"github.com/redis/go-redis/v9"
)

type (
	Options = redis.UniversalOptions
)

const (
	Nil         = redis.Nil
	TxFailedErr = redis.TxFailedErr
)

var (
	ErrClosed = redis.ErrClosed
)
