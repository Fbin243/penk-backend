package utils

import (
	"fmt"
	"time"
)

func Now() time.Time {
	// Why truncate? 2024-03-03 04:28:37.994545 +0000 UTC --> 2024-03-03 04:28:37 +0000 UTC
	return time.Now().Truncate(time.Second).UTC()
}

// MonthToIntMap maps month names to their respective integer values
var MonthToIntMap = map[string]int{
	"JANUARY":   1,
	"FEBRUARY":  2,
	"MARCH":     3,
	"APRIL":     4,
	"MAY":       5,
	"JUNE":      6,
	"JULY":      7,
	"AUGUST":    8,
	"SEPTEMBER": 9,
	"OCTOBER":   10,
	"NOVEMBER":  11,
	"DECEMBER":  12,
}

func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, time.UTC)
}

func ParseTime(timeStr string) time.Time {
	time, _ := time.Parse(time.RFC3339, timeStr)
	return time
}

func UnixToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

func PrintTimeSlice(ts []time.Time) {
	for _, t := range ts {
		fmt.Println(t)
	}
}

func CompareTwoTimeWithoutSecond(t1, t2 time.Time) bool {
	// Compare year, month, day, hour, minute
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day() && t1.Hour() == t2.Hour() && t1.Minute() == t2.Minute()
}

func IsToday(t time.Time) bool {
	return CompareTwoTimeWithoutSecond(t, Now())
}
