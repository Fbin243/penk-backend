package coredb

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id"`
	Name               string             `json:"name" bson:"name"`
	Email              string             `json:"email" bson:"email"`
	FirebaseUID        string             `json:"firebaseUID" bson:"firebase_uid"`
	ImageURL           string             `json:"imageURL" bson:"image_url"`
	CurrentCharacterID primitive.ObjectID `json:"currentCharacterID" bson:"current_character_id"`
	CreatedAt          time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt          time.Time          `json:"updatedAt" bson:"updated_at"`
}

// Character
type MetricProperty struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Type  string             `json:"type" bson:"type"`
	Value any                `json:"value" bson:"value"`
	Unit  string             `json:"unit" bson:"unit"`
}

type MetricStyle struct {
	Color string `json:"color" bson:"color"`
	Icon  string `json:"icon" bson:"icon"`
}

type CustomMetric struct {
	ID                    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name                  string             `json:"name" bson:"name"`
	Description           string             `json:"description" bson:"description"`
	Time                  int32              `json:"time" bson:"time"`
	Style                 MetricStyle        `json:"style" bson:"style"`
	Properties            []MetricProperty   `json:"properties" bson:"properties"`
	LimitedPropertyNumber int32              `json:"limitedPropertyNumber" bson:"limited_property_number"`
}

type Character struct {
	ID                  primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID              primitive.ObjectID `json:"userID,omitempty" bson:"user_id,omitempty"`
	Name                string             `json:"name" bson:"name"`
	Tags                []string           `json:"tags" bson:"tags"`
	TotalFocusedTime    int32              `json:"totalFocusedTime" bson:"total_focused_time"`
	CustomMetrics       []CustomMetric     `json:"customMetrics" bson:"custom_metrics"`
	LimitedMetricNumber int32              `json:"limitedMetricNumber" bson:"limited_metric_number"`
}

type TimeTracking struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CharacterID     primitive.ObjectID `json:"characterID" bson:"character_id"`
	CustomMetricID  primitive.ObjectID `json:"customMetricID" bson:"custom_metric_id"`
	StartTime       time.Time          `json:"startTime" bson:"start_time"`
	EndTime         time.Time          `json:"endTime" bson:"end_time"`
	MinDurationTime int32              `json:"minDurationTime" bson:"min_duration_time"`
	MaxDurationTime int32              `json:"maxDurationTime" bson:"max_duration_time"`
}
