package entity

import (
	"time"

	"tenkhours/pkg/db/base"
)

type Task struct {
	*base.BaseEntity `                     bson:",inline"`
	CharacterID      string     `json:"characterID"   bson:"character_id"`
	CategoryID       *string    `json:"categoryID"    bson:"category_id"`
	Name             string     `json:"name"          bson:"name"`
	Priority         int        `json:"priority"      bson:"priority"`
	CompletedTime    *time.Time `json:"completedTime" bson:"completed_time"`
	Subtasks         []Checkbox `json:"subtasks"      bson:"subtasks"`
	Description      *string    `json:"description"   bson:"description"`
	Deadline         *time.Time `json:"deadline"      bson:"deadline"`
}
