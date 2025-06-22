package entity

import (
	"time"

	"tenkhours/pkg/types"
)

type TaskInput struct {
	ID            *string         `json:"id,omitempty"`
	CategoryID    *string         `json:"categoryID,omitempty"`
	Name          string          `json:"name"`
	Priority      int             `json:"priority"`
	Subtasks      []CheckboxInput `json:"subtasks"`
	Description   *string         `json:"description"`
	CompletedTime *time.Time      `json:"completedTime"`
	Deadline      *time.Time      `json:"deadline"`
}

type TaskPipeline struct {
	Filter  *TaskFilter
	OrderBy *TaskOrderBy
	*types.Pagination
}

type TaskFilter struct {
	CharacterID  *string  `json:"character_id"`
	CharacterIDs []string `json:"character_ids"`
	IsCompleted  *bool    `json:"is_completed"`
}

type TaskOrderBy struct {
	Priority *types.SortOrder `json:"priority"`
}
