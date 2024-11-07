package repo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Fish struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ProfileID primitive.ObjectID `json:"profileID,omitempty" bson:"profile_id,omitempty"`
	Numbers   int32              `json:"numbers" bson:"numbers"`
	Type      string             `json:"type" bson:"type"`
}

type Cod struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ProfileID primitive.ObjectID `json:"profileID,omitempty" bson:"profile_id,omitempty"`
	Type      string             `json:"type" bson:"type"`
}
