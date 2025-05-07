package mongomodel

import (
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reminder struct {
	*mongodb.BaseEntity `                     bson:",inline"`
	CharacterOID        primitive.ObjectID `json:"characterID"   bson:"character_id"`
	Name                string             `bson:"name"`
	RemindTime          time.Time          `bson:"remind_time"`
	RRule               string             `bson:"rrule"`
	ReferenceOID        primitive.ObjectID `bson:"reference_id"`
	ReferenceType       entity.EntityType  `bson:"reference_type"`
}

func (r *Reminder) CharacterID(id string) {
	r.CharacterOID = mongodb.ToObjectID(id)
}

func (r *Reminder) ReferenceID(id string) {
	r.ReferenceOID = mongodb.ToObjectID(id)
}
