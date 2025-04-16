package entity

import (
	"time"

	"tenkhours/pkg/types"
)

type HabitLogInput struct {
	Timestamp string  `json:"timestamp"`
	HabitID   string  `json:"habitID"`
	Value     float64 `json:"value"`
}

type HabitLogFilter struct {
	HabitID   *string    `json:"habit_id"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}

type HabitLogOrderBy struct {
	Timestamp *types.SortOrder `json:"time_stamp"`
}
