package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *TimeTrackingRepo) UnassignCategory(ctx context.Context, categoryID string) error {
	return r.UpdateMany(ctx,
		bson.M{"category_id": mongodb.ToObjectID(categoryID)},
		bson.M{"$set": bson.M{"category_id": nil}})
}

func (r *TimeTrackingRepo) UnassignReference(ctx context.Context, referenceID string, referenceType entity.EntityType) error {
	return r.UpdateMany(ctx,
		bson.M{
			"reference_id":   mongodb.ToObjectID(referenceID),
			"reference_type": referenceType,
		},
		bson.M{"$set": bson.M{"reference_id": nil}})
}

func (r *TimeTrackingRepo) UpdateCategoryByReferenceID(ctx context.Context, referenceID string, categoryID *string) error {
	return r.UpdateMany(ctx,
		bson.M{"reference_id": mongodb.ToObjectID(referenceID)},
		bson.M{"$set": bson.M{"category_id": mongodb.ToObjectIDOrNil(categoryID)}})
}

func (r *TimeTrackingRepo) DeleteByCharacterID(ctx context.Context, characterID string) error {
	return r.DeleteMany(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
}

func (r TimeTrackingRepo) DeleteByCharacterIDs(ctx context.Context, characterIDs []string) error {
	return r.DeleteMany(ctx, bson.M{"character_id": bson.M{"$in": mongodb.ToObjectIDs(characterIDs)}})
}

func (r TimeTrackingRepo) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": mongodb.ToObjectIDs(ids)}})
}
