package mongomodel

import (
	"time"

	mongodb "tenkhours/pkg/db/mongo"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	*mongodb.BaseEntity `                     bson:",inline"`
	CharacterOID        primitive.ObjectID  `json:"characterID"   bson:"character_id"`
	CategoryOID         *primitive.ObjectID `json:"categoryID"    bson:"category_id"`
	Name                string              `json:"name"          bson:"name"`
	Priority            int                 `json:"priority"      bson:"priority"`
	CompletedTime       *time.Time          `json:"completedTime" bson:"completed_time"`
	Subtasks            []Checkbox          `json:"subtasks"      bson:"subtasks"`
	Description         *string             `json:"description"   bson:"description"`
	Deadline            *time.Time          `json:"deadline"      bson:"deadline"`
}

func (t *Task) CharacterID(id string) {
	t.CharacterOID = mongodb.ToObjectID(id)
}

func (t *Task) CategoryID(id *string) {
	if id != nil {
		t.CategoryOID = lo.ToPtr(mongodb.ToObjectID(*id))
	}
}
