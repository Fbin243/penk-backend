package mongomodel

import (
	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Character struct {
	*mongodb.BaseEntity `                           bson:",inline"`
	ProfileOID          primitive.ObjectID `json:"profileID,omitempty" bson:"profile_id"`
	Name                string             `json:"name,omitempty"      bson:"name"`
}

func (p *Character) ProfileID(id string) {
	p.ProfileOID = mongodb.ToObjectID(id)
}
