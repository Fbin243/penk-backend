package entity

import "time"

type GoalCustomMetricInput struct {
	ID         string                    `json:"id"`
	Properties []GoalMetricPropertyInput `json:"properties,omitempty"`
}

type GoalInput struct {
	ID          *string                 `json:"id,omitempty"`
	Name        string                  `json:"name"                  validate:"min=1,max=50"`
	Description *string                 `json:"description,omitempty" validate:"omitempty,max=255"`
	StartDate   time.Time               `json:"startDate"`
	EndDate     time.Time               `json:"endDate"`
	Target      []GoalCustomMetricInput `json:"target,omitempty"`
}

type GoalMetricPropertyInput struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}
