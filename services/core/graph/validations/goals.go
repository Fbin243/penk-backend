package validations

import (
	"fmt"
	"tenkhours/services/core/graph/model"
	"time"
)

func ValidateGoalInput(input model.GoalInput) error {
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
