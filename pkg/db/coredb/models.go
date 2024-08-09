package coredb

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                 primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name               string             `json:"name" bson:"name" validate:"required,min=1,max=50"`
	Email              string             `json:"email" bson:"email" validate:"required,email"`
	FirebaseUID        string             `json:"firebaseUID" bson:"firebase_uid"`
	ImageURL           string             `json:"imageURL" bson:"image_url"`
	CurrentCharacterID primitive.ObjectID `json:"currentCharacterID" bson:"current_character_id"`
	AvailableSnapshots int32              `json:"availableSnapshots" bson:"available_snapshots"`
	AutoSnapshot       bool               `json:"autoSnapshot" bson:"auto_snapshot"`
	CreatedAt          time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt          time.Time          `json:"updatedAt" bson:"updated_at"`
}

// Character
type MetricProperty struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name" validate:"required,min=1,max=50"`
	Type  string             `json:"type" bson:"type" validate:"required,min=1,max=20"`
	Value any                `json:"value" bson:"value" validate:"required"`
	Unit  string             `json:"unit,omitempty" bson:"unit,omitempty" validate:"omitempty,min=1,max=10"`
}

type MetricStyle struct {
	Color string `json:"color,omitempty" bson:"color,omitempty" validate:"omitempty,hexcolor"`
	Icon  string `json:"icon,omitempty" bson:"icon,omitempty"`
}

type CustomMetric struct {
	ID                    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name                  string             `json:"name" bson:"name" validate:"required,min=1,max=50"`
	Description           string             `json:"description" bson:"description" validate:"omitempty,max=255"`
	Time                  int32              `json:"time" bson:"time"`
	Style                 MetricStyle        `json:"style,omitempty" bson:"style,omitempty"`
	Properties            []MetricProperty   `json:"properties,omitempty" bson:"properties,omitempty" validate:"omitempty,dive"`
	LimitedPropertyNumber int32              `json:"limitedPropertyNumber" bson:"limited_property_number"`
}

type Character struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id"`
	UserID              primitive.ObjectID `json:"userID" bson:"user_id"`
	Name                string             `json:"name" bson:"name" validate:"required,min=1,max=50"`
	Gender              bool               `json:"gender" bson:"gender"`
	Avatar              string             `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Tags                []string           `json:"tags,omitempty" bson:"tags,omitempty" validate:"omitempty,tags_valid"`
	TotalFocusedTime    int32              `json:"totalFocusedTime" bson:"total_focused_time"`
	CustomMetrics       []CustomMetric     `json:"customMetrics,omitempty" bson:"custom_metrics,omitempty" validate:"omitempty,dive"`
	LimitedMetricNumber int32              `json:"limitedMetricNumber" bson:"limited_metric_number"`
}

type TimeTracking struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	CharacterID     primitive.ObjectID `json:"characterID" bson:"character_id"`
	CustomMetricID  primitive.ObjectID `json:"customMetricID,omitempty" bson:"custom_metric_id,omitempty"`
	StartTime       time.Time          `json:"startTime" bson:"start_time"`
	EndTime         time.Time          `json:"endTime,omitempty" bson:"end_time,omitempty"`
	MinDurationTime int32              `json:"minDurationTime" bson:"min_duration_time"`
	MaxDurationTime int32              `json:"maxDurationTime" bson:"max_duration_time"`
}
