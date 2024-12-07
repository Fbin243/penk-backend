package repo

import (
	"time"

	"tenkhours/services/core/graph/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Profile struct {
	ID                     primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Name                   string             `json:"name,omitempty" bson:"name"`
	Email                  string             `json:"email,omitempty" bson:"email"`
	FirebaseUID            string             `json:"firebaseUID,omitempty" bson:"firebase_uid"`
	ImageURL               string             `json:"imageURL,omitempty" bson:"image_url"`
	CurrentCharacterID     primitive.ObjectID `json:"currentCharacterID,omitempty" bson:"current_character_id"`
	AvailableSnapshots     int32              `json:"availableSnapshots,omitempty" bson:"available_snapshots`
	LimitedCharacterNumber int32              `json:"limitedCharacterNumber,omitempty" bson:"limited_character_number"`
	AutoSnapshot           bool               `json:"autoSnapshot,omitempty" bson:"auto_snapshot"`
	CreatedAt              time.Time          `json:"createdAt,omitempty" bson:"created_at"`
	UpdatedAt              time.Time          `json:"updatedAt,omitempty" bson:"updated_at"`
}

// Character
type MetricProperty struct {
	ID    primitive.ObjectID       `json:"id,omitempty" bson:"_id"`
	Name  string                   `json:"name,omitempty" bson:"name"`
	Type  model.MetricPropertyType `json:"type,omitempty" bson:"type"`
	Value string                   `json:"value,omitempty" bson:"value"`
	Unit  string                   `json:"unit,omitempty" bson:"unit"`
}

type MetricStyle struct {
	Color string `json:"color,omitempty" bson:"color"`
	Icon  string `json:"icon,omitempty" bson:"icon"`
}

type CustomMetric struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Name        string             `json:"name,omitempty" bson:"name"`
	Description string             `json:"description,omitempty" bson:"description"`
	Time        int32              `json:"time" bson:"time,omitempty"`
	Style       MetricStyle        `json:"style,omitempty" bson:"style"`
	Properties  []MetricProperty   `json:"properties,omitempty" bson:"properties"`
	// Should be migrated later
	LimitedPropertyNumber int32 `json:"limitedPropertyNumber" bson:"limited_property_number"`
}

type Character struct {
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ProfileID        primitive.ObjectID `json:"profileID,omitempty" bson:"profile_id"`
	Name             string             `json:"name,omitempty" bson:"name"`
	Gender           bool               `json:"gender,omitempty" bson:"gender"`
	Tags             []string           `json:"tags,omitempty" bson:"tags"`
	TotalFocusedTime int32              `json:"totalFocusedTime,omitempty" bson:"total_focused_time"`
	CustomMetrics    []CustomMetric     `json:"customMetrics,omitempty" bson:"custom_metrics"`
	// Should be migrated later
	LimitedMetricNumber int32 `json:"limitedMetricNumber" bson:"limited_metric_number"`
}

// Make Character satisfy the Entity interface required by federation
func (Character) IsEntity() {}
