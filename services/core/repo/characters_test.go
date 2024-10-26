package repo_test

import (
	"context"
	"testing"

	"tenkhours/services/core/repo"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type characterInputType struct {
	Name   string
	Gender bool
	Tags   []string
}

type styleInputType struct {
	Color string
	Icon  string
}

type metricInputType struct {
	Name        string
	Description string
	Style       styleInputType
	Properties  []repo.MetricProperty
}

var charInput = &characterInputType{
	Name:   "example",
	Gender: true,
	Tags:   []string{"#tag1", "#tag2"},
}

func newCharacterFromInput(input *characterInputType) *repo.Character {
	return &repo.Character{
		ID:                  primitive.NewObjectID(),
		ProfileID:           primitive.NewObjectID(),
		Name:                input.Name,
		Tags:                input.Tags,
		Gender:              input.Gender,
		TotalFocusedTime:    0,
		CustomMetrics:       []repo.CustomMetric{},
		LimitedMetricNumber: 2,
	}
}

func assertWithCharInput(t *testing.T, character *repo.Character, input *characterInputType) {
	assert.Equal(t, character.Name, input.Name)
	assert.Equal(t, character.Gender, input.Gender)
	assert.Equal(t, character.Tags, input.Tags)
}

var metricInput1 = &metricInputType{
	Name:        "Metric 1 example",
	Description: "metric 1",
	Style: styleInputType{
		Color: "red",
		Icon:  "1",
	},
	Properties: []repo.MetricProperty{
		{
			ID:    primitive.NewObjectID(),
			Name:  "Property 1",
			Type:  "int",
			Value: "10",
			Unit:  "units",
		},
	},
}

var metricInput2 = &metricInputType{
	Name:        "Metric 2 example",
	Description: "metric 2",
	Style: styleInputType{
		Color: "blue",
		Icon:  "2",
	},
	Properties: []repo.MetricProperty{
		{
			ID:    primitive.NewObjectID(),
			Name:  "Property 2",
			Type:  "float",
			Value: "5.5",
			Unit:  "units",
		},
	},
}

var metricInput3 = &metricInputType{
	Name:        "Metric 3 example",
	Description: "metric 3",
	Style: styleInputType{
		Color: "yellow",
		Icon:  "3",
	},
	Properties: []repo.MetricProperty{
		{
			ID:    primitive.NewObjectID(),
			Name:  "Property 3",
			Type:  "string",
			Value: "value",
			Unit:  "",
		},
	},
}

func newMetricFromInput(input *metricInputType) *repo.CustomMetric {
	return &repo.CustomMetric{
		ID:                    primitive.NewObjectID(),
		Name:                  input.Name,
		Description:           input.Description,
		Style:                 repo.MetricStyle(input.Style),
		Properties:            input.Properties,
		LimitedPropertyNumber: 2,
	}
}

func assertWithMetricInput(t *testing.T, metric *repo.CustomMetric, input *metricInputType) {
	assert.Equal(t, metric.Name, input.Name)
	assert.Equal(t, metric.Description, input.Description)
	assert.Equal(t, metric.Style, repo.MetricStyle(input.Style))
	assert.Equal(t, metric.Properties, input.Properties)
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

func TestGetCharactersByProfileID(t *testing.T) {
	character := newCharacterFromInput(charInput)

	_, err := charactersRepo.CreateCharacter(character)
	assert.Nil(t, err)

	characters, err := charactersRepo.GetCharactersByProfileID(character.ProfileID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(characters))
	assert.Equal(t, *character, characters[0])
}

func setupMultipleCharactersTest(t *testing.T) ([]*repo.Character, func()) {
	_, err := charactersRepo.Collection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		t.Fatalf("Failed to clean up characters collection: %v", err)
	}

	character1 := newCharacterFromInput(charInput)
	character2 := newCharacterFromInput(&characterInputType{
		Name:   "example2",
		Gender: false,
		Tags:   []string{"#tag3"},
	})

	_, err = charactersRepo.CreateCharacter(character1)
	if err != nil {
		t.Fatalf("failed to create character 1: %v", err)
	}

	_, err = charactersRepo.CreateCharacter(character2)
	if err != nil {
		t.Fatalf("failed to create character 2: %v", err)
	}

	cleanup := func() {
		_, err := charactersRepo.DeleteCharacter(character1.ID)
		if err != nil {
			t.Fatalf("failed to delete character 1: %v", err)
		}

		_, err = charactersRepo.DeleteCharacter(character2.ID)
		if err != nil {
			t.Fatalf("failed to delete character 2: %v", err)
		}
	}

	return []*repo.Character{character1, character2}, cleanup
}

func TestGetAllCharacters(t *testing.T) {
	characters, cleanup := setupMultipleCharactersTest(t)
	defer cleanup()

	retrievedCharacters, err := charactersRepo.GetAllCharacters()
	assert.Nil(t, err)

	assert.Equal(t, len(characters), len(retrievedCharacters))

	for _, char := range characters {
		assert.Contains(t, retrievedCharacters, *char)
	}
}

func TestUpdateCharacter(t *testing.T) {
	character := newCharacterFromInput(charInput)

	_, err := charactersRepo.CreateCharacter(character)
	assert.Nil(t, err)

	updateInput := &characterInputType{
		Name:   "updated",
		Gender: false,
		Tags:   []string{"#updatedTag1"},
	}

	character.Name = updateInput.Name
	character.Gender = updateInput.Gender
	character.Tags = updateInput.Tags

	updatedCharacter, err := charactersRepo.UpdateCharacter(character)
	assert.Nil(t, err)
	assertWithCharInput(t, updatedCharacter, updateInput)
}

func TestCreateCustomMetric(t *testing.T) {
	character := newCharacterFromInput(charInput)
	_, err := charactersRepo.CreateCharacter(character)
	assert.Nil(t, err)

	metric := newMetricFromInput(metricInput1)
	createdMetric, err := charactersRepo.CreateCustomMetric(character.ID, metric)
	assert.Nil(t, err)
	assertWithMetricInput(t, createdMetric, metricInput1)
}

func TestUpdateCustomMetric(t *testing.T) {
	character := newCharacterFromInput(charInput)
	_, err := charactersRepo.CreateCharacter(character)
	assert.Nil(t, err)

	metric := newMetricFromInput(metricInput1)
	_, err = charactersRepo.CreateCustomMetric(character.ID, metric)
	assert.Nil(t, err)

	updateInput := &metricInputType{
		Name:        "Updated Metric",
		Description: "Updated description",
		Style: styleInputType{
			Color: "green",
			Icon:  "updatedIcon",
		},
		Properties: []repo.MetricProperty{
			{
				ID:    primitive.NewObjectID(),
				Name:  "Updated Property",
				Type:  "int",
				Value: "20",
				Unit:  "updatedUnits",
			},
		},
	}

	metric.Name = updateInput.Name
	metric.Description = updateInput.Description
	metric.Style = repo.MetricStyle(updateInput.Style)
	metric.Properties = updateInput.Properties

	updatedMetric, err := charactersRepo.UpdateCustomMetric(character.ID, metric)
	assert.Nil(t, err)
	assertWithMetricInput(t, updatedMetric, updateInput)
}

func TestAddMetricProperty(t *testing.T) {
	character := newCharacterFromInput(charInput)
	_, err := charactersRepo.CreateCharacter(character)
	assert.Nil(t, err)

	metric := newMetricFromInput(metricInput1)
	_, err = charactersRepo.CreateCustomMetric(character.ID, metric)
	assert.Nil(t, err)

	property := repo.MetricProperty{
		ID:    primitive.NewObjectID(),
		Name:  "New Property",
		Type:  "int",
		Value: "200",
		Unit:  "units",
	}

	createdProperty, err := charactersRepo.CreateMetricProperty(character.ID, metric.ID, &property)
	assert.Nil(t, err)
	assert.Equal(t, createdProperty, &property)
}

func TestUpdateMetricProperty(t *testing.T) {
	character := newCharacterFromInput(charInput)
	_, err := charactersRepo.CreateCharacter(character)
	assert.Nil(t, err)

	metric := newMetricFromInput(metricInput1)
	_, err = charactersRepo.CreateCustomMetric(character.ID, metric)
	assert.Nil(t, err)

	property := repo.MetricProperty{
		ID:    primitive.NewObjectID(),
		Name:  "Property 1",
		Type:  "int",
		Value: "10",
		Unit:  "units",
	}

	_, err = charactersRepo.CreateMetricProperty(character.ID, metric.ID, &property)
	assert.Nil(t, err)

	updateProperty := repo.MetricProperty{
		ID:    property.ID,
		Name:  "Updated Property",
		Type:  "float",
		Value: "15.5",
		Unit:  "updatedUnits",
	}

	updatedProperty, err := charactersRepo.UpdateMetricProperty(character.ID, metric.ID, &updateProperty)
	assert.Nil(t, err)
	assert.Equal(t, updatedProperty, &updateProperty)
}
