package mongomodel

import (
	mongodb "tenkhours/pkg/db/mongo"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Metric struct {
	*mongodb.BaseEntity `                             bson:",inline"`
	CharacterOID        primitive.ObjectID  `json:"characterID,omitempty" bson:"character_id"`
	CategoryOID         *primitive.ObjectID `json:"categoryID,omitempty"  bson:"category_id"`
	Name                string              `json:"name,omitempty"        bson:"name"`
	Value               float64             `json:"value,omitempty"       bson:"value"`
	Unit                string              `json:"unit,omitempty"        bson:"unit"`
}

func (m *Metric) CategoryID(id *string) {
	if id != nil {
		m.CategoryOID = lo.ToPtr(mongodb.ToObjectID(*id))
	}
}

func (m *Metric) CharacterID(id string) {
	m.CharacterOID = mongodb.ToObjectID(id)
}
