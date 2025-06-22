package entity

import "tenkhours/pkg/db/base"

type Metric struct {
	*base.BaseEntity `                             bson:",inline"`
	CharacterID      string  `json:"characterID,omitempty" bson:"character_id"`
	CategoryID       *string `json:"categoryID,omitempty"  bson:"category_id"`
	Name             string  `json:"name,omitempty"        bson:"name"`
	Value            float64 `json:"value,omitempty"       bson:"value"`
	Unit             string  `json:"unit,omitempty"        bson:"unit"`
}
