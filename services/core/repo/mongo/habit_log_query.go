package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/types"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *HabitLogRepo) FindByHabitID(ctx context.Context, habitID string, timeFilter *types.TimeFilter) ([]entity.HabitLog, error) {
	filter := bson.M{"habit_id": mongodb.ToObjectID(habitID)}
	if timeFilter != nil {
		if timeFilter.StartTime != nil {
			filter["timestamp"] = bson.M{"$gte": *timeFilter.StartTime}
		}
		if timeFilter.EndTime != nil {
			filter["timestamp"] = bson.M{"$lte": *timeFilter.EndTime}
		}
	}

	return r.FindMany(ctx, filter)
}

func (r *HabitLogRepo) FindByHabitIDs(ctx context.Context, habitIDs []string, timeFilter *types.TimeFilter) ([]entity.HabitLog, error) {
	filter := bson.M{"habit_id": bson.M{"$in": mongodb.ToObjectIDs(habitIDs)}}
	if timeFilter != nil {
		if timeFilter.StartTime != nil {
			filter["timestamp"] = bson.M{"$gte": *timeFilter.StartTime}
		}
		if timeFilter.EndTime != nil {
			filter["timestamp"] = bson.M{"$lte": *timeFilter.EndTime}
		}
	}
	return r.FindMany(ctx, filter)
}
