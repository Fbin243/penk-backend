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
type MetricsType interface {
}

type CustomMetric struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CharacterID primitive.ObjectID `json:"character_id,omitempty" bson:"character_id,omitempty"`
	Type        string             `json:"type,omitempty" bson:"type,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Value       string             `json:"value,omitempty" bson:"value,omitempty"`
}

type Character struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID           string             `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Name             string             `json:"name,omitempty" bson:"name,omitempty"`
	Tags             []string           `json:"tags,omitempty" bson:"tags,omitempty"`
	TotalFocusedTime int32              `json:"total_focused_time,omitempty" bson:"total_focused_time,omitempty"`
	CustomMetrics    []CustomMetric     `json:"custom_metrics,omitempty" bson:"custom_metrics,omitempty"`
}
