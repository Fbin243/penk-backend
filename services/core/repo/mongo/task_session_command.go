package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *TaskSessionRepo) DeleteByTaskID(ctx context.Context, taskID string) error {
	return r.DeleteMany(ctx, bson.M{"task_id": mongodb.ToObjectID(taskID)})
}

func (r *TaskSessionRepo) DeleteByTaskIDs(ctx context.Context, taskIDs []string) error {
	return r.DeleteMany(ctx, bson.M{"task_id": bson.M{"$in": mongodb.ToObjectIDs(taskIDs)}})
}
