package rediscache_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/boxgo/box/v2/client/redis"
	"github.com/boxgo/box/v2/contrib/cache/rediscache"
	"github.com/boxgo/box/v2/util/strutil"
)

func TestRedisCache(t *testing.T) {
	val := ""
	ctx := context.Background()
	testKey := strutil.RandomAlphabet(10)
	testVal := strutil.RandomAlphabet(100)

	cache, err := rediscache.New(redis.New(&redis.Options{
		Addrs: []string{"127.0.0.1:6379"},
	}))

	if err != nil {
		t.Fatal(err)
	}

	if err = cache.Set(ctx, testKey, testVal, time.Second*5); err != nil {
		t.Fatal(err)
	}

	if err = cache.Get(ctx, testKey, &val); err != nil {
		t.Fatal(err)
	}

	if testVal != val {
		t.Fatalf("RedisCache\nexpect:%s\nreal:%s", testVal, val)
	}
}

func TestRedisCache1(t *testing.T) {
	type (
		testStruct struct {
			String      string
			Int         int
			StringSlice []string
			Boolean     bool
			Map         map[string]interface{}
		}
	)

	ctx := context.Background()
	testKey := strutil.RandomAlphabet(10)
	testVal := testStruct{
		String:      "box",
		Int:         1234,
		StringSlice: []string{"tom", "andy", "luck"},
		Boolean:     true,
		Map: map[string]interface{}{
			"int":      99,
			"sliceInt": []int{1, 2, 3},
			"mapStringString": map[string]string{
				"key1": "val1",
				"key2": "val2",
			},
		},
	}

	cache, err := rediscache.New(redis.New(&redis.Options{
		Addrs: []string{"127.0.0.1:6379"},
	}))

	err = cache.Set(ctx, testKey, testVal, time.Second*5)
	if err != nil {
		t.Fatal(err)
	}

	val := testStruct{}
	err = cache.Get(ctx, testKey, &val)
	a, _ := json.Marshal(val)
	b, _ := json.Marshal(testVal)

	if bytes.Compare(a, b) != 0 {
		t.Fatalf("RedisCache\nexpect:%s\nreal:%s", a, b)
	}
}
