package validations

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func GetValidator() *validator.Validate {
	if validate != nil {
		return validate
	}

	validate = validator.New()
	return validate
}
