package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/types"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *TaskSessionRepo) FindByTaskID(ctx context.Context, taskID string, timeFilter *types.TimeFilter) ([]entity.TaskSession, error) {
	filter := bson.M{"task_id": mongodb.ToObjectID(taskID)}
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

func (r *TaskSessionRepo) FindByTaskIDs(ctx context.Context, taskIDs []string, timeFilter *types.TimeFilter) ([]entity.TaskSession, error) {
	filter := bson.M{"task_id": bson.M{"$in": mongodb.ToObjectIDs(taskIDs)}}
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
