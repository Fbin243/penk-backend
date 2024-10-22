package validations

import (
	"testing"

	"tenkhours/pkg/utils"
	"tenkhours/services/core/repo"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestValidateProfile(t *testing.T) {
	type testCase struct {
		name     string
		profile  repo.Profile
		hasError bool
	}

	profile := repo.Profile{
		ID:                 primitive.NewObjectID(),
		Name:               "John Doe",
		Email:              "john@example.com",
		FirebaseUID:        "someFirebaseUID",
		ImageURL:           "http://example.com/image.png",
		CurrentCharacterID: primitive.NewObjectID(),
		CreatedAt:          utils.Now(),
		UpdatedAt:          utils.Now(),
	}

	tests := []testCase{
		{
			name:     "valid profile",
			profile:  profile,
			hasError: false,
		},
		{
			name: "empty name",
			profile: func(u repo.Profile) repo.Profile {
				u.Name = ""
				return u
			}(profile),
			hasError: true,
		},
		{
			name: "name too long",
			profile: func(u repo.Profile) repo.Profile {
				u.Name = "This is a very long name that exceeds the maximum allowed length of fifty characters"
				return u
			}(profile),
			hasError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateProfile(tc.profile)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
