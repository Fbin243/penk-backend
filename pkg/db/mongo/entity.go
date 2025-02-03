package mongodb

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mongo base entity
type BaseEntity struct {
	OID       primitive.ObjectID `json:"id"                  bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}

func (m *BaseEntity) ID(id string) {
	m.OID = ToObjectID(id)
}
