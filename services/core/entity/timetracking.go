package entity

import (
	"time"

	"tenkhours/pkg/db/base"
)

type TimeTracking struct {
	*base.BaseEntity `                               bson:",inline"`
	CharacterID      string      `json:"characterID,omitempty"   bson:"character_id"`
	CategoryID       *string     `json:"categoryID,omitempty"    bson:"category_id,omitempty"`
	ReferenceID      *string     `json:"referenceID,omitempty"   bson:"reference_id,omitempty"`
	ReferenceType    *EntityType `json:"referenceType,omitempty" bson:"reference_type,omitempty"`
	StartTime        time.Time   `json:"startTime,omitempty"     bson:"start_time"`
	EndTime          time.Time   `json:"endTime,omitempty"       bson:"end_time,omitempty"`
}
