package validations

import (
	"tenkhours/services/core/graph/model"

	"github.com/go-playground/validator/v10"
)

func ValidateProfileInput(profile model.ProfileInput) error {
	validate := validator.New()
	if err := validate.Struct(profile); err != nil {
		return err
	}

	return nil
}
