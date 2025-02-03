package entity

import (
	"time"
)

type TimeTracking struct {
	ID             string    `json:"id,omitempty"             bson:"_id"`
	CharacterID    string    `json:"characterID,omitempty"    bson:"character_id"`
	CustomMetricID string    `json:"customMetricID,omitempty" bson:"custom_metric_id,omitempty"`
	StartTime      time.Time `json:"startTime,omitempty"      bson:"start_time"`
	EndTime        time.Time `json:"endTime,omitempty"        bson:"end_time,omitempty"`
	// TODO: Will be move to the better place later on
	MinDurationTime int32 `json:"minDurationTime" bson:"min_duration_time"`
	MaxDurationTime int32 `json:"maxDurationTime" bson:"max_duration_time"`
}

type TimeTrackingWithFish struct {
	TimeTracking *TimeTracking `json:"timeTracking"`
	Normal       int           `json:"normal"`
	Gold         int           `json:"gold"`
}
