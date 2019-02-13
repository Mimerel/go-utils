package go_utils

import (
	"time"
)

func Timestamp2BeginningOfMonth(moment time.Time) time.Time {
	y, m, _ := moment.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, time.Local)
}

func Timestamp2BeginningOfDay(moment time.Time) time.Time {
	y, m, d := moment.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

func Timestamp2BeginningOfYear(moment time.Time) time.Time {
	_, m, d := moment.Date()
	return time.Date(1, m, d, 0, 0, 0, 0, time.Local)
}

func Timestamp2EndOfMonth(moment time.Time) time.Time {
	return Timestamp2BeginningOfMonth(moment).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

func Timestamp2EndOfDay(moment time.Time) time.Time {
  return Timestamp2BeginningOfDay(moment).AddDate(0, 0, 1).Add(-time.Nanosecond)
}

func Timestamp2EndOfYear(moment time.Time) time.Time {
  return Timestamp2BeginningOfYear(moment).AddDate(1, 0, 0).Add(-time.Nanosecond)
}

