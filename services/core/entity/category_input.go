package entity

import "tenkhours/pkg/types"

type CategoryInput struct {
	ID          *string             `json:"id,omitempty"`
	Name        string              `json:"name"                  validate:"min=1,max=50"`
	Description *string             `json:"description,omitempty" validate:"omitempty,max=255"`
	Style       *CategoryStyleInput `json:"style"`
}

type CategoryStyleInput struct {
	Color string `json:"color" validate:"hexcolor"`
	Icon  string `json:"icon"`
}

type CategoryPipeline struct {
	Filter  *CategoryFilter
	OrderBy *CategoryOrderBy
	*types.Pagination
}

type CategoryFilter struct {
	CharacterID  *string  `json:"character_id"`
	CharacterIDs []string `json:"character_ids"`
}

type CategoryOrderBy struct{}
