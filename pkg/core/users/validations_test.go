package users

import (
	"testing"
	"time"

	"tenkhours/pkg/db/coredb"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestValidateUser(t *testing.T) {
	type testCase struct {
		name     string
		user     coredb.User
		hasError bool
	}

	tests := []testCase{
		{
			name: "valid user",
			user: coredb.User{
				ID:                 primitive.NewObjectID(),
				Name:               "John Doe",
				Email:              "john@example.com",
				FirebaseUID:        "someFirebaseUID",
				ImageURL:           "http://example.com/image.png",
				CurrentCharacterID: primitive.NewObjectID(),
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			},
			hasError: false,
		},
		{
			name: "invalid user (empty name)",
			user: coredb.User{
				ID:                 primitive.NewObjectID(),
				Name:               "",
				Email:              "john@example.com",
				FirebaseUID:        "someFirebaseUID",
				ImageURL:           "http://example.com/image.png",
				CurrentCharacterID: primitive.NewObjectID(),
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			},
			hasError: true,
		},
		{
			name: "invalid user (name too long)",
			user: coredb.User{
				ID:                 primitive.NewObjectID(),
				Name:               "This is a very long name that exceeds the maximum allowed length of fifty characters",
				Email:              "john@example.com",
				FirebaseUID:        "someFirebaseUID",
				ImageURL:           "http://example.com/image.png",
				CurrentCharacterID: primitive.NewObjectID(),
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			},
			hasError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateUser(tc.user)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
