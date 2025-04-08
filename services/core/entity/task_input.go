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
