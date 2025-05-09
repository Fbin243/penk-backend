package entity

import (
	core_entity "tenkhours/services/core/entity"
)

type ReminderWithMetadata struct {
	Reminder core_entity.Reminder `json:"reminder" bson:"reminder"`
	Habit    *core_entity.Habit   `json:"habit"    bson:"habit,omitempty"`
	Task     *core_entity.Task    `json:"task"     bson:"task,omitempty"`
}
