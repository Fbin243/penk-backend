package validations

import (
	"fmt"
	"time"

	"tenkhours/services/core/entity"
)

func ValidateGoalInput(input entity.GoalInput) error {
	validate := GetValidator()
	if err := validate.Struct(input); err != nil {
		return err
	}

	if input.StartDate.After(input.EndDate) {
		return fmt.Errorf("start date must be before end date")
	}

	if input.EndDate.Before(time.Now()) {
		return fmt.Errorf("end date must be in the future")
	}

	return nil
}
