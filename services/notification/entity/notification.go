package entity

import (
	core_entity "tenkhours/services/core/entity"
)

// NotificationMessage represents a message to be sent to Kafka for push notification
type NotificationMessage struct {
	// User identification
	CharacterID string `json:"character_id"`

	// Notification content
	Name     string `json:"title"`
	Body     string `json:"body"`
	Priority string `json:"priority"` // high, normal, low

	// Reference information
	ReferenceID   *string                 `json:"reference_id,omitempty"`
	ReferenceType *core_entity.EntityType `json:"reference_type,omitempty"`

	// Additional data for deep linking
	Data map[string]interface{} `json:"data,omitempty"`
}
