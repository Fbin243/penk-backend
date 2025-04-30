package mongomodel

import (
	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HabitLog struct {
	*mongodb.BaseEntity `                 bson:",inline"`
	Timestamp           string             `json:"timestamp" bson:"timestamp"`
	HabitOID            primitive.ObjectID `json:"habitID"   bson:"habit_id"`
	Value               float64            `json:"value"     bson:"value"`
}

func (h *HabitLog) HabitID(id string) {
	h.HabitOID = mongodb.ToObjectID(id)
}
