package characters

import (
	"testing"

	"tenkhours/pkg/db/coredb"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestValidateCharacter(t *testing.T) {
	type testCase struct {
		name      string
		character coredb.Character
		hasError  bool
	}

	character := coredb.Character{
		ID:                  primitive.NewObjectID(),
		UserID:              primitive.NewObjectID(),
		Name:                "Hero",
		Gender:              true,
		Tags:                []string{"warrior", "legend"},
		TotalFocusedTime:    100,
		CustomMetrics:       []coredb.CustomMetric{},
		LimitedMetricNumber: 5,
	}

	tests := []testCase{
		{
			name:      "valid character",
			character: character,
			hasError:  false,
		},
		{
			name: "empty name",
			character: func(char coredb.Character) coredb.Character {
				char.Name = ""
				return char
			}(character),
			hasError: true,
		},
		{
			name: "name too long",
			character: func(char coredb.Character) coredb.Character {
				char.Name = "This is a very long name that exceeds the maximum allowed length of fifty characters"
				return char
			}(character),
			hasError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateCharacter(tc.character)
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
		metric   coredb.CustomMetric
		hasError bool
	}

	metric := coredb.CustomMetric{
		ID:                    primitive.NewObjectID(),
		Name:                  "Metric",
		Description:           "Metric description",
		Time:                  10,
		Style:                 coredb.MetricStyle{},
		Properties:            []coredb.MetricProperty{},
		LimitedPropertyNumber: 5,
	}

	tests := []testCase{
		{
			name:     "valid metric",
			metric:   metric,
			hasError: false,
		},
		{
			name: "description with maximum length",
			metric: func(m coredb.CustomMetric) coredb.CustomMetric {
				m.Description = "This is a enough description with maximum length. This is a enough description with maximum length. This is a enough description with maximum length. This is a enough description with maximum length. This is a enough description with maximum length. Th..."
				return m
			}(metric),
			hasError: false,
		},
		{
			name: "too long description",
			metric: func(m coredb.CustomMetric) coredb.CustomMetric {
				m.Description = "This is a very long description that exceeds the maximum allowed length of two hundred fifty five characters. This is a very long description that exceeds the maximum allowed length of two hundred fifty five characters. This is a very long description t..."
				return m
			}(metric),
			hasError: true,
		},
		{
			name: "color with valid hex",
			metric: func(m coredb.CustomMetric) coredb.CustomMetric {
				m.Style.Color = "#ff0000"
				return m
			}(metric),
			hasError: false,
		},
		{
			name: "color with invalid hex",
			metric: func(m coredb.CustomMetric) coredb.CustomMetric {
				m.Style.Color = "#invalid"
				return m
			}(metric),
			hasError: true,
		},
		{
			name: "invalid property",
			metric: func(m coredb.CustomMetric) coredb.CustomMetric {
				m.Properties = []coredb.MetricProperty{
					{
						ID:    primitive.NewObjectID(),
						Name:  "",
						Type:  "number",
						Value: 10,
						Unit:  "kg",
					},
				}
				return m
			}(metric),
			hasError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateCustomMetric(tc.metric)
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
		property coredb.MetricProperty
		hasError bool
	}

	property := coredb.MetricProperty{
		ID:    primitive.NewObjectID(),
		Name:  "Property",
		Type:  "number",
		Value: 10,
		Unit:  "kg",
	}

	tests := []testCase{
		{
			name:     "valid property",
			property: property,
			hasError: false,
		},
		{
			name: "without unit",
			property: func(p coredb.MetricProperty) coredb.MetricProperty {
				p.Unit = ""
				return p
			}(property),
			hasError: false,
		},
		{
			name: "missing name",
			property: func(p coredb.MetricProperty) coredb.MetricProperty {
				p.Name = ""
				return p
			}(property),
			hasError: true,
		},
		{
			name: "missing type",
			property: func(p coredb.MetricProperty) coredb.MetricProperty {
				p.Type = ""
				return p
			}(property),
			hasError: true,
		},
		{
			name: "missing value",
			property: func(p coredb.MetricProperty) coredb.MetricProperty {
				p.Value = nil
				return p
			}(property),
			hasError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMetricProperty(tc.property)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
