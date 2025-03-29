package mongorepo

import (
	mongodb "tenkhours/pkg/db/mongo"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reward struct {
	*mongodb.BaseEntity `                           bson:",inline"`
	ProfileOID          primitive.ObjectID `json:"profileID,omitempty" bson:"profile_id,omitempty"`
	FishCount           int32              `json:"fishCount"              bson:"fish_count"`
	ClaimedAt           time.Time          `json:"claimedAt,omitempty" bson:"claimed_at"`
	StreakCount         int32              `json:"streakCount,omitempty" bson:"streak_count"`
}

func (Reward) IsEntity() {}

func (f *Reward) ProfileID(id string) {
	f.ProfileOID = mongodb.ToObjectID(id)
}
