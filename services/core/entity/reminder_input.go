package entity

import (
	"tenkhours/pkg/types"
)

type ReminderInput struct {
	ID            *string     `json:"id"`
	Name          string      `json:"name"`
	RemindTimeStr string      `json:"remindTimeStr"`
	RRule         string      `json:"rrule"`
	ReferenceID   *string     `json:"referenceId"`
	ReferenceType *EntityType `json:"referenceType"`
}

type ReminderFilter struct {
	CharacterID *string `json:"character_id"`
}

type ReminderOrderBy struct{}

type ReminderPipeline struct {
	Filter  *ReminderFilter
	OrderBy *ReminderOrderBy
	*types.Pagination
}
