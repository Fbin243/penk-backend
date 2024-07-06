package coredb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userInputType struct {
	Name               string
	ImageURL           string
	CurrentCharacterID primitive.ObjectID
}

var userInput = &userInputType{
	Name:               "example",
	ImageURL:           "https://example.com",
	CurrentCharacterID: primitive.NewObjectID(),
}

func newUserFromInput(input *userInputType) *User {
	return &User{
		ID:                 primitive.NewObjectID(),
		Email:              primitive.NewObjectID().Hex() + "@gmail.com",
		FirebaseUID:        primitive.NewObjectID().Hex(),
		Name:               input.Name,
		ImageURL:           input.ImageURL,
		CurrentCharacterID: input.CurrentCharacterID,
	}
}

func assertWithUserInput(t *testing.T, user *User, input *userInputType) {
	assert.Equal(t, user.Name, input.Name)
	assert.Equal(t, user.ImageURL, input.ImageURL)
	assert.Equal(t, user.CurrentCharacterID, input.CurrentCharacterID)
}

func TestCreateNewUser(t *testing.T) {
	user := newUserFromInput(userInput)

	createdUser, err := usersRepo.CreateNewUser(user)
	assert.Nil(t, err)
	assertWithUserInput(t, createdUser, userInput)
}

func TestCreateSameUser(t *testing.T) {
	user := newUserFromInput(userInput)

	_, err := usersRepo.CreateNewUser(user)
	assert.Nil(t, err)

	_, err = usersRepo.CreateNewUser(user)
	assert.NotNil(t, err)
}

func TestGetUserByFirebaseUID(t *testing.T) {
	user := newUserFromInput(userInput)

	_, err := usersRepo.CreateNewUser(user)
	assert.Nil(t, err)

	queriedUser, err := usersRepo.GetUserByFirebaseUID(user.FirebaseUID)
	assert.Nil(t, err)
	assert.Equal(t, *queriedUser, *user)
}

func TestUpdateUser(t *testing.T) {
	user := newUserFromInput(userInput)

	_, err := usersRepo.CreateNewUser(user)
	assert.Nil(t, err)

	updateInput := &userInputType{
		Name:               "updated",
		ImageURL:           "https://updated.com",
		CurrentCharacterID: primitive.NewObjectID(),
	}

	user.Name = updateInput.Name
	user.ImageURL = updateInput.ImageURL
	user.CurrentCharacterID = updateInput.CurrentCharacterID

	updatedUser, err := usersRepo.UpdateUser(user)
	assert.Nil(t, err)
	assertWithUserInput(t, updatedUser, updateInput)
}
