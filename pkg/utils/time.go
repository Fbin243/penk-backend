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
