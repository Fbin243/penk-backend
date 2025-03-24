package mongomodel

import (
	"time"

	mongodb "tenkhours/pkg/db/mongo"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TimeTracking struct {
	*mongodb.BaseEntity `                             bson:",inline"`
	CharacterOID        primitive.ObjectID  `json:"characterID,omitempty" bson:"character_id"`
	CategoryOID         *primitive.ObjectID `json:"categoryID,omitempty"  bson:"category_id"`
	StartTime           time.Time           `json:"startTime,omitempty"   bson:"start_time"`
	EndTime             time.Time           `json:"endTime,omitempty"     bson:"end_time"`
}

func (t *TimeTracking) CharacterID(id string) {
	t.CharacterOID = mongodb.ToObjectID(id)
}

func (t *TimeTracking) CategoryID(id *string) {
	if id != nil {
		t.CategoryOID = lo.ToPtr(mongodb.ToObjectID(*id))
	}
}
