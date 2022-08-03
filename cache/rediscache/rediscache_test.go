package rediscache_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/boxgo/box/v2/cache/rediscache"
	"github.com/boxgo/box/v2/util/strutil"
)

var (
	inst = rediscache.StdConfig("default").Build()
)

func Example() {
	val := ""
	ctx := context.Background()
	testKey := strutil.RandomAlphabet(10)
	testVal := strutil.RandomAlphabet(100)

	if err := inst.Set(ctx, testKey, testVal, time.Second*5); err != nil {
		panic(err)
	}

	if err := inst.Get(ctx, testKey, &val); err != nil {
		panic(err)
	}

	fmt.Println(testVal == val)
	// Output: true
}

func ExampleGet() {
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

	err := inst.Set(ctx, testKey, testVal, time.Second*5)
	if err != nil {
		panic(err)
	}

	val := testStruct{}
	err = inst.Get(ctx, testKey, &val)
	a, _ := json.Marshal(val)
	b, _ := json.Marshal(testVal)

	fmt.Println(bytes.Compare(a, b) == 0)
	// Output: true
}
