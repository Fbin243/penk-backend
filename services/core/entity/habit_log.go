package entity

import (
	"time"

	"tenkhours/pkg/db/base"
)

type HabitLog struct {
	*base.BaseEntity `                 bson:",inline"`
	Timestamp        time.Time `json:"timestamp" bson:"timestamp"`
	HabitID          string    `json:"habitID"   bson:"habit_id"`
	Value            float64   `json:"value"     bson:"value"`
	Percent          float64   `json:"percent"   bson:"percent"`
}
