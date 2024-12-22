package repo

import (
	"context"
	"tenkhours/pkg/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FishRepo struct {
	*db.BaseRepo[Fish]
}

func NewFishRepo(mongodb *mongo.Database) *FishRepo {
	return &FishRepo{db.NewBaseRepo[Fish](mongodb.Collection(db.FishCollection))}
}

func (r *FishRepo) GetFishByProfileID(profileID primitive.ObjectID) (*Fish, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var fish Fish
	err := r.FindOne(ctx, bson.M{"profile_id": profileID}).Decode(&fish)
	if err != nil {
		return nil, err
	}

	return &fish, nil
}

func (r *FishRepo) UpdateFishByProfileID(profileID primitive.ObjectID, fish *Fish) (*Fish, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fish.SetUpdatedAtByNow()

	update := bson.M{
		"$set": fish,
	}

	_, err := r.UpdateOne(ctx, bson.M{"profile_id": profileID}, update)

	return fish, err
}

func (r *FishRepo) DeleteFishByProfileID(profileID primitive.ObjectID) (*Fish, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var fish Fish
	err := r.FindOneAndDelete(ctx, bson.M{"profile_id": profileID}).Decode(&fish)

	return &fish, err
}
