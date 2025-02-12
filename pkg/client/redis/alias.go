package redis

import (
	"github.com/redis/go-redis/v9"
)

const (
	Nil         = redis.Nil
	TxFailedErr = redis.TxFailedErr
)

var (
	ErrClosed = redis.ErrClosed
)
