package mongomodel

import (
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reminder struct {
	*mongodb.BaseEntity `                     bson:",inline"`
	CharacterOID        primitive.ObjectID  `json:"characterID"   bson:"character_id"`
	Name                string              `json:"name"          bson:"name"`
	RemindTime          *time.Time          `json:"remindTime"    bson:"remind_time"`
	RemindTimeStr       string              `json:"remindTimeStr" bson:"remind_time_str"`
	RRule               string              `json:"rrule"         bson:"rrule"`
	ReferenceOID        *primitive.ObjectID `json:"referenceID"   bson:"reference_id"`
	ReferenceType       *entity.EntityType  `json:"referenceType" bson:"reference_type"`
}

func (r *Reminder) CharacterID(id string) {
	r.CharacterOID = mongodb.ToObjectID(id)
}

func (r *Reminder) ReferenceID(id *string) {
	if id != nil {
		r.ReferenceOID = lo.ToPtr(mongodb.ToObjectID(*id))
	}
}
