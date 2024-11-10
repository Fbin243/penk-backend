package repo

import (
	"time"

	"tenkhours/services/core/graph/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Profile struct {
	ID                 primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name               string             `json:"name,omitempty" bson:"name,omitempty"`
	Email              string             `json:"email,omitempty" bson:"email,omitempty"`
	FirebaseUID        string             `json:"firebaseUID,omitempty" bson:"firebase_uid,omitempty"`
	ImageURL           string             `json:"imageURL,omitempty" bson:"image_url,omitempty"`
	CurrentCharacterID primitive.ObjectID `json:"currentCharacterID,omitempty" bson:"current_character_id,omitempty"`
	AvailableSnapshots int32              `json:"availableSnapshots,omitempty" bson:"available_snapshots,omitempty"`
	AutoSnapshot       bool               `json:"autoSnapshot,omitempty" bson:"auto_snapshot,omitempty"`
	CreatedAt          time.Time          `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt          time.Time          `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}

// Character
type MetricProperty struct {
	ID    primitive.ObjectID       `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string                   `json:"name,omitempty" bson:"name,omitempty"`
	Type  model.MetricPropertyType `json:"type,omitempty" bson:"type,omitempty"`
	Value string                   `json:"value,omitempty" bson:"value,omitempty"`
	Unit  string                   `json:"unit,omitempty" bson:"unit,omitempty"`
}

type MetricStyle struct {
	Color string `json:"color,omitempty" bson:"color,omitempty"`
	Icon  string `json:"icon,omitempty" bson:"icon,omitempty"`
}

type CustomMetric struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Time        int32              `json:"time" bson:"time,omitempty"`
	Style       MetricStyle        `json:"style,omitempty" bson:"style,omitempty"`
	Properties  []MetricProperty   `json:"properties,omitempty" bson:"properties,omitempty"`
	// Should be migrated later
	LimitedPropertyNumber int32 `json:"limitedPropertyNumber" bson:"limited_property_number"`
}

type Character struct {
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ProfileID        primitive.ObjectID `json:"profileID,omitempty" bson:"profile_id,omitempty"`
	Name             string             `json:"name,omitempty" bson:"name,omitempty"`
	Gender           bool               `json:"gender,omitempty" bson:"gender,omitempty"`
	Tags             []string           `json:"tags,omitempty" bson:"tags,omitempty"`
	TotalFocusedTime int32              `json:"totalFocusedTime,omitempty" bson:"total_focused_time,omitempty"`
	CustomMetrics    []CustomMetric     `json:"customMetrics,omitempty" bson:"custom_metrics,omitempty"`
	// Should be migrated later
	LimitedMetricNumber int32 `json:"limitedMetricNumber" bson:"limited_metric_number,omitempty"`
}

// Make Character satisfy the Entity interface required by federation
func (Character) IsEntity() {}
