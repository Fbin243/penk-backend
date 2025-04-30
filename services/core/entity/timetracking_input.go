package entity

import (
	"time"
)

type TimeTrackingInput struct {
	ReferenceID   string     `json:"referenceID,omitempty"`
	ReferenceType EntityType `json:"referenceType,omitempty"`
	Timestamp     time.Time  `json:"timestamp,omitempty"`
	Duration      int        `json:"duration,omitempty"`
}
