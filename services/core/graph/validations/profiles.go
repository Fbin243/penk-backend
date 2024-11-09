package validations

import (
	"tenkhours/pkg/utils"
	"tenkhours/services/core/graph/model"
)

func ValidateProfileInput(profile model.ProfileInput) error {
	if err := utils.NewValidateBuilder().OmitEmpty().Min(1).Max(50).Validate(profile.Name); err != nil {
		return err
	}

	return nil
}
