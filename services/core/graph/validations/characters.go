package validations

import (
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

func ValidateCharacterInput(characterInput model.CharacterInput) error {
	validate := GetValidator()
	validate.RegisterValidation("tags_valid", ValidateTag)
	if err := validate.Struct(characterInput); err != nil {
		return err
	}

	return nil
}

func ValidateCustomMetricInput(customMetricInput model.CustomMetricInput) error {
	if err := GetValidator().Struct(customMetricInput); err != nil {
		return err
	}

	return nil
}

func ValidateMetricPropertyInput(metricPropertyInput model.MetricPropertyInput) error {
	if err := GetValidator().Struct(metricPropertyInput); err != nil {
		return err
	}

	return nil
}
