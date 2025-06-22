package entity

import (
	"tenkhours/pkg/db/base"
)

type CompletionType string

const (
	CompletionTypeNumber CompletionType = "Number"
	CompletionTypeTime   CompletionType = "Time"
)

type Habit struct {
	*base.BaseEntity `                       bson:",inline"`
	CharacterID      string         `json:"characterID"     bson:"character_id"`
	CategoryID       *string        `json:"categoryID"      bson:"category_id"`
	CompletionType   CompletionType `json:"completionType"  bson:"completion_type"`
	Name             string         `json:"name"            bson:"name"`
	Value            float64        `json:"value,omitempty" bson:"value"`
	Unit             *string        `json:"unit,omitempty"  bson:"unit"`
	RRule            string         `json:"rrule"           bson:"rrule"`
	ResetDuration    HabitReset     `json:"resetDuration"   bson:"reset_duration"`
}

type HabitReset string

const (
	HabitResetDaily   HabitReset = "Daily"
	HabitResetWeekly  HabitReset = "Weekly"
	HabitResetMonthly HabitReset = "Monthly"
)
