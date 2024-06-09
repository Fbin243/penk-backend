package character

import "go.mongodb.org/mongo-driver/bson/primitive"

type MetricsType interface {
}

type CustomMetricData struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CharacterID primitive.ObjectID `json:"character_id,omitempty" bson:"character_id,omitempty"`
	Type        string             `json:"type,omitempty" bson:"type,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Value       string             `json:"value,omitempty" bson:"value,omitempty"`
}

type CharacterData struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID           string             `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Name             string             `json:"name,omitempty" bson:"name,omitempty"`
	Tags             []string           `json:"tags,omitempty" bson:"tags,omitempty"`
	TotalFocusedTime float64            `json:"total_focused_time,omitempty" bson:"total_focused_time,omitempty"`
	CustomMetrics    []CustomMetricData `json:"custom_metrics,omitempty" bson:"custom_metrics,omitempty"`
}
