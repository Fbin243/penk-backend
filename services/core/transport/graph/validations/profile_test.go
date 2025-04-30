package validations_test

import (
	"testing"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"
	"tenkhours/services/core/transport/graph/validations"

	"github.com/stretchr/testify/assert"
)

func TestValidateProfile(t *testing.T) {
	type testCase struct {
		name     string
		profile  entity.ProfileInput
		hasError bool
	}

	tests := []testCase{
		{
			name: "valid profile",
			profile: entity.ProfileInput{
				Name:               "John Doe",
				ImageURL:           "http://example.com/image.png",
				CurrentCharacterID: mongodb.GenObjectID(),
			},
			hasError: false,
		},
		{
			name: "empty name",
			profile: entity.ProfileInput{
				Name:               "",
				ImageURL:           "http://example.com/image.png",
				CurrentCharacterID: mongodb.GenObjectID(),
			},
			hasError: true,
		},
		{
			name: "name too long",
			profile: entity.ProfileInput{
				Name:               "This is a very long name that exceeds the maximum allowed length of fifty characters. This is a very long name that exceeds the maximum allowed length of fifty characters.",
				ImageURL:           "http://example.com/image.png",
				CurrentCharacterID: mongodb.GenObjectID(),
			},
			hasError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validations.ValidateProfileInput(tc.profile)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
