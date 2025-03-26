package entity

import (
	"time"

	"tenkhours/pkg/db/base"
)

type GoalStatus string

const (
	GoalStatusPlanned    GoalStatus = "Planned"
	GoalStatusInProgress GoalStatus = "InProgress"
	GoalStatusCompleted  GoalStatus = "Completed"
	GoalStatusOverdue    GoalStatus = "Overdue"
)

type Goal struct {
	*base.BaseEntity `                     bson:",inline"`
	CharacterID      string       `json:"characterID"   bson:"character_id"`
	Name             string       `json:"name"          bson:"name"`
	Description      string       `json:"description"   bson:"description"`
	StartTime        time.Time    `json:"startTime"     bson:"start_time"`
	EndTime          time.Time    `json:"endTime"       bson:"end_time"`
	CompletedTime    *time.Time   `json:"completedTime" bson:"completed_time"`
	Metrics          []GoalMetric `json:"metrics"       bson:"metrics"`
	Checkboxes       []Checkbox   `json:"checkboxes"    bson:"checkboxes"`
}

type GoalMetric struct {
	ID          string          `json:"id"          bson:"id"`
	Condition   MetricCondition `json:"condition"   bson:"condition"`
	TargetValue *float64        `json:"targetValue" bson:"target_value"`
	RangeValue  *Range          `json:"rangeValue"  bson:"range_value"`
}

type Range struct {
	Min float64 `json:"min" bson:"min"`
	Max float64 `json:"max" bson:"max"`
}

func (m *GoalMetric) Evaluate(currentMetric Metric) bool {
	switch m.Condition {
	case MetricConditionEqual:
		return currentMetric.Value == *m.TargetValue
	case MetricConditionGreaterThan:
		return currentMetric.Value > *m.TargetValue
	case MetricConditionLessThan:
		return currentMetric.Value < *m.TargetValue
	case MetricConditionGreaterThanEqual:
		return currentMetric.Value >= *m.TargetValue
	case MetricConditionLessThanEqual:
		return currentMetric.Value <= *m.TargetValue
	case MetricConditionInRange:
		return currentMetric.Value >= m.RangeValue.Min && currentMetric.Value <= m.RangeValue.Max
	default:
		return false
	}
}

func (g *Goal) IsCompleted(metricMap map[string]Metric) bool {
	for _, targetMetric := range g.Metrics {
		currentMetric := metricMap[targetMetric.ID]
		if !targetMetric.Evaluate(currentMetric) {
			return false
		}
	}

	for _, checkbox := range g.Checkboxes {
		if !checkbox.Value {
			return false
		}
	}

	return true
}

func (g *Goal) EvaluateStatus() GoalStatus {
	if g.CompletedTime != nil {
		return GoalStatusCompleted
	}

	if time.Now().After(g.EndTime) {
		return GoalStatusOverdue
	}

	if time.Now().After(g.StartTime) {
		return GoalStatusInProgress
	}

	return GoalStatusPlanned
}
