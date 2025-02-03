package validations_test

import (
	"testing"

	"tenkhours/services/core/entity"
	"tenkhours/services/core/transport/graph/validations"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestValidateCharacterInput(t *testing.T) {
	type testCase struct {
		name      string
		character entity.CharacterInput
		hasError  bool
	}

	tests := []testCase{
		{
			name:      "valid character",
			character: entity.CharacterInput{Name: "Character"},
			hasError:  false,
		},
		{
			name:      "empty name",
			character: entity.CharacterInput{Name: ""},
			hasError:  true,
		},
		{
			name:      "name too long",
			character: entity.CharacterInput{Name: "This is a very long name that exceeds the maximum allowed length of fifty characters. This is a very long name that exceeds the maximum allowed length of fifty characters."},
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
		metric   entity.CustomMetricInput
		hasError bool
	}

	tests := []testCase{
		{
			name: "valid metric",
			metric: entity.CustomMetricInput{
				Name:        "Metric",
				Description: lo.ToPtr("This is a description"),
			},
			hasError: false,
		},
		{
			name: "description with maximum length",
			metric: entity.CustomMetricInput{
				Name:        "Metric",
				Description: lo.ToPtr("This is a enough description with maximum length. This is a enough description with maximum length. This is a enough description with maximum length. This is a enough description with maximum length. This is a enough description with maximum length. Th..."),
			},
			hasError: false,
		},
		{
			name: "too long description",
			metric: entity.CustomMetricInput{
				Name:        "Metric",
				Description: lo.ToPtr("This is a very long description that exceeds the maximum allowed length of two hundred fifty five characters. This is a very long description that exceeds the maximum allowed length of two hundred fifty five characters. This is a very long description t..."),
			},
			hasError: true,
		},
		{
			name: "color with valid hex",
			metric: entity.CustomMetricInput{
				Name: "Metric",
				Style: lo.ToPtr(entity.MetricStyleInput{
					Color: "#ff0000",
				}),
			},
			hasError: false,
		},
		{
			name: "color with invalid hex",
			metric: entity.CustomMetricInput{
				Name: "Metric",
				Style: lo.ToPtr(entity.MetricStyleInput{
					Color: "ff0000",
				}),
			},
			hasError: true,
		},
		{
			name: "invalid property",
			metric: entity.CustomMetricInput{
				Name: "Metric",
				Style: lo.ToPtr(entity.MetricStyleInput{
					Color: "ff0000",
				}),
				Properties: []entity.MetricPropertyInput{
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
		property entity.MetricPropertyInput
		hasError bool
	}

	tests := []testCase{
		{
			name: "valid property",
			property: entity.MetricPropertyInput{
				Name:  "Property",
				Type:  entity.MetricPropertyTypeString,
				Value: "10",
				Unit:  "kg",
			},
			hasError: false,
		},
		{
			name: "without unit",
			property: entity.MetricPropertyInput{
				Name:  "Property",
				Type:  entity.MetricPropertyTypeString,
				Value: "10",
				Unit:  "",
			},
			hasError: false,
		},
		{
			name: "missing name",
			property: entity.MetricPropertyInput{
				Name:  "",
				Type:  entity.MetricPropertyTypeString,
				Value: "10",
				Unit:  "kg",
			},
			hasError: true,
		},
		{
			name: "missing value",
			property: entity.MetricPropertyInput{
				Name:  "Property",
				Type:  entity.MetricPropertyTypeString,
				Value: "",
				Unit:  "kg",
			},
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
