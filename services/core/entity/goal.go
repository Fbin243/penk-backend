package entity

import (
	"time"

	"tenkhours/pkg/db/base"
)

type Goal struct {
	*base.BaseEntity `                   bson:",inline"`
	CharacterID      string           `json:"characterID" bson:"character_id"`
	Name             string           `json:"name"        bson:"name"`
	Description      string           `json:"description" bson:"description"`
	StartDate        time.Time        `json:"startDate"   bson:"start_date"`
	EndDate          time.Time        `json:"endDate"     bson:"end_date"`
	Status           GoalFinishStatus `json:"status"      bson:"status"`
	Target           []Category       `json:"target"      bson:"target"`
}

type (
	GoalFinishStatus string
	GoalExpireStatus string
)

const (
	GoalFinishStatusFinished   GoalFinishStatus = "Finished"
	GoalFinishStatusUnfinished GoalFinishStatus = "Unfinished"
	GoalExpireStatusExpired    GoalExpireStatus = "Expired"
	GoalExpireStatusUnexpired  GoalExpireStatus = "Unexpired"
)

type GoalStatusFilter struct {
	FinishStatus *GoalFinishStatus
	ExpireStatus *GoalExpireStatus
}
