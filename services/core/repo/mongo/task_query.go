package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *TaskRepo) CountByCharacterID(ctx context.Context, characterID string) (int, error) {
	return r.Count(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
}

func (r *TaskRepo) CountByCategoryID(ctx context.Context, categoryID string) (int, error) {
	return r.Count(ctx, bson.M{"category_id": mongodb.ToObjectID(categoryID)})
}

func (r *TaskRepo) CountUnassigned(ctx context.Context, characterID string) (int, error) {
	return r.Count(ctx, bson.M{
		"character_id": mongodb.ToObjectID(characterID),
		"category_id":  nil,
	})
}

func (r *TaskRepo) Exist(ctx context.Context, characterID, taskID string) error {
	return r.Exists(ctx, bson.M{
		"_id":          mongodb.ToObjectID(taskID),
		"character_id": mongodb.ToObjectID(characterID),
	})
}

func (r *TaskRepo) FindByCharacterID(ctx context.Context, characterID string) ([]entity.Task, error) {
	return r.FindMany(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
}

func (r *TaskRepo) FindByCharacterIDs(ctx context.Context, characterIDs []string) ([]entity.Task, error) {
	return r.FindMany(ctx, bson.M{"character_id": bson.M{"$in": mongodb.ToObjectIDs(characterIDs)}})
}
