package go_utils

import (
	"strings"
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

func TimestampFromYYYYMMDD(moment string) (date time.Time, err error) {
  moment += "T12:00:00.000Z"
  layout := "2006-01-02T15:04:05.000Z"
  return time.ParseInLocation(layout, moment, time.Local)
}

func Timestamp2BeginningOfMonthUTC(moment time.Time) time.Time {
	y, m, _ := moment.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
}

func Timestamp2BeginningOfDayUTC(moment time.Time) time.Time {
	y, m, d := moment.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func Timestamp2BeginningOfYearUTC(moment time.Time) time.Time {
	_, m, d := moment.Date()
	return time.Date(1, m, d, 0, 0, 0, 0, time.UTC)
}

func Timestamp2EndOfMonthUTC(moment time.Time) time.Time {
	return Timestamp2BeginningOfMonthUTC(moment).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

func Timestamp2EndOfDayUTC(moment time.Time) time.Time {
  return Timestamp2BeginningOfDayUTC(moment).AddDate(0, 0, 1).Add(-time.Nanosecond)
}

func Timestamp2EndOfYearUTC(moment time.Time) time.Time {
  return Timestamp2BeginningOfYearUTC(moment).AddDate(1, 0, 0).Add(-time.Nanosecond)
}

func TimestampFromYYYYMMDDUTC(moment string) (date time.Time, err error) {
  moment += "T12:00:00.000Z"
  layout := "2006-01-02T15:04:05.000Z"
  return time.ParseInLocation(layout, moment, time.UTC)
}

func LocalTime() time.Time {
	return time.Now().In(Location())
}

func Location() *time.Location {
	loc, _ := time.LoadLocation("Europe/Paris")
	return loc
}

func localTimeString() string {
	value := LocalTime().Format(time.RFC3339)
	value = strings.Replace(value, "T", " ", -1)
	value = strings.Replace(value, "Z", "", -1)
	return value
}