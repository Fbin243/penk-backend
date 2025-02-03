package mongorepo_test

import (
	"context"
	"testing"
	"time"

	"tenkhours/pkg/db/base"
	"tenkhours/services/core/entity"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewProfile() *entity.Profile {
	return &entity.Profile{
		BaseEntity: &base.BaseEntity{
			ID:        primitive.NewObjectID().Hex(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:               "example",
		ImageURL:           "https://example.com",
		CurrentCharacterID: primitive.NewObjectID().Hex(),
		Email:              primitive.NewObjectID().Hex() + "@gmail.com",
		FirebaseUID:        primitive.NewObjectID().Hex(),
	}
}

func assertProfile(t *testing.T, expected, actual *entity.Profile) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.ImageURL, actual.ImageURL)
	assert.Equal(t, expected.CurrentCharacterID, actual.CurrentCharacterID)
	assert.Equal(t, expected.Email, actual.Email)
	assert.Equal(t, expected.FirebaseUID, actual.FirebaseUID)
}

func TestCreateNewProfile(t *testing.T) {
	profile := NewProfile()
	createdProfile, err := profileRepo.InsertOne(context.Background(), profile)
	defer cleanUpProfile(createdProfile.ID)

	assert.Nil(t, err)
	assertProfile(t, profile, createdProfile)
}

func TestCreateSameProfile(t *testing.T) {
	profile := NewProfile()
	_, err := profileRepo.InsertOne(context.Background(), profile)
	defer cleanUpProfile(profile.ID)
	assert.Nil(t, err)

	_, err = profileRepo.InsertOne(context.Background(), profile)
	assert.NotNil(t, err)
}

func TestGetProfileByFirebaseUID(t *testing.T) {
	profile := NewProfile()
	_, err := profileRepo.InsertOne(context.Background(), profile)
	defer cleanUpProfile(profile.ID)
	assert.Nil(t, err)

	queriedProfile, err := profileRepo.GetProfileByFirebaseUID(context.Background(), profile.FirebaseUID)
	assert.Nil(t, err)
	assertProfile(t, profile, queriedProfile)
}

func TestUpdateProfile(t *testing.T) {
	profile := NewProfile()
	_, err := profileRepo.InsertOne(context.Background(), profile)
	defer cleanUpProfile(profile.ID)
	assert.Nil(t, err)

	updateInput := map[string]any{
		"name":                 "updated",
		"image_url":            "https://updated.com",
		"current_character_id": primitive.NewObjectID().Hex(),
	}

	profile.Name = updateInput["name"].(string)
	profile.ImageURL = updateInput["image_url"].(string)
	profile.CurrentCharacterID = updateInput["current_character_id"].(string)

	updatedProfile, err := profileRepo.UpdateByID(context.Background(), profile.ID, profile)
	assert.Nil(t, err)
	assertProfile(t, profile, updatedProfile)
}

func cleanUpProfile(id string) {
	// Delete profile from database
	_, err := profileRepo.DeleteByID(context.Background(), id)
	if err != nil {
		panic(err)
	}
}
