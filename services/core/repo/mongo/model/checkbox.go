package mongomodel

import (
	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Checkbox struct {
	OID   primitive.ObjectID `json:"id"    bson:"id"`
	Name  string             `json:"name"  bson:"name"`
	Value bool               `json:"value" bson:"value"`
}

func (c *Checkbox) ID(id string) {
	c.OID = mongodb.ToObjectID(id)
}
