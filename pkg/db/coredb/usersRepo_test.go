package coredb

import (
	"testing"

	"tenkhours/pkg/db"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	testdb    = db.InitDBFromURL("mongodb://localhost:27017", "test")
	usersRepo = NewUsersRepo(testdb)
)

func TestGetUserByFirebaseUID(t *testing.T) {
	user := &User{
		ID:          primitive.NewObjectID(),
		FirebaseUID: primitive.NewObjectID().Hex(),
		Email:       primitive.NewObjectID().Hex() + "@gmail.com",
	}

	_, err := usersRepo.CreateNewUser(user)
	if err != nil {
		t.Logf("failed to create user: %v\n", err)
	}

	queriedUser, err := usersRepo.GetUserByFirebaseUID(user.FirebaseUID)
	assert.Nil(t, err)
	assert.Equal(t, *user, *queriedUser)
}

func TestCreateNewUser(t *testing.T) {
	user := &User{
		ID:          primitive.NewObjectID(),
		Name:        "example",
		Email:       primitive.NewObjectID().Hex() + "@gmail.com",
		FirebaseUID: primitive.NewObjectID().Hex(),
	}

	createdUser, err := usersRepo.CreateNewUser(user)
	assert.Nil(t, err)
	assert.Equal(t, *user, *createdUser)
}

func TestCreateSameUser(t *testing.T) {
	user := &User{
		ID:          primitive.NewObjectID(),
		Name:        "example",
		Email:       primitive.NewObjectID().Hex() + "@gmail.com",
		FirebaseUID: primitive.NewObjectID().Hex(),
	}

	_, err := usersRepo.CreateNewUser(user)
	assert.Nil(t, err)

	_, err = usersRepo.CreateNewUser(user)
	assert.NotNil(t, err)
}

func TestUpdateUser(t *testing.T) {
	user := &User{
		ID:                 primitive.NewObjectID(),
		Name:               "example",
		Email:              primitive.NewObjectID().Hex() + "@gmail.com",
		FirebaseUID:        primitive.NewObjectID().Hex(),
		ImageURL:           "https://example.com",
		CurrentCharacterID: primitive.NewObjectID(),
	}

	_, err := usersRepo.CreateNewUser(user)
	assert.Nil(t, err)

	newCurrentCharacterID := primitive.NewObjectID()
	user.Name = "updated"
	user.ImageURL = "https://updated.com"
	user.CurrentCharacterID = newCurrentCharacterID

	updatedUser, err := usersRepo.UpdateUserByID(user.ID, user)
	assert.Nil(t, err)
	assert.Equal(t, updatedUser.Name, "updated")
	assert.Equal(t, updatedUser.ImageURL, "https://updated.com")
	assert.Equal(t, updatedUser.CurrentCharacterID, newCurrentCharacterID)
}
