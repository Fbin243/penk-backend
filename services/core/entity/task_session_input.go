package entity

import "time"

type TaskSessionInput struct {
	ID            *string    `json:"id,omitempty"`
	TaskID        string     `json:"taskID"`
	StartTime     time.Time  `json:"startTime"`
	EndTime       time.Time  `json:"endTime"`
	CompletedTime *time.Time `json:"completedTime"`
}

type TaskSessionPineline struct {
	Filter  *TaskSessionFilter
	OrderBy *TaskSessionOrderBy
	Limit   *int
	Offset  *int
}

type TaskSessionFilter struct {
	TaskIDs     []string   `json:"task_ids"`
	TaskID      *string    `json:"task_id"`
	StartTime   *time.Time `json:"startTime"`
	EndTime     *time.Time `json:"endTime"`
	IsCompleted *bool      `json:"isCompleted"`
}

type TaskSessionOrderBy struct{}
