package rediscache

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/boxgo/box/pkg/util/strutil"
	"github.com/boxgo/box/pkg/util/testutil"
)

func TestCacheGetSetSimple(t *testing.T) {
	ctx := context.Background()
	testKey := strutil.RandomAlphabet(10)
	testVal := strutil.RandomAlphabet(100)

	err := Set(ctx, testKey, testVal, time.Second*5)
	testutil.ExpectEqual(t, err, nil)

	val := ""
	err = Get(ctx, testKey, &val)
	testutil.ExpectEqual(t, err, nil)
	testutil.ExpectEqual(t, val, testVal)
}

func TestCacheGetSetStruct(t *testing.T) {
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

	err := Set(ctx, testKey, testVal, time.Second*5)
	testutil.ExpectEqual(t, err, nil)

	val := testStruct{}
	err = Get(ctx, testKey, &val)
	a, _ := json.Marshal(val)
	b, _ := json.Marshal(testVal)
	testutil.ExpectEqual(t, err, nil)
	testutil.ExpectEqual(t, a, b)
}
