package validations

import (
	"testing"

	"tenkhours/services/core/graph/model"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestValidateProfile(t *testing.T) {
	type testCase struct {
		name     string
		profile  model.ProfileInput
		hasError bool
	}

	profile := model.ProfileInput{
		Name:               toPtr("John Doe"),
		ImageURL:           toPtr("http://example.com/image.png"),
		CurrentCharacterID: toPtr(primitive.NewObjectID()),
	}

	tests := []testCase{
		{
			name:     "valid profile",
			profile:  profile,
			hasError: false,
		},
		{
			name: "empty name",
			profile: func(u model.ProfileInput) model.ProfileInput {
				u.Name = nil
				return u
			}(profile),
			hasError: true,
		},
		{
			name: "name too long",
			profile: func(u model.ProfileInput) model.ProfileInput {
				u.Name = toPtr("This is a very long name that exceeds the maximum allowed length of fifty characters")
				return u
			}(profile),
			hasError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateProfileInput(tc.profile)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
