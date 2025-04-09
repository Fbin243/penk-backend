package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *TaskSessionRepo) FindByTaskID(ctx context.Context, taskID string) ([]entity.TaskSession, error) {
	return r.FindMany(ctx, bson.M{"task_id": mongodb.ToObjectID(taskID)})
}

func (r *TaskSessionRepo) FindByTaskIDs(ctx context.Context, taskIDs []string) ([]entity.TaskSession, error) {
	return r.FindMany(ctx, bson.M{"task_id": bson.M{"$in": mongodb.ToObjectIDs(taskIDs)}})
}
