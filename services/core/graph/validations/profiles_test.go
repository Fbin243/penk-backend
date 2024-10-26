package validations_test

import (
	"testing"

	"tenkhours/services/core/graph/validations"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestValidateProfile(t *testing.T) {
	type testCase struct {
		name     string
		profile  ProfileInput
		hasError bool
	}

	profile := NewProfileInput().Name("John Doe").ImageURL("http://example.com/image.png").CurrentCharacterID(primitive.NewObjectID())

	tests := []testCase{
		{
			name:     "valid profile",
			profile:  profile,
			hasError: false,
		},
		{
			name:     "empty name",
			profile:  profile.Name(""),
			hasError: true,
		},
		{
			name:     "name too long",
			profile:  profile.Name("This is a very long name that exceeds the maximum allowed length of fifty characters"),
			hasError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validations.ValidateProfileInput(tc.profile.ProfileInput)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
