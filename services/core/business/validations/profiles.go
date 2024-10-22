package validations

import (
	"tenkhours/services/core/repo"

	"github.com/go-playground/validator/v10"
)

func ValidateProfile(profile repo.Profile) error {
	validate := validator.New()

	return validate.Struct(profile)
}
