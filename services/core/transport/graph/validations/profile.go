package validations

import "tenkhours/services/core/entity"

func ValidateProfileInput(profile entity.ProfileInput) error {
	if err := GetValidator().Struct(profile); err != nil {
		return err
	}

	return nil
}
