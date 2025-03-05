package entity

import (
	"time"

	"tenkhours/pkg/db/base"
)

type (
	GoalFinishStatus string
	GoalExpireStatus string
)

const (
	GoalFinishStatusFinished   GoalFinishStatus = "Finished"
	GoalFinishStatusUnfinished GoalFinishStatus = "Unfinished"
	GoalExpireStatusExpired    GoalExpireStatus = "Expired"
	GoalExpireStatusUnexpired  GoalExpireStatus = "Unexpired"
)

type GoalStatusFilter struct {
	FinishStatus *GoalFinishStatus
	ExpireStatus *GoalExpireStatus
}

type Goal struct {
	*base.BaseEntity `                   bson:",inline"`
	CharacterID      string           `json:"characterID" bson:"character_id"`
	Name             string           `json:"name"        bson:"name"`
	Description      string           `json:"description" bson:"description"`
	StartTime        time.Time        `json:"startTime"   bson:"start_time"`
	EndTime          time.Time        `json:"endTime"     bson:"end_time"`
	Status           GoalFinishStatus `json:"status"      bson:"status"`
	Target           GoalTarget       `json:"target"      bson:"target"`
}

type GoalTarget struct {
	Categories []GoalCategory `json:"categories" bson:"categories"`
	Metrics    []GoalMetric   `json:"metrics"    bson:"metrics"`
	Checkboxes []Checkbox     `json:"checkboxes" bson:"checkboxes"`
}

type GoalCategory struct {
	*Category `               bson:",inline"`
	Metrics   []GoalMetric `json:"metrics" bson:"metrics"`
}

type GoalMetric struct {
	*Metric     `                   bson:",inline"`
	Condition   MetricCondition `json:"condition"   bson:"condition"`
	TargetValue *float64        `json:"targetValue" bson:"target_value"`
	RangeValue  *Range          `json:"rangeValue"  bson:"range_value"`
}

type Range struct {
	Min float64 `json:"min" bson:"min"`
	Max float64 `json:"max" bson:"max"`
}

type Checkbox struct {
	ID    string `json:"id"    bson:"id"`
	Name  string `json:"name"  bson:"name"`
	Value bool   `json:"value" bson:"value"`
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

func (g *Goal) UpdateStatus(metricMap map[string]Metric) {
	for _, category := range g.Target.Categories {
		for _, targetMetric := range category.Metrics {
			currentMetric := metricMap[targetMetric.ID]
			if !targetMetric.Evaluate(currentMetric) {
				return
			}
		}
	}

	for _, targetMetric := range g.Target.Metrics {
		currentMetric := metricMap[targetMetric.ID]
		if !targetMetric.Evaluate(currentMetric) {
			return
		}
	}

	for _, checkbox := range g.Target.Checkboxes {
		if !checkbox.Value {
			return
		}
	}

	g.Status = GoalFinishStatusFinished
}
