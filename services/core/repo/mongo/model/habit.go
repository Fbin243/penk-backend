package mongomodel

import (
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Habit struct {
	*mongodb.BaseEntity `                       bson:",inline"`
	CharacterOID        primitive.ObjectID    `json:"characterID"     bson:"character_id"`
	CategoryOID         *primitive.ObjectID   `json:"categoryID"      bson:"category_id"`
	CompletionType      entity.CompletionType `json:"completionType"  bson:"completion_type"`
	Name                string                `json:"name"            bson:"name"`
	Value               float64               `json:"value,omitempty" bson:"value"`
	Unit                *string               `json:"unit,omitempty"  bson:"unit"`
	RRule               string                `json:"rrule"           bson:"rrule"`
	ResetDuration       entity.HabitReset     `json:"resetDuration"   bson:"reset_duration"`
}

func (h *Habit) CharacterID(id string) {
	h.CharacterOID = mongodb.ToObjectID(id)
}

func (h *Habit) CategoryID(id *string) {
	if id != nil {
		h.CategoryOID = lo.ToPtr(mongodb.ToObjectID(*id))
	}
}
