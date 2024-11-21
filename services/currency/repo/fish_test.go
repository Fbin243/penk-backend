package repo_test

import (
	"context"
	"testing"

	"tenkhours/services/currency/repo"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupTestRepo(t *testing.T) *repo.FishRepo {
	// Assuming `fishRepo` is globally accessible and properly initialized
	_, err := fishRepo.Collection.DeleteMany(context.Background(), bson.M{})
	assert.Nil(t, err)
	return fishRepo
}

func newFish(profileID primitive.ObjectID, gold, normal int32) *repo.Fish {
	return &repo.Fish{
		ID:        primitive.NewObjectID(),
		ProfileID: profileID,
		Gold:      gold,
		Normal:    normal,
	}
}

func TestCreateFish(t *testing.T) {
	repo := setupTestRepo(t)
	fish := newFish(primitive.NewObjectID(), 10, 5)

	createdFish, err := repo.CreateFish(fish)
	assert.Nil(t, err)
	assert.Equal(t, fish.Gold, createdFish.Gold)
	assert.Equal(t, fish.Normal, createdFish.Normal)
}

func TestGetFishByProfileID(t *testing.T) {
	repo := setupTestRepo(t)
	profileID := primitive.NewObjectID()
	fish := newFish(profileID, 10, 5)

	_, err := repo.CreateFish(fish)
	assert.Nil(t, err)

	retrievedFish, err := repo.GetFishByProfileID(profileID)
	assert.Nil(t, err)
	assert.Equal(t, fish.Gold, retrievedFish.Gold)
	assert.Equal(t, fish.Normal, retrievedFish.Normal)
}

func TestUpdateFish(t *testing.T) {
	repo := setupTestRepo(t)
	profileID := primitive.NewObjectID()
	fish := newFish(profileID, 10, 5)

	_, err := repo.CreateFish(fish)
	assert.Nil(t, err)

	fish.Gold = 15
	fish.Normal = 8

	updatedFish, err := repo.UpdateFish(fish, profileID)
	assert.Nil(t, err)

	t.Logf("Updated Fish - Gold: %v, Normal: %v", updatedFish.Gold, updatedFish.Normal)

	assert.Equal(t, fish.Gold, updatedFish.Gold)
	assert.Equal(t, fish.Normal, updatedFish.Normal)
}
