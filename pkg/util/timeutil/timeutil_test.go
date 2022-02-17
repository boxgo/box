package timeutil_test

import (
	"testing"
	"time"

	"github.com/boxgo/box/pkg/util/timeutil"
)

func TestName(t *testing.T) {
	now := time.Now()

	t.Log(timeutil.BeginOfMinute(now))
	t.Log(timeutil.EndOfMinute(now))

	t.Log(timeutil.BeginOfHour(now))
	t.Log(timeutil.EndOfHour(now))

	t.Log(timeutil.BeginOfDay(now))
	t.Log(timeutil.EndOfDay(now))

	t.Log(timeutil.BeginOfWeek(now))
	t.Log(timeutil.EndOfWeek(now))

	t.Log(timeutil.BeginOfMonth(now))
	t.Log(timeutil.EndOfMonth(now))

	t.Log(timeutil.BeginOfYear(now))
	t.Log(timeutil.EndOfYear(now))
}
