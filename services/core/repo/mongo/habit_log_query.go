package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *HabitLogRepo) FindByHabitID(ctx context.Context, habitID string) ([]entity.HabitLog, error) {
	return r.FindMany(ctx, bson.M{"habit_id": mongodb.ToObjectID(habitID)})
}

func (r *HabitLogRepo) FindByHabitIDs(ctx context.Context, habitIDs []string) ([]entity.HabitLog, error) {
	return r.FindMany(ctx, bson.M{"habit_id": bson.M{"$in": mongodb.ToObjectIDs(habitIDs)}})
}
