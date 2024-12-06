package validations

import (
	"tenkhours/pkg/utils"
	"tenkhours/services/core/graph/model"

	"github.com/go-playground/validator/v10"
)

func ValidateTag(fl validator.FieldLevel) bool {
	tags, ok := fl.Field().Interface().([]string)
	if !ok {
		return false
	}

	for _, tag := range tags {
		if len(tag) < 1 || len(tag) > 20 {
			return false
		}
	}
	return true
}

func ValidateCreateCharacterInput(characterInput model.CharacterInput) error {
	if err := utils.NewValidateBuilder().Required().Validate(characterInput.Name); err != nil {
		return err
	}

	return ValidateUpdateCharacterInput(characterInput)
}

func ValidateUpdateCharacterInput(characterInput model.CharacterInput) error {
	if characterInput.Name != nil {
		err := utils.NewValidateBuilder().Required().Min(1).Max(50).Validate(*characterInput.Name)
		if err != nil {
			return err
		}
	}

	err := utils.NewValidateBuilder().OmitEmpty().Custom("tags_valid", ValidateTag).Validate(characterInput.Tags)
	if err != nil {
		return err
	}

	if characterInput.CustomMetrics != nil {
		for _, customMetric := range characterInput.CustomMetrics {
			if err := ValidateCreateCustomMetricInput(customMetric); err != nil {
				return err
			}
		}
	}

	return nil
}

func ValidateCreateCustomMetricInput(customMetricInput model.CustomMetricInput) error {
	if err := utils.NewValidateBuilder().Required().Validate(customMetricInput.Name); err != nil {
		return err
	}

	return ValidateUpdateCustomMetricInput(customMetricInput)
}

func ValidateUpdateCustomMetricInput(customMetricInput model.CustomMetricInput) error {
	if customMetricInput.Name != nil {
		err := utils.NewValidateBuilder().Required().Min(1).Max(50).Validate(*customMetricInput.Name)
		if err != nil {
			return err
		}
	}

	if customMetricInput.Description != nil {
		if err := utils.NewValidateBuilder().OmitEmpty().Max(255).Validate(*customMetricInput.Description); err != nil {
			return err
		}
	}

	if customMetricInput.Style != nil {
		if err := ValidateMetricStyleInput(*customMetricInput.Style); err != nil {
			return err
		}
	}

	if customMetricInput.Properties != nil {
		for _, property := range customMetricInput.Properties {
			if err := ValidateCreateMetricPropertyInput(property); err != nil {
				return err
			}
		}
	}

	return nil
}

func ValidateMetricStyleInput(metricStyleInput model.MetricStyleInput) error {
	if metricStyleInput.Color != nil {
		if err := utils.NewValidateBuilder().OmitEmpty().Color().Validate(*metricStyleInput.Color); err != nil {
			return err
		}
	}

	return nil
}

func ValidateCreateMetricPropertyInput(metricPropertyInput model.MetricPropertyInput) error {
	if err := utils.NewValidateBuilder().Required().Validate(metricPropertyInput.Name); err != nil {
		return err
	}

	if err := utils.NewValidateBuilder().Required().Validate(metricPropertyInput.Type); err != nil {
		return err
	}

	return ValidateUpdateMetricPropertyInput(metricPropertyInput)
}

func ValidateUpdateMetricPropertyInput(metricPropertyInput model.MetricPropertyInput) error {
	if metricPropertyInput.Name != nil {
		if err := utils.NewValidateBuilder().Required().Min(1).Max(50).Validate(*metricPropertyInput.Name); err != nil {
			return err
		}
	}

	if metricPropertyInput.Unit != nil {
		if err := utils.NewValidateBuilder().OmitEmpty().Max(50).Validate(*metricPropertyInput.Unit); err != nil {
			return err
		}
	}

	return nil
}
