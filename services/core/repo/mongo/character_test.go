package mongorepo_test

import (
	"context"
	"testing"
	"time"

	"tenkhours/pkg/db/base"
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"github.com/stretchr/testify/assert"
)

func NewCharacter() *entity.Character {
	return &entity.Character{
		BaseEntity: &base.BaseEntity{
			ID:        mongodb.GenObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ProfileID: mongodb.GenObjectID(),
		Name:      "Example",
	}
}

func assertCharacter(t *testing.T, expected, actual *entity.Character) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.ProfileID, actual.ProfileID)
	assert.Equal(t, expected.Name, actual.Name)
}

func TestCreateNewCharacter(t *testing.T) {
	character := NewCharacter()
	createdCharacter, err := characterRepo.InsertOne(context.Background(), character)
	defer cleanUpCharacter(createdCharacter.ID)

	assert.Nil(t, err)
	assertCharacter(t, character, createdCharacter)
}

func TestGetCharacterByID(t *testing.T) {
	character := NewCharacter()
	_, err := characterRepo.InsertOne(context.Background(), character)
	defer cleanUpCharacter(character.ID)
	assert.Nil(t, err)

	queriedCharacter, err := characterRepo.FindByID(context.Background(), character.ID)
	assert.Nil(t, err)
	assertCharacter(t, character, queriedCharacter)
}

func TestGetCharactersByProfileID(t *testing.T) {
	character := NewCharacter()
	_, err := characterRepo.InsertOne(context.Background(), character)
	defer cleanUpCharacter(character.ID)
	assert.Nil(t, err)

	characters, err := characterRepo.GetCharactersByProfileID(context.Background(), character.ProfileID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(characters))
	assertCharacter(t, character, &characters[0])
}

func TestUpdateCharacter(t *testing.T) {
	character := NewCharacter()
	_, err := characterRepo.InsertOne(context.Background(), character)
	defer cleanUpCharacter(character.ID)
	assert.Nil(t, err)

	updateInput := map[string]interface{}{
		"name":   "updated",
		"gender": false,
		"tags":   []string{"#updatedTag1"},
	}

	character.Name = updateInput["name"].(string)

	updatedCharacter, err := characterRepo.FindAndUpdateByID(context.Background(), character.ID, character)
	assert.Nil(t, err)
	assertCharacter(t, character, updatedCharacter)
}

func cleanUpCharacter(id string) {
	_, err := characterRepo.DeleteCharacter(context.Background(), id)
	if err != nil {
		panic(err)
	}
}
