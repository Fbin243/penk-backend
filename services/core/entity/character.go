package entity

import "tenkhours/pkg/db/base"

type Character struct {
	*base.BaseEntity `                           bson:",inline"`
	ProfileID        string `json:"profileID,omitempty" bson:"profile_id"`
	Name             string `json:"name,omitempty"      bson:"name"`
}
