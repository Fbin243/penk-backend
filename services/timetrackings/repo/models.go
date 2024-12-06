package repo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TimeTracking struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	CharacterID    primitive.ObjectID `json:"characterID,omitempty" bson:"character_id"`
	CustomMetricID primitive.ObjectID `json:"customMetricID,omitempty" bson:"custom_metric_id,omitempty"`
	StartTime      time.Time          `json:"startTime,omitempty" bson:"start_time"`
	EndTime        time.Time          `json:"endTime,omitempty" bson:"end_time,omitempty"`
	// Will be move to the better place later on
	MinDurationTime int32 `json:"minDurationTime" bson:"min_duration_time"`
	MaxDurationTime int32 `json:"maxDurationTime" bson:"max_duration_time"`
}
