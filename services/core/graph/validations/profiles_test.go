package validations_test

import (
	"testing"

	"tenkhours/services/core/graph/model"
	"tenkhours/services/core/graph/validations"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestValidateProfile(t *testing.T) {
	type testCase struct {
		name     string
		profile  model.ProfileInput
		hasError bool
	}

	tests := []testCase{
		{
			name: "valid profile",
			profile: model.ProfileInput{
				Name:               "John Doe",
				ImageURL:           "http://example.com/image.png",
				CurrentCharacterID: lo.ToPtr(primitive.NewObjectID()),
			},
			hasError: false,
		},
		{
			name: "empty name",
			profile: model.ProfileInput{
				Name:               "",
				ImageURL:           "http://example.com/image.png",
				CurrentCharacterID: lo.ToPtr(primitive.NewObjectID()),
			},
			hasError: true,
		},
		{
			name: "name too long",
			profile: model.ProfileInput{
				Name:               "This is a very long name that exceeds the maximum allowed length of fifty characters. This is a very long name that exceeds the maximum allowed length of fifty characters.",
				ImageURL:           "http://example.com/image.png",
				CurrentCharacterID: lo.ToPtr(primitive.NewObjectID()),
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
