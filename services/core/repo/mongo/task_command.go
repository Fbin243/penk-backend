package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *TaskRepo) DeleteByCharacterID(ctx context.Context, characterID string) error {
	return r.DeleteOne(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
}

func (r *TaskRepo) DeleteByCharacterIDs(ctx context.Context, characterIDs []string) error {
	return r.DeleteMany(ctx, bson.M{"character_id": bson.M{"$in": mongodb.ToObjectIDs(characterIDs)}})
}

func (r *TaskRepo) UnassignCategory(ctx context.Context, categoryID string) error {
	return r.UpdateMany(ctx,
		bson.M{"category_id": mongodb.ToObjectID(categoryID)},
		bson.M{"$set": bson.M{"category_id": nil}})
}
