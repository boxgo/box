package redislocker_test

import (
	"context"
	"testing"
	"time"

	"github.com/boxgo/box/pkg/locker/redislocker"
	"github.com/boxgo/box/pkg/util/strutil"
	"github.com/boxgo/box/pkg/util/testutil"
)

var (
	inst = redislocker.StdConfig("default").Build()
)

func TestLock(t *testing.T) {
	ctx := context.Background()
	testKey := strutil.RandomAlphabet(10)

	locked, err := inst.Lock(ctx, testKey, time.Second*5)
	testutil.ExpectEqual(t, err, nil)
	testutil.ExpectEqual(t, locked, true)

	time.Sleep(time.Second * 3)

	locked, err = inst.IsLocked(ctx, testKey)
	testutil.ExpectEqual(t, err, nil)
	testutil.ExpectEqual(t, locked, true)

	time.Sleep(time.Second * 2)
	locked, err = inst.IsLocked(ctx, testKey)
	testutil.ExpectEqual(t, err, nil)
	testutil.ExpectEqual(t, locked, false)
}

func TestUnlock(t *testing.T) {
	ctx := context.Background()
	testKey := strutil.RandomAlphabet(10)

	locked, err := inst.Lock(ctx, testKey, time.Second*5)
	testutil.ExpectEqual(t, err, nil)
	testutil.ExpectEqual(t, locked, true)

	locked, err = inst.IsLocked(ctx, testKey)
	testutil.ExpectEqual(t, err, nil)
	testutil.ExpectEqual(t, locked, true)

	err = inst.UnLock(ctx, testKey)
	testutil.ExpectEqual(t, err, nil)

	locked, err = inst.IsLocked(ctx, testKey)
	testutil.ExpectEqual(t, err, nil)
	testutil.ExpectEqual(t, locked, false)
}
