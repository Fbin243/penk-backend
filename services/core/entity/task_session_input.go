package entity

import "time"

type TaskSessionInput struct {
	ID            *string    `json:"id,omitempty"`
	TaskID        string     `json:"taskID"`
	StartTime     time.Time  `json:"startTime"`
	EndTime       time.Time  `json:"endTime"`
	CompletedTime *time.Time `json:"completedTime"`
}
