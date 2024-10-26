package validations_test

import (
	"testing"

	"tenkhours/services/core/graph/model"
	"tenkhours/services/core/graph/validations"

	"github.com/stretchr/testify/assert"
)

func TestValidateCharacterInput(t *testing.T) {
	type testCase struct {
		name      string
		character CharacterInput
		hasError  bool
	}

	characterInput := NewCharacterInput().Name("Hero").Gender(true).Tags([]string{"warrior", "legend"})

	tests := []testCase{
		{
			name:      "valid character",
			character: characterInput,
			hasError:  false,
		},
		{
			name:      "empty name",
			character: characterInput.Name(""),
			hasError:  true,
		},
		{
			name:      "name too long",
			character: characterInput.Name("This is a very long name that exceeds the maximum allowed length of fifty characters."),
			hasError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validations.ValidateUpdateCharacterInput(tc.character.CharacterInput)
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
		metric   CustomMetricInput
		hasError bool
	}

	metricInput := NewCustomMetricInput().Name("Metric").Description("Metric description").Style(model.MetricStyleInput{}).Properties([]model.MetricPropertyInput{})

	tests := []testCase{
		{
			name:     "valid metric",
			metric:   metricInput,
			hasError: false,
		},
		{
			name:     "description with maximum length",
			metric:   metricInput.Description("This is a enough description with maximum length. This is a enough description with maximum length. This is a enough description with maximum length. This is a enough description with maximum length. This is a enough description with maximum length. Th..."),
			hasError: false,
		},
		{
			name:     "too long description",
			metric:   metricInput.Description("This is a very long description that exceeds the maximum allowed length of two hundred fifty five characters. This is a very long description that exceeds the maximum allowed length of two hundred fifty five characters. This is a very long description t..."),
			hasError: true,
		},
		{
			name:     "color with valid hex",
			metric:   metricInput.Style(model.MetricStyleInput{Color: toPtr("#ff0000")}),
			hasError: false,
		},
		{
			name:     "color with invalid hex",
			metric:   metricInput.Style(model.MetricStyleInput{Color: toPtr("#ff000")}),
			hasError: true,
		},
		{
			name: "invalid property",
			metric: metricInput.Properties([]model.MetricPropertyInput{
				NewMetricPropertyInput().Name("").Type(model.MetricPropertyTypeNumber).Value("10").Unit("kg").MetricPropertyInput,
			}),
			hasError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validations.ValidateUpdateCustomMetricInput(tc.metric.CustomMetricInput)
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
		property MetricPropertyInput
		hasError bool
	}

	propertyInput := NewMetricPropertyInput().Name("Property").Type(model.MetricPropertyTypeNumber).Value("10").Unit("kg")

	tests := []testCase{
		{
			name:     "valid property",
			property: propertyInput,
			hasError: false,
		},
		{
			name:     "without unit",
			property: propertyInput.Unit(""),
			hasError: false,
		},
		{
			name:     "missing name",
			property: propertyInput.Name(""),
			hasError: true,
		},
		{
			name:     "missing value",
			property: propertyInput.Value(""),
			hasError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validations.ValidateUpdateMetricPropertyInput(tc.property.MetricPropertyInput)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
