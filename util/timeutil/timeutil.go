package timeutil

import (
	"time"
)

// BeginOfMinute begin of minute
func BeginOfMinute(now time.Time) time.Time {
	return now.Truncate(time.Minute)
}

// EndOfMinute end of minute
func EndOfMinute(now time.Time) time.Time {
	return BeginOfMinute(now).Add(time.Minute - time.Nanosecond)
}

// BeginOfHour begin of hour
func BeginOfHour(now time.Time) time.Time {
	y, m, d := now.Date()
	return time.Date(y, m, d, now.Hour(), 0, 0, 0, now.Location())
}

// EndOfHour end of hour
func EndOfHour(now time.Time) time.Time {
	return BeginOfHour(now).Add(time.Hour - time.Nanosecond)
}

// BeginOfDay 获取一天的起始时间
func BeginOfDay(now time.Time) time.Time {
	y, m, d := now.Date()

	return time.Date(y, m, d, 0, 0, 0, 0, now.Location())
}

// EndOfDay 获取一天的结束时间
func EndOfDay(now time.Time) time.Time {
	y, m, d := now.Date()

	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), now.Location())
}

// BeginOfWeek 获取一周的起始时间
func BeginOfWeek(now time.Time) time.Time {
	begin := BeginOfDay(now)

	weekday := int(now.Weekday())

	return begin.AddDate(0, 0, -weekday+1) // 周一是开始
}

// EndOfWeek 获取一周的结束时间
func EndOfWeek(now time.Time) time.Time {
	return BeginOfWeek(now).AddDate(0, 0, 7).Add(-time.Nanosecond)
}

// BeginOfMonth 获取月的起始时间
func BeginOfMonth(now time.Time) time.Time {
	y, m, _ := now.Date()

	return time.Date(y, m, 1, 0, 0, 0, 0, now.Location())
}

// EndOfMonth 获取月的结束时间
func EndOfMonth(now time.Time) time.Time {
	y, m, _ := now.Date()

	return time.Date(y, m+1, 0, 23, 59, 59, int(time.Second-time.Nanosecond), now.Location())
}

// BeginOfYear beginning of year
func BeginOfYear(now time.Time) time.Time {
	y, _, _ := now.Date()
	return time.Date(y, time.January, 1, 0, 0, 0, 0, now.Location())
}

// EndOfYear end of year
func EndOfYear(now time.Time) time.Time {
	return BeginOfYear(now).AddDate(1, 0, 0).Add(-time.Nanosecond)
}
