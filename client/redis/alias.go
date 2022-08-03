package redis

import (
	"github.com/go-redis/redis/v8"
)

const (
	Nil         = redis.Nil
	TxFailedErr = redis.TxFailedErr
)

var (
	ErrClosed = redis.ErrClosed
)
