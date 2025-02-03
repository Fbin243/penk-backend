package mongorepo_test

import (
	"context"
	"testing"

	"tenkhours/pkg/db/base"
	"tenkhours/services/currency/entity"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type fishInputType struct {
	ProfileID string
	Gold      int32
	Normal    int32
}

var fishInput = &fishInputType{
	ProfileID: primitive.NewObjectID().Hex(),
	Gold:      10,
	Normal:    5,
}

func newFishFromInput(input *fishInputType) *entity.Fish {
	return &entity.Fish{
		BaseEntity: &base.BaseEntity{},
		ProfileID:  input.ProfileID,
		Gold:       input.Gold,
		Normal:     input.Normal,
	}
}

func cleanUp(profileID string) {
	_, err := fishRepo.DeleteFishByProfileID(context.Background(), profileID)
	if err != nil {
		panic(err)
	}
}

func assertWithFishInput(t *testing.T, fish *entity.Fish, input *fishInputType) {
	assert.Equal(t, fish.ProfileID, input.ProfileID)
	assert.Equal(t, fish.Gold, input.Gold)
	assert.Equal(t, fish.Normal, input.Normal)
}

func TestCreateFish(t *testing.T) {
	fish := newFishFromInput(fishInput)
	createdFish, err := fishRepo.InsertOne(context.Background(), fish)
	defer cleanUp(createdFish.ProfileID)

	assert.Nil(t, err)
	assertWithFishInput(t, createdFish, fishInput)
}

func TestGetFishByProfileID(t *testing.T) {
	fish := newFishFromInput(fishInput)
	createdFish, err := fishRepo.InsertOne(context.Background(), fish)
	assert.Nil(t, err)
	defer cleanUp(createdFish.ProfileID)

	retrievedFish, err := fishRepo.GetFishByProfileID(context.Background(), fish.ProfileID)
	assert.Nil(t, err)
	assertWithFishInput(t, retrievedFish, fishInput)
}

func TestUpdateFish(t *testing.T) {
	fish := newFishFromInput(fishInput)
	createdFish, err := fishRepo.InsertOne(context.Background(), fish)
	assert.Nil(t, err)
	defer cleanUp(createdFish.ProfileID)

	updateInput := &fishInputType{
		ProfileID: fish.ProfileID,
		Gold:      20,
		Normal:    15,
	}

	fish.Gold = updateInput.Gold
	fish.Normal = updateInput.Normal

	updatedFish, err := fishRepo.UpdateFishByProfileID(context.Background(), fish.ProfileID, fish)
	assert.Nil(t, err)

	assertWithFishInput(t, updatedFish, updateInput)
}
