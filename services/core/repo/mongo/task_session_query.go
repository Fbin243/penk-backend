package mongorepo

import (
	"context"
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *TaskSessionRepo) FindByTaskID(ctx context.Context, taskID string) ([]entity.TaskSession, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.Find(ctx, bson.M{"task_id": mongodb.ToObjectID(taskID)})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	taskSessions := []entity.TaskSession{}
	err = cursor.All(ctx, &taskSessions)
	if err != nil {
		return nil, err
	}

	return taskSessions, nil
}

func (r *TaskSessionRepo) FindByTaskIDs(ctx context.Context, taskIDs []string) ([]entity.TaskSession, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.Find(ctx, bson.M{"task_id": bson.M{"$in": mongodb.ToObjectIDs(taskIDs)}})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	taskSessions := []entity.TaskSession{}
	err = cursor.All(ctx, &taskSessions)
	if err != nil {
		return nil, err
	}

	return taskSessions, nil
}
