package entity

import (
	"tenkhours/pkg/db/base"
)

type Fish struct {
	*base.BaseEntity `                           bson:",inline"`
	ProfileID        string `json:"profileID,omitempty" bson:"profile_id,omitempty"`
	Gold             int32  `json:"gold"                bson:"gold"`
	Normal           int32  `json:"normal"              bson:"normal"`
}

func (Fish) IsEntity() {}

type FishType string

const (
	FishTypeNone   FishType = "None"
	FishTypeGold   FishType = "Gold"
	FishTypeNormal FishType = "Normal"
)

type CatchFishResult struct {
	FishType FishType `json:"fishType"`
	Number   int32    `json:"number"`
}
