package mongodb

import (
	"tenkhours/pkg/db/base"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GenObjectID() string {
	return primitive.NewObjectID().Hex()
}

func ToObjectID(id string) primitive.ObjectID {
	oid, _ := primitive.ObjectIDFromHex(id)
	return oid
}

func addMissingFields[M base.IBaseEntity](m M) {
	if m.GetID() == "" {
		m.SetID(primitive.NewObjectID().Hex())
	}
	if m.GetCreatedAt().IsZero() {
		m.SetCreatedAtByNow()
	}
	if m.GetUpdatedAt().IsZero() {
		m.SetUpdatedAtByNow()
	}
}
