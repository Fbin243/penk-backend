package analyticsdb

import (
	"time"

	"tenkhours/pkg/db/coredb"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Metadata struct {
	CharacterID primitive.ObjectID `json:"characterID" bson:"character_id"`
	ProfileID   primitive.ObjectID `json:"profileID" bson:"profile_id"`
}

type Snapshot struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Timestamp   time.Time          `json:"timestamp" bson:"timestamp"`
	Metadata    Metadata           `json:"metadata" bson:"metadata"`
	Character   coredb.Character   `json:"character" bson:"character"`
	Description string             `json:"description,omitempty" bson:"description,omitempty" validate:"omitempty,max=255"`
	Asset       interface{}        `json:"asset,omitempty" bson:"asset,omitempty"`
}
