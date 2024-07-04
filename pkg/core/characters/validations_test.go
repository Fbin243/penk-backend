package characters

import (
	"testing"

	"tenkhours/pkg/db/coredb"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestValidateCharacter(t *testing.T) {
	type testCase struct {
		name      string
		character coredb.Character
		hasError  bool
	}

	tests := []testCase{
		{
			name: "valid character",
			character: coredb.Character{
				ID:                  primitive.NewObjectID(),
				UserID:              primitive.NewObjectID(),
				Name:                "Hero",
				Gender:              true,
				Tags:                []string{"warrior", "legend"},
				TotalFocusedTime:    100,
				CustomMetrics:       []coredb.CustomMetric{},
				LimitedMetricNumber: 5,
			},
			hasError: false,
		},
		{
			name: "invalid character (empty name)",
			character: coredb.Character{
				ID:                  primitive.NewObjectID(),
				UserID:              primitive.NewObjectID(),
				Name:                "",
				Gender:              true,
				Tags:                []string{"warrior", "legend"},
				TotalFocusedTime:    100,
				CustomMetrics:       []coredb.CustomMetric{},
				LimitedMetricNumber: 5,
			},
			hasError: true,
		},
		{
			name: "invalid character (name too long)",
			character: coredb.Character{
				ID:                  primitive.NewObjectID(),
				UserID:              primitive.NewObjectID(),
				Name:                "This is a very long name that exceeds the maximum allowed length of fifty characters",
				Gender:              true,
				Tags:                []string{"warrior", "legend"},
				TotalFocusedTime:    100,
				CustomMetrics:       []coredb.CustomMetric{},
				LimitedMetricNumber: 5,
			},
			hasError: true,
		},
		{
			name: "invalid character (empty tag)",
			character: coredb.Character{
				ID:                  primitive.NewObjectID(),
				UserID:              primitive.NewObjectID(),
				Name:                "Hero",
				Gender:              true,
				Tags:                []string{"warrior", ""},
				TotalFocusedTime:    100,
				CustomMetrics:       []coredb.CustomMetric{},
				LimitedMetricNumber: 5,
			},
			hasError: true,
		},
		{
			name: "invalid character (empty custom metric name)",
			character: coredb.Character{
				ID:                  primitive.NewObjectID(),
				UserID:              primitive.NewObjectID(),
				Name:                "Hero",
				Gender:              true,
				Tags:                []string{"warrior", "legend"},
				TotalFocusedTime:    100,
				CustomMetrics:       []coredb.CustomMetric{{Name: "", Description: "Power level", Time: 10}},
				LimitedMetricNumber: 5,
			},
			hasError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateCharacter(tc.character)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
