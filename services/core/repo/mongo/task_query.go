package mongorepo

import (
	"context"
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *TaskRepo) CountByCharacterID(ctx context.Context, characterID string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	count, err := r.CountDocuments(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
	return int(count), err
}

func (r *TaskRepo) Exist(ctx context.Context, characterID, taskID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	count, err := r.CountDocuments(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID), "_id": mongodb.ToObjectID(taskID)})
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.ErrMongoNotFound
	}

	return nil
}

func (r *TaskRepo) FindByCharacterID(ctx context.Context, characterID string) ([]entity.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.Find(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	tasks := []entity.Task{}
	err = cursor.All(ctx, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepo) FindByCharacterIDs(ctx context.Context, characterIDs []string) ([]entity.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.Find(ctx, bson.M{"character_id": bson.M{"$in": mongodb.ToObjectIDs(characterIDs)}})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	tasks := []entity.Task{}
	err = cursor.All(ctx, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
