package entity

import "time"

type GoalCategoryInput struct {
	ID      string            `json:"id"`
	Metrics []GoalMetricInput `json:"metrics,omitempty"`
}

type GoalInput struct {
	ID          *string             `json:"id,omitempty"`
	Name        string              `json:"name"                  validate:"min=1,max=50"`
	Description *string             `json:"description,omitempty" validate:"omitempty,max=255"`
	StartDate   time.Time           `json:"startDate"`
	EndDate     time.Time           `json:"endDate"`
	Target      []GoalCategoryInput `json:"target,omitempty"`
}

type GoalMetricInput struct {
	ID    string  `json:"id"`
	Value float64 `json:"value"`
}
