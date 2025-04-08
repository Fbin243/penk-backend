package mongorepo

import (
	"context"
	"log"
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"
	mongomodel "tenkhours/services/core/repo/mongo/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type HabitRepo struct {
	*mongodb.BaseRepo[entity.Habit, mongomodel.Habit]
}

func NewHabitRepo(db *mongo.Database) *HabitRepo {
	habitCollection := db.Collection(mongodb.HabitsCollection)
	_, err := habitCollection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "character_id", Value: 1}}},
		{Keys: bson.D{{Key: "category_id", Value: 1}}},
	})
	if err != nil {
		log.Printf("failed to create indexes for %s collection\n", mongodb.HabitsCollection)
		return nil
	}

	return &HabitRepo{mongodb.NewBaseRepo(
		habitCollection,
		&mongodb.Mapper[entity.Habit, mongomodel.Habit]{},
		true,
	)}
}

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

func (r *HabitRepo) DeleteByCharacterID(ctx context.Context, characterID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.DeleteMany(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
	return err
}

func (r *HabitRepo) DeleteByCharacterIDs(ctx context.Context, characterIDs []string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.DeleteMany(ctx, bson.M{"character_id": bson.M{"$in": mongodb.ToObjectIDs(characterIDs)}})
	return err
}
