package coredb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type characterInputType struct {
	Name   string
	Avatar string
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
	Properties  []MetricProperty
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

var metricInput1 = &metricInputType{
	Name:        "Metric 1 example",
	Description: "metric 1",
	Style: styleInputType{
		Color: "red",
		Icon:  "1",
	},
	Properties: []MetricProperty{
		{
			ID:    primitive.NewObjectID(),
			Name:  "Property 1",
			Type:  "int",
			Value: 10,
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
	Properties: []MetricProperty{
		{
			ID:    primitive.NewObjectID(),
			Name:  "Property 2",
			Type:  "float",
			Value: 5.5,
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
	Properties: []MetricProperty{
		{
			ID:    primitive.NewObjectID(),
			Name:  "Property 3",
			Type:  "string",
			Value: "value",
			Unit:  "",
		},
	},
}

func newMetricFromInput(input *metricInputType) *CustomMetric {
	return &CustomMetric{
		ID:                    primitive.NewObjectID(),
		Name:                  input.Name,
		Description:           input.Description,
		Style:                 MetricStyle(input.Style),
		Properties:            input.Properties,
		LimitedPropertyNumber: 2,
	}
}

func assertWithMetricInput(t *testing.T, metric *CustomMetric, input *metricInputType) {
	assert.Equal(t, metric.Name, input.Name)
	assert.Equal(t, metric.Description, input.Description)
	assert.Equal(t, metric.Style, MetricStyle(input.Style))
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

func TestCreateCustomMetric(t *testing.T) {
	character := newCharacterFromInput(charInput)
	_, err := charactersRepo.CreateCharacter(character)
	assert.Nil(t, err)

	metric := newMetricFromInput(metricInput1)
	createdMetric, err := charactersRepo.CreateCustomMetric(character.ID, metric)
	assert.Nil(t, err)
	assertWithMetricInput(t, createdMetric, metricInput1)
}

func TestDeleteCustomMetric(t *testing.T) {
	character := newCharacterFromInput(charInput)
	_, err := charactersRepo.CreateCharacter(character)
	assert.Nil(t, err)

	metric := newMetricFromInput(metricInput1)
	_, err = charactersRepo.CreateCustomMetric(character.ID, metric)
	assert.Nil(t, err)

	deletedMetric, err := charactersRepo.DeleteCustomMetric(character.ID, metric.ID)
	assert.Nil(t, err)
	assertWithMetricInput(t, deletedMetric, metricInput1)
}

func TestAddMetricProperty(t *testing.T) {
	character := newCharacterFromInput(charInput)
	_, err := charactersRepo.CreateCharacter(character)
	assert.Nil(t, err)

	metric := newMetricFromInput(metricInput1)
	_, err = charactersRepo.CreateCustomMetric(character.ID, metric)
	assert.Nil(t, err)

	property := MetricProperty{
		ID:    primitive.NewObjectID(),
		Name:  "New Property",
		Type:  "int",
		Value: 100,
		Unit:  "units",
	}

	metric.Properties = append(metric.Properties, property)

	updatedMetric, err := charactersRepo.UpdateCharacter(character)
	assert.Nil(t, err)
	assert.Equal(t, updatedMetric.CustomMetrics[0].Properties, metric.Properties)
}

func TestRemoveMetricProperty(t *testing.T) {
	character := newCharacterFromInput(charInput)
	_, err := charactersRepo.CreateCharacter(character)
	assert.Nil(t, err)

	metric := newMetricFromInput(metricInput1)
	_, err = charactersRepo.CreateCustomMetric(character.ID, metric)
	assert.Nil(t, err)

	property := MetricProperty{
		ID:    primitive.NewObjectID(),
		Name:  "New Property",
		Type:  "int",
		Value: 100,
		Unit:  "units",
	}

	metric.Properties = append(metric.Properties, property)

	// Update the metric with the new property
	updatedMetric, err := charactersRepo.UpdateCharacter(character)
	assert.Nil(t, err)
	assert.Equal(t, updatedMetric.CustomMetrics[0].Properties, metric.Properties)

	// Now remove the property
	metric.Properties = metric.Properties[:len(metric.Properties)-1]

	updatedMetric, err = charactersRepo.UpdateCharacter(character)
	assert.Nil(t, err)
	assert.Equal(t, updatedMetric.CustomMetrics[0].Properties, metric.Properties)
}
