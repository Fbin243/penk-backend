package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/types"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *HabitLogRepo) FindByHabitID(ctx context.Context, habitID string, timeFilter *types.TimeFilter) ([]entity.HabitLog, error) {
	return r.FindMany(ctx, bson.M{
		"habit_id":  mongodb.ToObjectID(habitID),
		"timestamp": mongodb.MakeTimeFilter(timeFilter),
	})
}

func (r *HabitLogRepo) FindByHabitIDs(ctx context.Context, habitIDs []string, timeFilter *types.TimeFilter) ([]entity.HabitLog, error) {
	return r.FindMany(ctx, bson.M{
		"habit_id":  bson.M{"$in": mongodb.ToObjectIDs(habitIDs)},
		"timestamp": mongodb.MakeTimeFilter(timeFilter),
	})
}
