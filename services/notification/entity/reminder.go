package entity

import (
	"tenkhours/pkg/db/base"
	"time"
)

type Reminder struct {
	*base.BaseEntity `                            bson:",inline"`
	ProfileID        string       `json:"profile_id,omitempty" bson:"profile_id"`
	Type             ReminderType `json:"type,omitempty"       bson:"type"`
	Title            string       `json:"title,omitempty"      bson:"title"`
	RemindTime       time.Time    `json:"remind_time,omitempty" bson:"remind_time"`
	Recurrence       string       `json:"recurrence,omitempty"  bson:"recurrence"`
	LinkedItemID     string       `json:"linked_item_id,omitempty" bson:"linked_item_id"`
}
