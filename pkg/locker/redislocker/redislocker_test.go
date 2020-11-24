package redislocker

import (
	"testing"
	"time"

	"github.com/boxgo/box/pkg/util/strutil"
	"github.com/boxgo/box/pkg/util/testutil"
)

func TestLock(t *testing.T) {
	testKey := strutil.RandomAlphabet(10)

	locked, err := Lock(testKey, time.Second*5)
	testutil.ExpectEqual(t, err, nil)
	testutil.ExpectEqual(t, locked, true)

	time.Sleep(time.Second * 3)

	locked, err = IsLocked(testKey)
	testutil.ExpectEqual(t, err, nil)
	testutil.ExpectEqual(t, locked, true)

	time.Sleep(time.Second * 2)
	locked, err = IsLocked(testKey)
	testutil.ExpectEqual(t, err, nil)
	testutil.ExpectEqual(t, locked, false)
}

func TestUnlock(t *testing.T) {
	testKey := strutil.RandomAlphabet(10)

	locked, err := Lock(testKey, time.Second*5)
	testutil.ExpectEqual(t, err, nil)
	testutil.ExpectEqual(t, locked, true)

	locked, err = IsLocked(testKey)
	testutil.ExpectEqual(t, err, nil)
	testutil.ExpectEqual(t, locked, true)

	err = UnLock(testKey)
	testutil.ExpectEqual(t, err, nil)

	locked, err = IsLocked(testKey)
	testutil.ExpectEqual(t, err, nil)
	testutil.ExpectEqual(t, locked, false)
}
