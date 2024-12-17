package repo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Fish struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ProfileID primitive.ObjectID `json:"profileID,omitempty" bson:"profile_id,omitempty"`
	Gold      int32              `json:"gold" bson:"gold"`
	Normal    int32              `json:"normal" bson:"normal"`
}
