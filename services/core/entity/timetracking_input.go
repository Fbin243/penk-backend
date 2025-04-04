package entity

import "time"

type TimeTrackingInput struct {
	CategoryID    *string     `json:"categoryID,omitempty"`
	ReferenceID   *string     `json:"referenceID,omitempty"`
	ReferenceType *EntityType `json:"referenceType,omitempty"`
	StartTime     time.Time   `json:"startTime"`
}
