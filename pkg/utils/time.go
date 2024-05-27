package utils

import "time"

func Now() time.Time {
	// Why truncate? 2024-03-03 04:28:37.994545 +0000 UTC --> 2024-03-03 04:28:37 +0000 UTC
	return time.Now().Truncate(time.Second).UTC()
}
