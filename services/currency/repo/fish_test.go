package repo_test

import (
	"testing"

	"tenkhours/pkg/db"
	"tenkhours/services/currency/repo"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type fishInputType struct {
	ProfileID primitive.ObjectID
	Gold      int32
	Normal    int32
}

var fishInput = &fishInputType{
	ProfileID: primitive.NewObjectID(),
	Gold:      10,
	Normal:    5,
}

func newFishFromInput(input *fishInputType) *repo.Fish {
	return &repo.Fish{
		BaseModel: &db.BaseModel{},
		ProfileID: input.ProfileID,
		Gold:      input.Gold,
		Normal:    input.Normal,
	}
}

func cleanUp(profileID primitive.ObjectID) {
	_, err := fishRepo.DeleteFishByProfileID(profileID)
	if err != nil {
		panic(err)
	}
}

func assertWithFishInput(t *testing.T, fish *repo.Fish, input *fishInputType) {
	assert.Equal(t, fish.ProfileID, input.ProfileID)
	assert.Equal(t, fish.Gold, input.Gold)
	assert.Equal(t, fish.Normal, input.Normal)
}

func TestCreateFish(t *testing.T) {
	fish := newFishFromInput(fishInput)
	createdFish, err := fishRepo.InsertOne(fish)
	defer cleanUp(createdFish.ProfileID)

	assert.Nil(t, err)
	assertWithFishInput(t, createdFish, fishInput)
}

func TestGetFishByProfileID(t *testing.T) {
	fish := newFishFromInput(fishInput)
	createdFish, err := fishRepo.InsertOne(fish)
	defer cleanUp(createdFish.ProfileID)

	retrievedFish, err := fishRepo.GetFishByProfileID(fish.ProfileID)
	assert.Nil(t, err)
	assertWithFishInput(t, retrievedFish, fishInput)
}

func TestUpdateFish(t *testing.T) {
	fish := newFishFromInput(fishInput)
	createdFish, err := fishRepo.InsertOne(fish)
	defer cleanUp(createdFish.ProfileID)

	updateInput := &fishInputType{
		ProfileID: fish.ProfileID,
		Gold:      20,
		Normal:    15,
	}

	fish.Gold = updateInput.Gold
	fish.Normal = updateInput.Normal

	updatedFish, err := fishRepo.UpdateFishByProfileID(fish.ProfileID, fish)
	assert.Nil(t, err)

	assertWithFishInput(t, updatedFish, updateInput)
}
