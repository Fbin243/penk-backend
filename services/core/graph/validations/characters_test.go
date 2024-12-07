package validations_test

import (
	"testing"

	"tenkhours/services/core/graph/model"
	"tenkhours/services/core/graph/validations"
	"tenkhours/services/core/repo"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestValidateCharacterInput(t *testing.T) {
	type testCase struct {
		name      string
		character model.CharacterInput
		hasError  bool
	}

	tests := []testCase{
		{
			name:      "valid character",
			character: model.CharacterInput{Name: "Character"},
			hasError:  false,
		},
		{
			name:      "empty name",
			character: model.CharacterInput{Name: ""},
			hasError:  true,
		},
		{
			name:      "name too long",
			character: model.CharacterInput{Name: "This is a very long name that exceeds the maximum allowed length of fifty characters. This is a very long name that exceeds the maximum allowed length of fifty characters."},
			hasError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validations.ValidateCharacterInput(tc.character)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateCustomMetric(t *testing.T) {
	type testCase struct {
		name     string
		metric   model.CustomMetricInput
		hasError bool
	}

	tests := []testCase{
		{
			name: "valid metric",
			metric: model.CustomMetricInput{
				Name:        "Metric",
				Description: lo.ToPtr("This is a description"),
			},
			hasError: false,
		},
		{
			name: "description with maximum length",
			metric: model.CustomMetricInput{
				Name:        "Metric",
				Description: lo.ToPtr("This is a enough description with maximum length. This is a enough description with maximum length. This is a enough description with maximum length. This is a enough description with maximum length. This is a enough description with maximum length. Th..."),
			},
			hasError: false,
		},
		{
			name: "too long description",
			metric: model.CustomMetricInput{
				Name:        "Metric",
				Description: lo.ToPtr("This is a very long description that exceeds the maximum allowed length of two hundred fifty five characters. This is a very long description that exceeds the maximum allowed length of two hundred fifty five characters. This is a very long description t..."),
			},
			hasError: true,
		},
		{
			name: "color with valid hex",
			metric: model.CustomMetricInput{
				Name: "Metric",
				Style: lo.ToPtr(model.MetricStyleInput{
					Color: "#ff0000",
				}),
			},
			hasError: false,
		},
		{
			name: "color with invalid hex",
			metric: model.CustomMetricInput{
				Name: "Metric",
				Style: lo.ToPtr(model.MetricStyleInput{
					Color: "ff0000",
				}),
			},
			hasError: true,
		},
		{
			name: "invalid property",
			metric: model.CustomMetricInput{
				Name: "Metric",
				Style: lo.ToPtr(model.MetricStyleInput{
					Color: "ff0000",
				}),
				Properties: []model.MetricPropertyInput{
					{
						Name:  "Property",
						Type:  "",
						Value: "10",
						Unit:  "kg",
					},
				},
			},
			hasError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validations.ValidateCustomMetricInput(tc.metric)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateMetricProperty(t *testing.T) {
	type testCase struct {
		name     string
		property model.MetricPropertyInput
		hasError bool
	}

	tests := []testCase{
		{
			name: "valid property",
			property: model.MetricPropertyInput{
				Name:  "Property",
				Type:  repo.MetricPropertyTypeString,
				Value: "10",
				Unit:  "kg",
			},
			hasError: false,
		},
		{
			name: "without unit",
			property: model.MetricPropertyInput{
				Name:  "Property",
				Type:  repo.MetricPropertyTypeString,
				Value: "10",
				Unit:  "",
			},
			hasError: true,
		},
		{
			name: "missing name",
			property: model.MetricPropertyInput{
				Name:  "",
				Type:  repo.MetricPropertyTypeString,
				Value: "10",
				Unit:  "kg",
			},
			hasError: true,
		},
		{
			name: "missing value",
			property: model.MetricPropertyInput{
				Name:  "Property",
				Type:  repo.MetricPropertyTypeString,
				Value: "",
				Unit:  "kg"},
			hasError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validations.ValidateMetricPropertyInput(tc.property)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
