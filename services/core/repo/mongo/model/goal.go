package mongomodel

import (
	"time"

	"tenkhours/services/core/entity"

	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Goal struct {
	*mongodb.BaseEntity `                     bson:",inline"`
	CharacterOID        primitive.ObjectID `json:"characterID"   bson:"character_id"`
	Name                string             `json:"name"          bson:"name"`
	Description         string             `json:"description"   bson:"description"`
	StartTime           time.Time          `json:"startTime"     bson:"start_time"`
	EndTime             time.Time          `json:"endTime"       bson:"end_time"`
	Metrics             []GoalMetric       `json:"metrics"       bson:"metrics"`
	Checkboxes          []Checkbox         `json:"checkboxes"    bson:"checkboxes"`
	CompletedTime       *time.Time         `json:"completedTime" bson:"completed_time"`
}

func (g *Goal) CharacterID(id string) {
	g.CharacterOID = mongodb.ToObjectID(id)
}

type GoalMetric struct {
	OID         primitive.ObjectID     `json:"id"          bson:"id"`
	Condition   entity.MetricCondition `json:"condition"   bson:"condition"`
	TargetValue *float64               `json:"targetValue" bson:"target_value"`
	RangeValue  *entity.Range          `json:"rangeValue"  bson:"range_value"`
}

func (m *GoalMetric) ID(id string) {
	m.OID = mongodb.ToObjectID(id)
}
