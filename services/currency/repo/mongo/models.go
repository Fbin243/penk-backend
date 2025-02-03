package mongorepo

import (
	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Fish struct {
	*mongodb.BaseEntity `                           bson:",inline"`
	ProfileOID          primitive.ObjectID `json:"profileID,omitempty" bson:"profile_id,omitempty"`
	Gold                int32              `json:"gold"                bson:"gold"`
	Normal              int32              `json:"normal"              bson:"normal"`
}

func (Fish) IsEntity() {}

func (f *Fish) ProfileID(id string) {
	f.ProfileOID = mongodb.ToObjectID(id)
}
