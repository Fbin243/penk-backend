package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *TaskSessionRepo) FindByTaskID(ctx context.Context, taskID string, filter *entity.TaskSessionFilter) ([]entity.TaskSession, error) {
	mongoFilter := bson.M{
		"task_id": mongodb.ToObjectID(taskID),
	}

	if filter != nil {
		mongoFilter["start_time"] = bson.M{}

		if filter.StartTime != nil {
			mongoFilter["start_time"].(bson.M)["$gte"] = filter.StartTime
		}
		if filter.EndTime != nil {
			mongoFilter["start_time"].(bson.M)["$lte"] = filter.EndTime
		}
	}

	return r.FindMany(ctx, mongoFilter)
}

func (r *TaskSessionRepo) FindByTaskIDs(ctx context.Context, taskIDs []string, filter *entity.TaskSessionFilter) ([]entity.TaskSession, error) {
	mongoFilter := bson.M{
		"task_id": bson.M{
			"$in": mongodb.ToObjectIDs(taskIDs),
		},
	}

	if filter != nil {
		mongoFilter["start_time"] = bson.M{}

		if filter.StartTime != nil {
			mongoFilter["start_time"].(bson.M)["$gte"] = filter.StartTime
		}
		if filter.EndTime != nil {
			mongoFilter["start_time"].(bson.M)["$lte"] = filter.EndTime
		}
	}

	return r.FindMany(ctx, mongoFilter)
}
