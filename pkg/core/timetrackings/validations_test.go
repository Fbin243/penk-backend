package timetrackings

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateDuration(t *testing.T) {
	type testCase struct {
		name      string
		startTime time.Time
		duration  int32
		hasError  bool
	}

	tests := []testCase{
		{
			name:      "valid duration less than time since start",
			startTime: time.Now().Add(-time.Hour),
			duration:  60 * 30,
			hasError:  false,
		},
		{
			name:      "valid duration equal to time since start",
			startTime: time.Now().Add(-time.Hour),
			duration:  60 * 60,
			hasError:  false,
		},
		{
			name:      "negative duration",
			startTime: time.Now(),
			duration:  -10,
			hasError:  true,
		},
		{
			name:      "duration longer than time since start",
			startTime: time.Now().Add(-time.Hour),
			duration:  60*60 + 1,
			hasError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateDuration(tc.startTime, tc.duration)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
