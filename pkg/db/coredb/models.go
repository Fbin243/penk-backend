package coredb

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MetricPropertyType string

const (
	STRING MetricPropertyType = "STRING"
	NUMBER MetricPropertyType = "NUMBER"
)

type Profile struct {
	ID                 primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name               string             `json:"name" bson:"name" validate:"required,min=1,max=50"`
	Email              string             `json:"email" bson:"email" validate:"required,email"`
	FirebaseUID        string             `json:"firebaseUID" bson:"firebase_uid"`
	ImageURL           string             `json:"imageURL" bson:"image_url"`
	CurrentCharacterID primitive.ObjectID `json:"currentCharacterID,omitempty" bson:"current_character_id,omitempty"`
	AvailableSnapshots int32              `json:"availableSnapshots" bson:"available_snapshots"`
	AutoSnapshot       bool               `json:"autoSnapshot" bson:"auto_snapshot"`
	CreatedAt          time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt          time.Time          `json:"updatedAt" bson:"updated_at"`
}

// Character
type MetricProperty struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name" validate:"required,min=1,max=50"`
	Type  MetricPropertyType `json:"type" bson:"type" validate:"required,min=1,max=20"`
	Value string             `json:"value" bson:"value" validate:"required"`
	Unit  string             `json:"unit,omitempty" bson:"unit,omitempty" validate:"omitempty,min=1,max=10"`
}

type MetricStyle struct {
	Color string `json:"color,omitempty" bson:"color,omitempty" validate:"omitempty,hexcolor"`
	Icon  string `json:"icon,omitempty" bson:"icon,omitempty"`
}

type CustomMetric struct {
	ID                    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name                  string             `json:"name" bson:"name" validate:"required,min=1,max=50"`
	Description           string             `json:"description,omitempty" bson:"description" validate:"omitempty,max=255"`
	Time                  int32              `json:"time" bson:"time"`
	Style                 MetricStyle        `json:"style,omitempty" bson:"style,omitempty"`
	Properties            []MetricProperty   `json:"properties,omitempty" bson:"properties,omitempty" validate:"omitempty,dive"`
	LimitedPropertyNumber int32              `json:"limitedPropertyNumber" bson:"limited_property_number"`
}

type Character struct {
	ID                  primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ProfileID           primitive.ObjectID `json:"profileID,omitempty" bson:"profile_id,omitempty"`
	Name                string             `json:"name" bson:"name" validate:"required,min=1,max=50"`
	Gender              bool               `json:"gender" bson:"gender"`
	Tags                []string           `json:"tags,omitempty" bson:"tags,omitempty" validate:"omitempty,tags_valid"`
	TotalFocusedTime    int32              `json:"totalFocusedTime" bson:"total_focused_time"`
	CustomMetrics       []CustomMetric     `json:"customMetrics,omitempty" bson:"custom_metrics,omitempty" validate:"omitempty,dive"`
	LimitedMetricNumber int32              `json:"limitedMetricNumber" bson:"limited_metric_number"`
}

// Make Character satisfy the Entity interface required by federation
func (Character) IsEntity() {}
