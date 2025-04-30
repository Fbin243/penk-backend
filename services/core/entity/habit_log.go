package entity

import (
	"tenkhours/pkg/db/base"
)

type HabitLog struct {
	*base.BaseEntity `                 bson:",inline"`
	Timestamp        string  `json:"timestamp" bson:"timestamp"`
	HabitID          string  `json:"habitID"   bson:"habit_id"`
	Value            float64 `json:"value"     bson:"value"`
}
