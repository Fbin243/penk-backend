package entity

import "tenkhours/pkg/types"

type MetricInput struct {
	ID         *string `json:"id,omitempty"`
	CategoryID *string `json:"categoryID,omitempty"`
	Name       string  `json:"name"                 validate:"min=1,max=50"`
	Value      float64 `json:"value"`
	Unit       string  `json:"unit"                 validate:"omitempty,max=50"`
}

type MetricPipeline struct {
	Filter  *MetricFilter
	OrderBy *MetricOrderBy
	*types.Pagination
}

type MetricFilter struct {
	CharacterID *string `json:"character_id"`
}

type MetricOrderBy struct{}
