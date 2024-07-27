package analyticsdb

import (
	"time"

	"tenkhours/pkg/db/coredb"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Metadata struct {
	CharacterID primitive.ObjectID `json:"characterID" bson:"character_id"`
	UserID      primitive.ObjectID `json:"userID" bson:"user_id"`
}

type Snapshot struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
	Metadata  Metadata           `json:"metadata" bson:"metadata"`
	Character coredb.Character   `json:"character" bson:"character"`
	Asset     interface{}        `json:"asset" bson:"asset"`
}
