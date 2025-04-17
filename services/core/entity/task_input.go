package entity

import "time"

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

type TaskPineline struct {
	Filter  *TaskFilter
	OrderBy *TaskOrderBy
	Limit   *int
	Offset  *int
}

type TaskFilter struct {
	CharacterID  *string  `json:"character_id"`
	CharacterIDs []string `json:"character_ids"`
	IsCompleted  *bool    `json:"is_completed"`
}

type TaskOrderBy struct{}
