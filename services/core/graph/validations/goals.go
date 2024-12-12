package validations

import (
	"fmt"
	"tenkhours/services/core/graph/model"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func ValidateGoalInput(input model.GoalInput) error {
	// Validate the goal struct
	validate = validator.New()
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
