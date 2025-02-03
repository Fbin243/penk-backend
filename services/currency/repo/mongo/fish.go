package mongorepo

import (
	"context"
	"time"

	"tenkhours/services/currency/entity"

	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FishRepo struct {
	*mongodb.BaseRepo[entity.Fish, Fish]
}

func NewFishRepo(db *mongo.Database) *FishRepo {
	return &FishRepo{mongodb.NewBaseRepo(
		db.Collection(mongodb.FishCollection),
		&mongodb.Mapper[entity.Fish, Fish]{},
	)}
}

func (r *FishRepo) GetFishByProfileID(ctx context.Context, profileID string) (*entity.Fish, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var fish entity.Fish
	err := r.FindOne(ctx, bson.M{"profile_id": mongodb.ToObjectID(profileID)}).Decode(&fish)
	if err != nil {
		return nil, err
	}

	return &fish, nil
}

func (r *FishRepo) UpdateFishByProfileID(ctx context.Context, profileID string, fish *entity.Fish) (*entity.Fish, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fish.SetUpdatedAtByNow()

	update := bson.M{
		"$set": r.ToMongoEntity(fish),
	}

	_, err := r.UpdateOne(ctx, bson.M{"profile_id": mongodb.ToObjectID(profileID)}, update)

	return fish, err
}

func (r *FishRepo) DeleteFishByProfileID(ctx context.Context, profileID string) (*entity.Fish, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var fish entity.Fish
	err := r.FindOneAndDelete(ctx, bson.M{"profile_id": mongodb.ToObjectID(profileID)}).Decode(&fish)

	return &fish, err
}
