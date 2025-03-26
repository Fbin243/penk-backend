package entity

import (
	"tenkhours/pkg/db/base"
	"time"
)

type Reward struct {
	*base.BaseEntity `                           bson:",inline"`
	ProfileID        string    `json:"profileID,omitempty" bson:"profile_id,omitempty"`
	FishCount        int32     `json:"fishCount"              bson:"fish_count"`
	ClaimedAt        time.Time `json:"claimedAt,omitempty" bson:"claimed_at,omitempty"`
	StreakCount      int32     `json:"streakCount,omitempty" bson:"streak_count"`
}

func (Reward) IsEntity() {}
