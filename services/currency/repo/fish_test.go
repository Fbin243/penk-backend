package repo_test

// import (
// 	"context"
// 	"testing"

// 	"tenkhours/services/currency/repo"

// 	"github.com/stretchr/testify/assert"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type fishInputType struct {
// 	ProfileID primitive.ObjectID
// 	Numbers   int32
// 	Type      string
// }

// var fishInput = &fishInputType{
// 	ProfileID: primitive.NewObjectID(),
// 	Numbers:   3,
// 	Type:      "gold",
// }

// func newFishFromInput(input *fishInputType) *repo.Fish {
// 	return &repo.Fish{
// 		ID:        primitive.NewObjectID(),
// 		ProfileID: input.ProfileID,
// 		Numbers:   input.Numbers,
// 		Type:      input.Type,
// 	}
// }

// func assertWithFishInput(t *testing.T, fish *repo.Fish, input *fishInputType) {
// 	assert.Equal(t, fish.ProfileID, input.ProfileID)
// 	assert.Equal(t, fish.Numbers, input.Numbers)
// 	assert.Equal(t, fish.Type, input.Type)
// }

// func setupTest(t *testing.T) *repo.Fish {
// 	_, err := fishRepo.Collection.DeleteMany(context.Background(), bson.M{})
// 	if err != nil {
// 		t.Fatalf("Failed to clean up collection: %v", err)
// 	}

// 	fish := newFishFromInput(fishInput)

// 	_, err = fishRepo.CreateFish(fish)
// 	if err != nil {
// 		t.Fatalf("Failed to create fish: %v", err)
// 	}

// 	return fish
// }

// func TestCreateFish(t *testing.T) {
// 	fish := newFishFromInput(fishInput)

// 	createdFish, err := fishRepo.CreateFish(fish)
// 	assert.Nil(t, err)
// 	assertWithFishInput(t, createdFish, fishInput)
// }

// func TestGetFishByProfileID(t *testing.T) {
// 	fish := newFishFromInput(fishInput)

// 	_, err := fishRepo.CreateFish(fish)
// 	assert.Nil(t, err)

// 	queriedFish, err := fishRepo.GetFishByProfileID(fish.ProfileID, fish.Type)
// 	assert.Nil(t, err)

// 	queriedFish.ID = fish.ID
// 	assert.Equal(t, *queriedFish, *fish)
// }

// func TestUpdateFish(t *testing.T) {
// 	fish := setupTest(t)

// 	// Update the Numbers field
// 	updatedNumbers := fish.Numbers + 5
// 	fish.Numbers = updatedNumbers

// 	updatedFish, err := fishRepo.UpdateFish(fish, fish.ProfileID)
// 	assert.Nil(t, err)
// 	assert.Equal(t, updatedFish.Numbers, updatedNumbers)
// 	assert.Equal(t, updatedFish.Type, fish.Type)
// 	assert.Equal(t, updatedFish.ProfileID, fish.ProfileID)
// }
