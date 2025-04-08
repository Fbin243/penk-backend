package mongorepo

import (
	"context"
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *HabitRepo) CountByCharacterID(ctx context.Context, characterID string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	count, err := r.CountDocuments(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
	return int(count), err
}

func (r *HabitRepo) Exist(ctx context.Context, characterID, habitID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	count, err := r.CountDocuments(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID), "_id": mongodb.ToObjectID(habitID)})
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.ErrMongoNotFound
	}

	return nil
}

func (r *HabitRepo) FindByCharacterID(ctx context.Context, characterID string) ([]entity.Habit, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.Find(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	habits := []entity.Habit{}
	err = cursor.All(ctx, &habits)
	if err != nil {
		return nil, err
	}

	return habits, nil
}

func (r *HabitRepo) FindByCharacterIDs(ctx context.Context, characterIDs []string) ([]entity.Habit, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.Find(ctx, bson.M{"character_id": bson.M{"$in": mongodb.ToObjectIDs(characterIDs)}})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	habits := []entity.Habit{}
	err = cursor.All(ctx, &habits)
	if err != nil {
		return nil, err
	}

	return habits, nil
}
