package entity

import (
	"time"

	"tenkhours/pkg/db/base"
)

type Reminder struct {
	*base.BaseEntity `bson:",inline"`
	CharacterID      string      `bson:"character_id"`
	Name             string      `bson:"name"`
	RemindTime       time.Time   `bson:"remind_time"`
	RRule            string      `bson:"rrule"`
	ReferenceID      *string     `bson:"reference_id"`
	ReferenceType    *EntityType `bson:"reference_type"`
}
