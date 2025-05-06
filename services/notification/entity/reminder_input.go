package entity

import "time"

type ReminderInput struct {
	Type         ReminderType `json:"type,omitempty"       validate:"required"`
	Title        string       `json:"title,omitempty"      validate:"required"`
	RemindTime   time.Time    `json:"remind_time,omitempty" validate:"required"`
	Recurrence   string       `json:"recurrence,omitempty"  validate:"required"`
	LinkedItemID string       `json:"linked_item_id,omitempty" validate:"required"`
}

type ReminderType string

const (
	ReminderTypeTask  ReminderType = "TASK"
	ReminderTypeEvent ReminderType = "EVENT"
	ReminderTypeHabit ReminderType = "HABIT"
)

type ReminderStatus string

const (
	ReminderStatusPending   ReminderStatus = "PENDING"
	ReminderStatusCompleted ReminderStatus = "COMPLETED"
	ReminderStatusCanceled  ReminderStatus = "CANCELED"
)
