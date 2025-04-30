package entity

import (
	"time"

	"tenkhours/pkg/types"
)

type GoalInput struct {
	ID          *string           `json:"id,omitempty"`
	Name        string            `json:"name"                  validate:"min=1,max=50"`
	Description *string           `json:"description,omitempty" validate:"omitempty,max=255"`
	StartTime   time.Time         `json:"startTime"`
	EndTime     time.Time         `json:"endTime"`
	Metrics     []GoalMetricInput `json:"metrics"`
	Checkboxes  []CheckboxInput   `json:"checkboxes"`
}

type GoalMetricInput struct {
	ID          string          `json:"id"`
	Condition   MetricCondition `json:"condition"`
	TargetValue *float64        `json:"targetValue"`
	RangeValue  *RangeInput     `json:"rangeValue,omitempty"`
}

type RangeInput struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type CheckboxInput struct {
	ID    *string `json:"id"`
	Name  string  `json:"name"`
	Value bool    `json:"value"`
}

type MetricCondition string

const (
	MetricConditionGreaterThan      MetricCondition = "gt"
	MetricConditionLessThan         MetricCondition = "lt"
	MetricConditionEqual            MetricCondition = "eq"
	MetricConditionLessThanEqual    MetricCondition = "lte"
	MetricConditionGreaterThanEqual MetricCondition = "gte"
	MetricConditionInRange          MetricCondition = "ir"
)

type GoalPipeline struct {
	Filter  *GoalFilter
	OrderBy *GoalOrderBy
	*types.Pagination
}

type GoalFilter struct {
	CharacterID *string     `json:"character_id"`
	Status      *GoalStatus `json:"status"`
}

type GoalOrderBy struct{}
