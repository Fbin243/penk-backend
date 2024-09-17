package validations

import (
	"fmt"
	"time"
)

func ValidateDuration(startTime time.Time, duration int32) error {
	if duration >= 0 && duration <= int32(time.Since(startTime).Seconds()) {
		return nil
	}
	return fmt.Errorf("invalid duration")
}
