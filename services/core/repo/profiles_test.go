package repo_test

import (
	"context"
	"testing"

	"tenkhours/services/core/repo"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type profileInputType struct {
	Name               string
	ImageURL           string
	CurrentCharacterID primitive.ObjectID
}

var profileInput = &profileInputType{
	Name:               "example",
	ImageURL:           "https://example.com",
	CurrentCharacterID: primitive.NewObjectID(),
}

func newProfileFromInput(input *profileInputType) *repo.Profile {
	return &repo.Profile{
		ID:                 primitive.NewObjectID(),
		Email:              primitive.NewObjectID().Hex() + "@gmail.com",
		FirebaseUID:        primitive.NewObjectID().Hex(),
		Name:               input.Name,
		ImageURL:           input.ImageURL,
		CurrentCharacterID: input.CurrentCharacterID,
	}
}

func assertWithProfileInput(t *testing.T, profile *repo.Profile, input *profileInputType) {
	assert.Equal(t, profile.Name, input.Name)
	assert.Equal(t, profile.ImageURL, input.ImageURL)
	assert.Equal(t, profile.CurrentCharacterID, input.CurrentCharacterID)
}

func TestCreateNewProfile(t *testing.T) {
	profile := newProfileFromInput(profileInput)

	createdProfile, err := profilesRepo.CreateNewProfile(profile)
	defer cleanUpProfile(createdProfile)

	assert.Nil(t, err)
	assertWithProfileInput(t, createdProfile, profileInput)
}

func TestCreateSameProfile(t *testing.T) {
	profile := newProfileFromInput(profileInput)

	_, err := profilesRepo.CreateNewProfile(profile)
	defer cleanUpProfile(profile)
	assert.Nil(t, err)

	_, err = profilesRepo.CreateNewProfile(profile)
	assert.NotNil(t, err)
}

func TestGetProfileByFirebaseUID(t *testing.T) {
	profile := newProfileFromInput(profileInput)

	_, err := profilesRepo.CreateNewProfile(profile)
	defer cleanUpProfile(profile)
	assert.Nil(t, err)

	queriedProfile, err := profilesRepo.GetProfileByFirebaseUID(profile.FirebaseUID)
	assert.Nil(t, err)
	assert.Equal(t, *queriedProfile, *profile)
}

func TestUpdateProfile(t *testing.T) {
	profile := newProfileFromInput(profileInput)

	_, err := profilesRepo.CreateNewProfile(profile)
	defer cleanUpProfile(profile)
	assert.Nil(t, err)

	updateInput := &profileInputType{
		Name:               "updated",
		ImageURL:           "https://updated.com",
		CurrentCharacterID: primitive.NewObjectID(),
	}

	profile.Name = updateInput.Name
	profile.ImageURL = updateInput.ImageURL
	profile.CurrentCharacterID = updateInput.CurrentCharacterID

	updatedProfile, err := profilesRepo.UpdateProfile(profile)
	assert.Nil(t, err)
	assertWithProfileInput(t, updatedProfile, updateInput)
}

func cleanUpProfile(profile *repo.Profile) {
	// Delete profile from database
	_, err := profilesRepo.Collection.DeleteOne(context.Background(), bson.M{"_id": profile.ID})
	if err != nil {
		panic(err)
	}

	// Delete profile from Redis
	_, err = profilesRepo.Del(context.Background(), profile.FirebaseUID).Result()
	if err != nil {
		panic(err)
	}
}
