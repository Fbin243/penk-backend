package mongorepo

import (
	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Metric struct {
	*mongodb.BaseEntity `                             bson:",inline"`
	CharacterOID        primitive.ObjectID `json:"characterID,omitempty" bson:"character_id"`
	Name                string             `json:"name,omitempty"        bson:"name"`
	Value               float64            `json:"value,omitempty"       bson:"value"`
	Unit                string             `json:"unit,omitempty"        bson:"unit"`
}

func (m *Metric) CharacterID(id string) {
	m.CharacterOID = mongodb.ToObjectID(id)
}
