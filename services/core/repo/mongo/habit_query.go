package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *HabitRepo) CountByCharacterID(ctx context.Context, characterID string) (int, error) {
	return r.Count(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
}

func (r *HabitRepo) CountByCategoryID(ctx context.Context, categoryID string) (int, error) {
	return r.Count(ctx, bson.M{"category_id": mongodb.ToObjectID(categoryID)})
}

func (r *HabitRepo) CountUnassigned(ctx context.Context, characterID string) (int, error) {
	return r.Count(ctx, bson.M{
		"character_id": mongodb.ToObjectID(characterID),
		"category_id":  nil,
	})
}

func (r *HabitRepo) Exist(ctx context.Context, characterID, habitID string) error {
	return r.Exists(ctx, bson.M{
		"_id":          mongodb.ToObjectID(habitID),
		"character_id": mongodb.ToObjectID(characterID),
	})
}

func (r *HabitRepo) FindByCharacterID(ctx context.Context, characterID string) ([]entity.Habit, error) {
	return r.FindMany(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
}

func (r *HabitRepo) FindByCharacterIDs(ctx context.Context, characterIDs []string) ([]entity.Habit, error) {
	return r.FindMany(ctx, bson.M{"character_id": bson.M{"$in": mongodb.ToObjectIDs(characterIDs)}})
}
