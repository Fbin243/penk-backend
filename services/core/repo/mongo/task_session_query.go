package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/types"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *TaskSessionRepo) FindByTaskID(ctx context.Context, taskID string, timeFilter *types.TimeFilter) ([]entity.TaskSession, error) {
	return r.FindMany(ctx, bson.M{
		"task_id":    mongodb.ToObjectID(taskID),
		"start_time": mongodb.MakeTimeFilter(timeFilter),
	})
}

func (r *TaskSessionRepo) FindByTaskIDs(ctx context.Context, taskIDs []string, timeFilter *types.TimeFilter) ([]entity.TaskSession, error) {
	return r.FindMany(ctx, bson.M{
		"task_id":    bson.M{"$in": mongodb.ToObjectIDs(taskIDs)},
		"start_time": mongodb.MakeTimeFilter(timeFilter),
	})
}
