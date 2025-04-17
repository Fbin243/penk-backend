package entity

import (
	"time"

	"tenkhours/pkg/types"
)

type HabitLogPineline struct {
	Filter  *HabitLogFilter
	OrderBy *HabitLogOrderBy
	Limit   *int
	Offset  *int
}

type HabitLogFilter struct {
	HabitIDs  []string    `json:"habit_ids"`
	HabitID   *string     `json:"habit_id"`
	StartTime *time.Time  `json:"start_time"`
	EndTime   *time.Time  `json:"end_time"`
	Reset     *HabitReset `json:"reset"`
}

type HabitLogOrderBy struct {
	Timestamp *types.SortOrder `json:"time_stamp"`
}
