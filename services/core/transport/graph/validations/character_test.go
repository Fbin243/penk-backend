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

func TestValidateCategoryInput(t *testing.T) {
	type testCase struct {
		name     string
		category entity.CategoryInput
		hasError bool
	}

	tests := []testCase{
		{
			name: "valid category",
			category: entity.CategoryInput{
				Name:        "Category",
				Description: lo.ToPtr("This is a description"),
			},
			hasError: false,
		},
		{
			name: "description with maximum length",
			category: entity.CategoryInput{
				Name:        "Category",
				Description: lo.ToPtr("This is a enough description with maximum length. This is a enough description with maximum length. This is a enough description with maximum length. This is a enough description with maximum length. This is a enough description with maximum length. Th..."),
			},
			hasError: false,
		},
		{
			name: "too long description",
			category: entity.CategoryInput{
				Name:        "Category",
				Description: lo.ToPtr("This is a very long description that exceeds the maximum allowed length of two hundred fifty five characters. This is a very long description that exceeds the maximum allowed length of two hundred fifty five characters. This is a very long description t..."),
			},
			hasError: true,
		},
		{
			name: "color with valid hex",
			category: entity.CategoryInput{
				Name: "Category",
				Style: lo.ToPtr(entity.CategoryStyleInput{
					Color: "#ff0000",
				}),
			},
			hasError: false,
		},
		{
			name: "color with invalid hex",
			category: entity.CategoryInput{
				Name: "Category",
				Style: lo.ToPtr(entity.CategoryStyleInput{
					Color: "ff0000",
				}),
			},
			hasError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validations.ValidateCategoryInput(tc.category)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateMetric(t *testing.T) {
	type testCase struct {
		name     string
		metric   entity.MetricInput
		hasError bool
	}

	tests := []testCase{
		{
			name: "valid metric",
			metric: entity.MetricInput{
				Name:  "Category",
				Value: 10,
				Unit:  "kg",
			},
			hasError: false,
		},
		{
			name: "without unit",
			metric: entity.MetricInput{
				Name:  "Category",
				Value: 10,
				Unit:  "",
			},
			hasError: false,
		},
		{
			name: "missing name",
			metric: entity.MetricInput{
				Name:  "",
				Value: 10,
				Unit:  "kg",
			},
			hasError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validations.ValidateMetricInput(tc.metric)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
