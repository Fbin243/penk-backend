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

func ToObjectIDOrNil(id *string) *primitive.ObjectID {
	if id == nil {
		return nil
	}

	return lo.ToPtr(ToObjectID(*id))
}

func ToObjectIDs(ids []string) []primitive.ObjectID {
	return lo.Map(ids, func(id string, _ int) primitive.ObjectID {
		return ToObjectID(id)
	})
}

func addMissingFields[M base.IBaseEntity](m M, withTimestamp bool) {
	if m.GetID() == "" {
		m.SetID(GenObjectID())
	}
	if withTimestamp {
		if m.GetCreatedAt().IsZero() {
			m.SetCreatedAtByNow()
		}
		if m.GetUpdatedAt().IsZero() {
			m.SetUpdatedAtByNow()
		}
	}
}
