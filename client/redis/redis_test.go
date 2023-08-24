package redis_test

import (
	"context"
	"fmt"
	"time"

	"github.com/boxgo/box/v2/client/redis"
)

func Example() {
	ctx := context.TODO()

	client := redis.New(&redis.Options{
		Addrs: []string{"127.0.0.1:6379"},
	})

	if err := client.Set(ctx, "key", "value", time.Minute).Err(); err != nil {
		panic(err)
	}

	if val, err := client.Get(ctx, "key").Result(); err != nil {
		panic(err)
	} else if val != "value" {
		panic(err)
	} else {
		fmt.Println(val)
	}

	if err := client.Del(ctx, "key").Err(); err != nil {
		panic(err)
	}

	// Output:
	// value
}

func ExampleRedis_NewScript() {
	ctx := context.TODO()

	client := redis.New(&redis.Options{
		Addrs: []string{"127.0.0.1:6379"},
	})

	IncrByXX := client.NewScript(`
		return redis.call("INCRBY", KEYS[1], ARGV[1])
	`)

	defer client.Del(ctx, "xx_counter")

	if n, err := IncrByXX.Run(ctx, []string{"xx_counter"}, 100).Result(); err != nil {
		panic(err)
	} else {
		fmt.Println(n)
	}

	if err := client.Set(ctx, "xx_counter", "40", 0).Err(); err != nil {
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
