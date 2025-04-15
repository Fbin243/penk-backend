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
	Timestamp        time.Time   `json:"timestamp,omitempty"     bson:"timestamp"`
	Duration         int         `json:"duration,omitempty"      bson:"duration"`
}
