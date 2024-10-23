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
	validate := validator.New()
	validate.RegisterValidation("tags_valid", ValidateTag)

	return validate.Struct(characterInput)
}

func ValidateCustomMetricInput(customMetricInput model.CustomMetricInput) error {
	return validator.New().Struct(customMetricInput)
}

func ValidateMetricPropertyInput(metricPropertyInput model.MetricPropertyInput) error {
	return validator.New().Struct(metricPropertyInput)
}
