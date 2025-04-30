package mongomodel

import (
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TimeTracking struct {
	*mongodb.BaseEntity `                               bson:",inline"`
	CharacterOID        primitive.ObjectID  `json:"characterID,omitempty"   bson:"character_id"`
	CategoryOID         *primitive.ObjectID `json:"categoryID,omitempty"    bson:"category_id"`
	ReferenceOID        *primitive.ObjectID `json:"referenceID,omitempty"   bson:"reference_id"`
	ReferenceType       *entity.EntityType  `json:"referenceType,omitempty" bson:"reference_type"`
	Timestamp           time.Time           `json:"timestamp"               bson:"timestamp"`
	Duration            int                 `json:"duration"                bson:"duration"`
}

func (t *TimeTracking) CharacterID(id string) {
	t.CharacterOID = mongodb.ToObjectID(id)
}

func (t *TimeTracking) CategoryID(id *string) {
	if id != nil {
		t.CategoryOID = lo.ToPtr(mongodb.ToObjectID(*id))
	}
}

func (t *TimeTracking) ReferenceID(id *string) {
	if id != nil {
		t.ReferenceOID = lo.ToPtr(mongodb.ToObjectID(*id))
	}
}
