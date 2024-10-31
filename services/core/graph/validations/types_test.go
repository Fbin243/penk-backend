package validations_test

import (
	"tenkhours/services/core/graph/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProfileInput struct {
	model.ProfileInput
}

func NewProfileInput() ProfileInput {
	return ProfileInput{model.ProfileInput{}}
}

func (p ProfileInput) Name(name string) ProfileInput {
	p.ProfileInput.Name = toPtr(name)
	return p
}

func (p ProfileInput) ImageURL(imageURL string) ProfileInput {
	p.ProfileInput.ImageURL = toPtr(imageURL)
	return p
}

func (p ProfileInput) CurrentCharacterID(currentCharacterID primitive.ObjectID) ProfileInput {
	p.ProfileInput.CurrentCharacterID = toPtr(currentCharacterID)
	return p
}

type CharacterInput struct {
	model.CharacterInput
}

func NewCharacterInput() CharacterInput {
	return CharacterInput{model.CharacterInput{}}
}

func (c CharacterInput) Name(name string) CharacterInput {
	c.CharacterInput.Name = toPtr(name)
	return c
}

func (c CharacterInput) Gender(gender bool) CharacterInput {
	c.CharacterInput.Gender = toPtr(gender)
	return c
}

func (c CharacterInput) Tags(tags []string) CharacterInput {
	c.CharacterInput.Tags = tags
	return c
}

type CustomMetricInput struct {
	model.CustomMetricInput
}

func NewCustomMetricInput() CustomMetricInput {
	return CustomMetricInput{model.CustomMetricInput{}}
}

func (c CustomMetricInput) Name(name string) CustomMetricInput {
	c.CustomMetricInput.Name = toPtr(name)
	return c
}

func (c CustomMetricInput) Description(description string) CustomMetricInput {
	c.CustomMetricInput.Description = toPtr(description)
	return c
}

func (c CustomMetricInput) Style(style model.MetricStyleInput) CustomMetricInput {
	c.CustomMetricInput.Style = toPtr(style)
	return c
}

func (c CustomMetricInput) Properties(properties []model.MetricPropertyInput) CustomMetricInput {
	c.CustomMetricInput.Properties = properties
	return c
}

type MetricPropertyInput struct {
	model.MetricPropertyInput
}

func NewMetricPropertyInput() MetricPropertyInput {
	return MetricPropertyInput{model.MetricPropertyInput{}}
}

func (m MetricPropertyInput) Name(name string) MetricPropertyInput {
	m.MetricPropertyInput.Name = toPtr(name)
	return m
}

func (m MetricPropertyInput) Type(t model.MetricPropertyType) MetricPropertyInput {
	m.MetricPropertyInput.Type = toPtr(t)
	return m
}

func (m MetricPropertyInput) Value(value string) MetricPropertyInput {
	m.MetricPropertyInput.Value = toPtr(value)
	return m
}

func (m MetricPropertyInput) Unit(unit string) MetricPropertyInput {
	m.MetricPropertyInput.Unit = toPtr(unit)
	return m
}

func toPtr[M any](s M) *M {
	return &s
}
