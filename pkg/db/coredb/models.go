package coredb

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Email       string             `json:"email" bson:"email"`
	FirebaseUID string             `json:"firebaseUID" bson:"firebase_uid"`
	ImageURL    string             `json:"imageURL" bson:"image_url"`
	CreatedAt   time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updated_at"`
}

// Character
type MetricProperty struct {
	Name  string `json:"name" bson:"name"`
	Type  string `json:"type" bson:"type"`
	Value any    `json:"value" bson:"value"`
	Unit  string `json:"unit" bson:"unit"`
}

type StyleType struct {
	Color string `json:"color" bson:"color"`
	Icon  string `json:"icon" bson:"icon"`
}

type CustomMetric struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CharacterID string             `json:"character_id,omitempty" bson:"character_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Time        int32              `json:"time" bson:"time"` // will change after we merge time service
	Style       StyleType          `json:"style" bson:"style"`
	Properties  []MetricProperty   `json:"properties" bson:"properties"`
}

type Character struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID           string             `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Name             string             `json:"name" bson:"name"`
	Tags             []string           `json:"tags" bson:"tags"`
	TotalFocusedTime int32              `json:"total_focused_time" bson:"total_focused_time"`
	CustomMetrics    []CustomMetric     `json:"custom_metrics" bson:"custom_metrics"`
}
