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

var metric = entity.CustomMetric{
	ID:          primitive.NewObjectID().Hex(),
	Name:        "Metric name",
	Description: "Metric desc",
	Style: entity.MetricStyle{
		Color: "red",
		Icon:  "icon.png",
	},
	Properties: []entity.MetricProperty{
		{
			ID:    primitive.NewObjectID().Hex(),
			Name:  "Property name",
			Type:  "int",
			Value: "10",
			Unit:  "units",
		},
	},
}

func NewCharacter() *entity.Character {
	return &entity.Character{
		BaseEntity: &base.BaseEntity{
			ID:        primitive.NewObjectID().Hex(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ProfileID:        primitive.NewObjectID().Hex(),
		Name:             "example",
		Tags:             []string{"#tag1", "#tag2"},
		Gender:           false,
		TotalFocusedTime: 0,
		CustomMetrics: []entity.CustomMetric{
			metric, metric, metric,
		},
		LimitedMetricNumber: 2,
		Vision: entity.Vision{
			Name:        "Vision name",
			Description: "Vision desc",
		},
	}
}

func assertCharacter(t *testing.T, expected, actual *entity.Character) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.ProfileID, actual.ProfileID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Tags, actual.Tags)
	assert.Equal(t, expected.Gender, actual.Gender)
	assert.Equal(t, expected.TotalFocusedTime, actual.TotalFocusedTime)
	assert.Equal(t, expected.CustomMetrics, actual.CustomMetrics)
	assert.Equal(t, expected.LimitedMetricNumber, actual.LimitedMetricNumber)
	assert.Equal(t, expected.Vision, actual.Vision)
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

func TestGetAllCharacters(t *testing.T) {
	characterMap := map[string]*entity.Character{}
	for i := 0; i < 3; i++ {
		character := NewCharacter()
		characterMap[character.ID] = character
		createdCharacter, err := characterRepo.InsertOne(context.Background(), character)
		defer cleanUpCharacter(createdCharacter.ID)
		assert.Nil(t, err)
	}

	retrievedCharacters, err := characterRepo.GetAllCharacters(context.Background())
	assert.Nil(t, err)

	assert.Equal(t, 3, len(retrievedCharacters))
	for i := 0; i < 3; i++ {
		assertCharacter(t, characterMap[retrievedCharacters[i].ID], &retrievedCharacters[i])
	}
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
	character.Gender = updateInput["gender"].(bool)
	character.Tags = updateInput["tags"].([]string)

	updatedCharacter, err := characterRepo.UpdateByID(context.Background(), character.ID, character)
	assert.Nil(t, err)
	assertCharacter(t, character, updatedCharacter)
}

func cleanUpCharacter(id string) {
	_, err := characterRepo.DeleteCharacter(context.Background(), id)
	if err != nil {
		panic(err)
	}
}
