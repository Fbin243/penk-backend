package entity

import "tenkhours/pkg/db/base"

type Category struct {
	*base.BaseEntity `                             bson:",inline"`
	CharacterID      string        `json:"characterID,omitempty" bson:"character_id"`
	Name             string        `json:"name,omitempty"        bson:"name"`
	Description      string        `json:"description,omitempty" bson:"description"`
	Style            CategoryStyle `json:"style,omitempty"       bson:"style"`
}

type CategoryStyle struct {
	Color string `json:"color,omitempty" bson:"color"`
	Icon  string `json:"icon,omitempty"  bson:"icon"`
}
