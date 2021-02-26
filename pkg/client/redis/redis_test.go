package redis_test

import (
	"context"
	"fmt"
	"time"

	"github.com/boxgo/box/pkg/client/redis"
)

func Example() {
	ctx := context.TODO()

	if err := redis.Set(ctx, "key", "value", time.Minute).Err(); err != nil {
		panic(err)
	}

	if err := redis.Del(ctx, "key").Err(); err != nil {
		panic(err)
	}
}

func ExampleRedis_NewScript() {
	ctx := context.TODO()

	IncrByXX := redis.Default.NewScript(`
		return redis.call("INCRBY", KEYS[1], ARGV[1])
	`)

	defer redis.Del(ctx, "xx_counter")

	if n, err := IncrByXX.Run(ctx, []string{"xx_counter"}, 100).Result(); err != nil {
		panic(err)
	} else {
		fmt.Println(n)
	}

	if err := redis.Set(ctx, "xx_counter", "40", 0).Err(); err != nil {
		panic(err)
	}

	if n, err := IncrByXX.Run(ctx, []string{"xx_counter"}, 2).Result(); err != nil {
		panic(err)
	} else {
		fmt.Println(n)
	}

	// Output:
	// 100
	// 42
}
