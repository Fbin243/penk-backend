package entity

import (
	"time"

	"tenkhours/pkg/db/base"
)

type Reminder struct {
	*base.BaseEntity `json:",inline"         bson:",inline"`
	CharacterID      string      `json:"character_id"    bson:"character_id"`
	Name             string      `json:"name"            bson:"name"`
	RemindTime       *time.Time  `json:"remind_time"     bson:"remind_time"`
	RemindTimeStr    string      `json:"remind_time_str" bson:"remind_time_str"`
	RRule            string      `json:"rrule"           bson:"rrule"`
	ReferenceID      *string     `json:"reference_id"    bson:"reference_id"`
	ReferenceType    *EntityType `json:"reference_type"  bson:"reference_type"`
}
