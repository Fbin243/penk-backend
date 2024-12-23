package validations

import (
	"tenkhours/services/core/graph/model"
)

func ValidateProfileInput(profile model.ProfileInput) error {
	if err := GetValidator().Struct(profile); err != nil {
		return err
	}

	return nil
}
