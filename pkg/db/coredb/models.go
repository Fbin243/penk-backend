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
<<<<<<< HEAD

// Character
type MetricsType interface {
}

type CustomMetric struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CharacterID primitive.ObjectID `json:"character_id,omitempty" bson:"character_id,omitempty"`
	Type        string             `json:"type" bson:"type"`
	Name        string             `json:"name" bson:"name"`
	Value       string             `json:"value" bson:"value"`
}

type Character struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID           string             `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Name             string             `json:"name" bson:"name"`
	Tags             []string           `json:"tags" bson:"tags"`
	TotalFocusedTime int32              `json:"total_focused_time" bson:"total_focused_time"`
	CustomMetrics    []CustomMetric     `json:"custom_metricsy" bson:"custom_metrics"`
}
=======
>>>>>>> dev
