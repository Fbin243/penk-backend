package mongodb

import (
	"tenkhours/pkg/db/base"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GenObjectID() string {
	return primitive.NewObjectID().Hex()
}

func ToObjectID(id string) primitive.ObjectID {
	oid, _ := primitive.ObjectIDFromHex(id)
	return oid
}

func ToObjectIDs(ids []string) []primitive.ObjectID {
	return lo.Map(ids, func(id string, _ int) primitive.ObjectID {
		return ToObjectID(id)
	})
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
