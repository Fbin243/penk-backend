package validations

import (
	"testing"

	"tenkhours/services/core/graph/model"

	"github.com/stretchr/testify/assert"
)

func TestValidateCharacter(t *testing.T) {
	type testCase struct {
		name      string
		character model.CharacterInput
		hasError  bool
	}

	character := model.CharacterInput{
		Name:   toPtr("Hero"),
		Gender: toPtr(true),
		Tags:   []string{"warrior", "legend"},
	}

	tests := []testCase{
		{
			name:      "valid character",
			character: character,
			hasError:  false,
		},
		{
			name: "empty name",
			character: func(char model.CharacterInput) model.CharacterInput {
				char.Name = toPtr("")
				return char
			}(character),
			hasError: true,
		},
		{
			name: "name too long",
			character: func(char model.CharacterInput) model.CharacterInput {
				char.Name = toPtr("This is a very long name that exceeds the maximum allowed length of fifty characters")
				return char
			}(character),
			hasError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateCharacterInput(tc.character)
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

	metric := model.CustomMetricInput{
		Name:        toPtr("Metric"),
		Description: toPtr("Metric description"),
		Style:       &model.MetricStyleInput{},
		Properties:  []model.MetricPropertyInput{},
	}

	tests := []testCase{
		{
			name:     "valid metric",
			metric:   metric,
			hasError: false,
		},
		{
			name: "description with maximum length",
			metric: func(m model.CustomMetricInput) model.CustomMetricInput {
				m.Description = toPtr("This is a enough description with maximum length. This is a enough description with maximum length. This is a enough description with maximum length. This is a enough description with maximum length. This is a enough description with maximum length. Th...")
				return m
			}(metric),
			hasError: false,
		},
		{
			name: "too long description",
			metric: func(m model.CustomMetricInput) model.CustomMetricInput {
				m.Description = toPtr("This is a very long description that exceeds the maximum allowed length of two hundred fifty five characters. This is a very long description that exceeds the maximum allowed length of two hundred fifty five characters. This is a very long description t...")
				return m
			}(metric),
			hasError: true,
		},
		{
			name: "color with valid hex",
			metric: func(m model.CustomMetricInput) model.CustomMetricInput {
				m.Style.Color = toPtr("#ff0000")
				return m
			}(metric),
			hasError: false,
		},
		{
			name: "color with invalid hex",
			metric: func(m model.CustomMetricInput) model.CustomMetricInput {
				m.Style.Color = toPtr("#invalid")
				return m
			}(metric),
			hasError: true,
		},
		{
			name: "invalid property",
			metric: func(m model.CustomMetricInput) model.CustomMetricInput {
				m.Properties = []model.MetricPropertyInput{
					{
						Name:  toPtr(""),
						Type:  toPtr(model.MetricPropertyTypeString),
						Value: toPtr("10"),
						Unit:  toPtr("kg"),
					},
				}
				return m
			}(metric),
			hasError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateCustomMetricInput(tc.metric)
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

	property := model.MetricPropertyInput{
		Name:  toPtr("Property"),
		Type:  toPtr(model.MetricPropertyTypeNumber),
		Value: toPtr("10"),
		Unit:  toPtr("kg"),
	}

	tests := []testCase{
		{
			name:     "valid property",
			property: property,
			hasError: false,
		},
		{
			name: "without unit",
			property: func(p model.MetricPropertyInput) model.MetricPropertyInput {
				p.Unit = nil
				return p
			}(property),
			hasError: false,
		},
		{
			name: "missing name",
			property: func(p model.MetricPropertyInput) model.MetricPropertyInput {
				p.Name = nil
				return p
			}(property),
			hasError: true,
		},
		{
			name: "missing type",
			property: func(p model.MetricPropertyInput) model.MetricPropertyInput {
				p.Type = nil
				return p
			}(property),
			hasError: true,
		},
		{
			name: "missing value",
			property: func(p model.MetricPropertyInput) model.MetricPropertyInput {
				p.Value = nil
				return p
			}(property),
			hasError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMetricPropertyInput(tc.property)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func toPtr[M any](s M) *M {
	return &s
}
