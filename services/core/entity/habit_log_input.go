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

type HabitLogPipeline struct {
	Filter  *HabitLogFilter
	OrderBy *HabitLogOrderBy
	*types.Pagination
}

type HabitLogFilter struct {
	HabitIDs      []string    `json:"habit_ids"`
	HabitID       *string     `json:"habit_id"`
	StartTime     *time.Time  `json:"start_time"`
	EndTime       *time.Time  `json:"end_time"`
	ResetDuration *HabitReset `json:"resetDuration"`
}

type HabitLogOrderBy struct {
	Timestamp *types.SortOrder `json:"time_stamp"`
}
