package entity

import (
	"time"
)

type TimeTracking struct {
	ID          string    `json:"id,omitempty"          bson:"_id"`
	CharacterID string    `json:"characterID,omitempty" bson:"character_id"`
	CategoryID  string    `json:"categoryID,omitempty"  bson:"category_id,omitempty"`
	StartTime   time.Time `json:"startTime,omitempty"   bson:"start_time"`
	EndTime     time.Time `json:"endTime,omitempty"     bson:"end_time,omitempty"`
}

type TimeTrackingWithFish struct {
	TimeTracking *TimeTracking `json:"timeTracking"`
	Normal       int           `json:"normal"`
	Gold         int           `json:"gold"`
}
