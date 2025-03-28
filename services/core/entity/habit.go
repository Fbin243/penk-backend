package entity

import (
	"time"

	"tenkhours/pkg/db/base"
)

type CompletionType string

const (
	CompletionTypeMetric   CompletionType = "Metric"
	CompletionTypeCheckbox CompletionType = "Checkbox"
	CompletionTypeTime     CompletionType = "Time"
)

type Habit struct {
	*base.BaseEntity `                           bson:",inline"`
	CharacterID      string         `json:"characterID"         bson:"character_id"`
	CategoryID       *string        `json:"categoryID"          bson:"category_id"`
	CompletionType   CompletionType `json:"completionType"      bson:"completion_type"`
	Name             string         `json:"name"                bson:"name"`
	Value            float64        `json:"value,omitempty"     bson:"value"`
	Unit             *string        `json:"unit,omitempty"      bson:"unit"`
	StartTime        time.Time      `json:"startTime,omitempty" bson:"start_time"`
	EndTime          *time.Time     `json:"endTime,omitempty"   bson:"end_time"`
	Frequency        string         `json:"frequency"           bson:"frequency"`
}
