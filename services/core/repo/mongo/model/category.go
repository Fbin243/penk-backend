package mongomodel

import (
	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	*mongodb.BaseEntity `                             bson:",inline"`
	CharacterOID        primitive.ObjectID `json:"characterID,omitempty" bson:"character_id"`
	Name                string             `json:"name,omitempty"        bson:"name"`
	Description         string             `json:"description,omitempty" bson:"description"`
	Style               CategoryStyle      `json:"style,omitempty"       bson:"style"`
}

func (c *Category) CharacterID(id string) {
	c.CharacterOID = mongodb.ToObjectID(id)
}

type CategoryStyle struct {
	Color string `json:"color,omitempty" bson:"color"`
	Icon  string `json:"icon,omitempty"  bson:"icon"`
}
