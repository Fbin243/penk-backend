package validations

import (
	"tenkhours/services/currency/graph/model"

	"github.com/go-playground/validator/v10"
)

func ValidateFishInput(fish model.FishInput) error {
	validate := validator.New()

	return validate.Struct(fish)
}
