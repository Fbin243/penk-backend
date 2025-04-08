package mongorepo

import (
	"context"
	"time"

	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *TaskSessionRepo) DeleteByTaskID(ctx context.Context, taskID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.DeleteMany(ctx, bson.M{"task_id": mongodb.ToObjectID(taskID)})
	return err
}

func (r *TaskSessionRepo) DeleteByTaskIDs(ctx context.Context, taskIDs []string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.DeleteMany(ctx, bson.M{"task_id": bson.M{"$in": mongodb.ToObjectIDs(taskIDs)}})
	return err
}
