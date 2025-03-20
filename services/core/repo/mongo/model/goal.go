package mongorepo

import (
	"time"

	"tenkhours/services/core/entity"

	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Goal struct {
	*mongodb.BaseEntity `                   bson:",inline"`
	CharacterOID        primitive.ObjectID `json:"characterID" bson:"character_id"`
	Name                string             `json:"name"        bson:"name"`
	Description         string             `json:"description" bson:"description"`
	StartTime           time.Time          `json:"startTime"   bson:"start_time"`
	EndTime             time.Time          `json:"endTime"     bson:"end_time"`
	Status              entity.GoalStatus  `json:"status"      bson:"status"`
	Categories          []Category         `json:"categories"  bson:"categories"`
	Metrics             []GoalMetric       `json:"metrics"     bson:"metrics"`
	Checkboxes          []Checkbox         `json:"checkboxes"  bson:"checkboxes"`
}

func (g *Goal) CharacterID(id string) {
	g.CharacterOID = mongodb.ToObjectID(id)
}

type GoalMetric struct {
	*Metric     `                   bson:",inline"`
	Condition   entity.MetricCondition `json:"condition"   bson:"condition"`
	TargetValue *float64               `json:"targetValue" bson:"target_value,omitempty"`
	RangeValue  *entity.Range          `json:"rangeValue"  bson:"range_value,omitempty"`
}
