package redislocker

import (
	"context"
	"testing"
	"time"

	"github.com/boxgo/box/pkg/util/strutil"
	"github.com/boxgo/box/pkg/util/testutil"
)

func TestLock(t *testing.T) {
	ctx := context.Background()
	testKey := strutil.RandomAlphabet(10)

	locked, err := Lock(ctx, testKey, time.Second*5)
	testutil.ExpectEqual(t, err, nil)
	testutil.ExpectEqual(t, locked, true)

	time.Sleep(time.Second * 3)

	locked, err = IsLocked(ctx, testKey)
	testutil.ExpectEqual(t, err, nil)
	testutil.ExpectEqual(t, locked, true)

	time.Sleep(time.Second * 2)
	locked, err = IsLocked(ctx, testKey)
	testutil.ExpectEqual(t, err, nil)
	testutil.ExpectEqual(t, locked, false)
}

func TestUnlock(t *testing.T) {
	ctx := context.Background()
	testKey := strutil.RandomAlphabet(10)

	locked, err := Lock(ctx, testKey, time.Second*5)
	testutil.ExpectEqual(t, err, nil)
	testutil.ExpectEqual(t, locked, true)

	locked, err = IsLocked(ctx, testKey)
	testutil.ExpectEqual(t, err, nil)
	testutil.ExpectEqual(t, locked, true)

	err = UnLock(ctx, testKey)
	testutil.ExpectEqual(t, err, nil)

	locked, err = IsLocked(ctx, testKey)
	testutil.ExpectEqual(t, err, nil)
	testutil.ExpectEqual(t, locked, false)
}
