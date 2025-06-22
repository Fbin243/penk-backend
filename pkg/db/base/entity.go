package base

import (
	"time"
)

// Interface of base entity to be embedded to concrete entity structs
type IBaseEntity interface {
	GetID() string
	SetID(id string)
	GetCreatedAt() time.Time
	SetCreatedAtByNow()
	GetUpdatedAt() time.Time
	SetUpdatedAtByNow()
}

type BaseEntity struct {
	ID        string    `json:"id,omitempty"        bson:"_id,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}

func (m *BaseEntity) GetID() string {
	return m.ID
}

func (m *BaseEntity) SetID(id string) {
	m.ID = id
}

func (m *BaseEntity) GetCreatedAt() time.Time {
	return m.CreatedAt
}

func (m *BaseEntity) SetCreatedAtByNow() {
	m.CreatedAt = time.Now()
}

func (m *BaseEntity) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}

func (m *BaseEntity) SetUpdatedAtByNow() {
	m.UpdatedAt = time.Now()
}
