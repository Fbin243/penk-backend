package coredb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type characterInputType = struct {
	Name   string
	Avatar string
	Gender bool
	Tags   []string
}

var charInput = &characterInputType{
	Name:   "example",
	Avatar: "https://example.com",
	Gender: true,
	Tags:   []string{"#tag1", "#tag2"},
}

func newCharacterFromInput(input *characterInputType) *Character {
	return &Character{
		ID:                  primitive.NewObjectID(),
		UserID:              primitive.NewObjectID(),
		Name:                input.Name,
		Avatar:              input.Avatar,
		Tags:                input.Tags,
		Gender:              input.Gender,
		TotalFocusedTime:    0,
		CustomMetrics:       []CustomMetric{},
		LimitedMetricNumber: 2,
	}
}

func assertWithCharInput(t *testing.T, character *Character, input *characterInputType) {
	assert.Equal(t, character.Name, input.Name)
	assert.Equal(t, character.Avatar, input.Avatar)
	assert.Equal(t, character.Gender, input.Gender)
	assert.Equal(t, character.Tags, input.Tags)
}

func TestCreateNewCharacter(t *testing.T) {
	character := newCharacterFromInput(charInput)

	createdCharacter, err := charactersRepo.CreateCharacter(character)
	assert.Nil(t, err)
	assertWithCharInput(t, createdCharacter, charInput)
}

func TestGetCharacterByID(t *testing.T) {
	character := newCharacterFromInput(charInput)

	_, err := charactersRepo.CreateCharacter(character)
	assert.Nil(t, err)

	queriedCharacter, err := charactersRepo.GetCharacterByID(character.ID)
	assert.Nil(t, err)
	assertWithCharInput(t, queriedCharacter, charInput)
}

func TestGetCharactersByUserID(t *testing.T) {
	character := newCharacterFromInput(charInput)

	_, err := charactersRepo.CreateCharacter(character)
	assert.Nil(t, err)

	characters, err := charactersRepo.GetCharactersByUserID(character.UserID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(characters))
	assert.Equal(t, *character, characters[0])
}

func TestUpdateCharacter(t *testing.T) {
	character := newCharacterFromInput(charInput)

	_, err := charactersRepo.CreateCharacter(character)
	assert.Nil(t, err)

	updateInput := &characterInputType{
		Name:   "updated",
		Avatar: "https://updated.com",
		Gender: false,
		Tags:   []string{"#updatedTag1"},
	}

	character.Name = updateInput.Name
	character.Gender = updateInput.Gender
	character.Avatar = updateInput.Avatar
	character.Tags = updateInput.Tags

	updatedCharacter, err := charactersRepo.UpdateCharacter(character)
	assert.Nil(t, err)
	assertWithCharInput(t, updatedCharacter, updateInput)
}
