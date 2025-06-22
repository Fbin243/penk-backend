package entity

import "tenkhours/pkg/types"

type HabitInput struct {
	ID             *string        `json:"id,omitempty"`
	CategoryID     *string        `json:"categoryID"`
	CompletionType CompletionType `json:"completionType"`
	Name           string         `json:"name"`
	Value          float64        `json:"value"`
	Unit           *string        `json:"unit"`
	RRule          string         `json:"rrule"`
	ResetDuration  HabitReset     `json:"resetDuration"`
}

type HabitPipeline struct {
	Filter  *HabitFilter
	OrderBy *HabitOrderBy
	*types.Pagination
}

type HabitFilter struct {
	CharacterID  *string  `json:"character_id"`
	CharacterIDs []string `json:"character_ids"`
}

type HabitOrderBy struct{}
