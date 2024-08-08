package users

import (
	"testing"

	"tenkhours/pkg/db/coredb"
	"tenkhours/pkg/utils"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestValidateUser(t *testing.T) {
	type testCase struct {
		name     string
		user     coredb.User
		hasError bool
	}

	user := coredb.User{
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
			name:     "valid user",
			user:     user,
			hasError: false,
		},
		{
			name: "empty name",
			user: func(u coredb.User) coredb.User {
				u.Name = ""
				return u
			}(user),
			hasError: true,
		},
		{
			name: "name too long",
			user: func(u coredb.User) coredb.User {
				u.Name = "This is a very long name that exceeds the maximum allowed length of fifty characters"
				return u
			}(user),
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
