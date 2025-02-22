package validations

import (
	"time"

	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"
)

func ValidateGoalInput(input entity.GoalInput) error {
	validate := GetValidator()
	if err := validate.Struct(input); err != nil {
		return err
	}

	if input.StartTime.After(input.EndTime) {
		return errors.NewGQLError(errors.ErrCodeBadRequest, "start time must be before end time")
	}

	if input.EndTime.Before(time.Now()) {
		return errors.NewGQLError(errors.ErrCodeBadRequest, "end time must be in the future")
	}

	for _, metric := range input.Metrics {
		if metric.Condition == entity.MetricConditionInRange {
			if metric.RangeValue == nil || metric.TargetValue != nil {
				return errors.NewGQLError(errors.ErrCodeBadRequest, "condition and values not match")
			}

			if metric.RangeValue.Min > metric.RangeValue.Max {
				return errors.NewGQLError(errors.ErrCodeBadRequest, "min value must be less than max value")
			}
		} else {
			if metric.TargetValue == nil || metric.RangeValue != nil {
				return errors.NewGQLError(errors.ErrCodeBadRequest, "condition and values not match")
			}
		}
	}

	return nil
}
