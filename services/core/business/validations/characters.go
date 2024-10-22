package validations

import (
	"tenkhours/services/core/repo"

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

func ValidateCharacter(character repo.Character) error {
	validate := validator.New()
	validate.RegisterValidation("tags_valid", ValidateTag)

	return validate.Struct(character)
}

func ValidateCustomMetric(customMetric repo.CustomMetric) error {
	return validator.New().Struct(customMetric)
}

func ValidateMetricProperty(metricProperty repo.MetricProperty) error {
	return validator.New().Struct(metricProperty)
}
