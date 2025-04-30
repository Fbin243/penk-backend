package utils

import (
	"slices"
	"time"

	"github.com/teambition/rrule-go"
)

// FindTimestamp checks if a given timestamp is present in the RRule's occurrences.
func FindTimestamp(rule *rrule.RRule, timestamp time.Time) (int, bool) {
	if rule == nil {
		return -1, false
	}

	// Set the end date of RRule to tommorrow
	rule.Until(timestamp.AddDate(0, 0, 1))

	return slices.BinarySearchFunc(rule.All(), timestamp, func(occur, target time.Time) int {
		occurStr := occur.Format(time.DateOnly)
		targetStr := target.Format(time.DateOnly)
		if occurStr < targetStr {
			return -1
		}
		if occurStr > targetStr {
			return 1
		}
		return 0
	})
}
