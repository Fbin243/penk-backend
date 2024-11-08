package repo

import (
	"context"
	"fmt"
	"tenkhours/pkg/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FishRepo struct {
	*mongo.Collection
}

func NewFishRepo(mongodb *mongo.Database) *FishRepo {
	return &FishRepo{mongodb.Collection(db.FishCollection)}
}

func (r *FishRepo) GetFishByProfileID(id primitive.ObjectID, fishType string) (*Fish, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fish := Fish{}
	err := r.FindOne(ctx, bson.M{"profile_id": id, "type": fishType}).Decode(&fish)

	return &fish, err
}

func (r *FishRepo) CreateFish(fish *Fish) (*Fish, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.InsertOne(ctx, fish)

	return fish, err
}

func (r *FishRepo) UpdateFish(fish *Fish, profileID primitive.ObjectID) (*Fish, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// update fish based on profile id and type
	filter := bson.M{"profile_id": profileID, "type": fish.Type}
	update := bson.M{"$set": fish}

	err := r.FindOneAndUpdate(ctx, filter, update, db.FindOneAndUpdateOptions).Decode(fish)
	if err != nil {
		return nil, fmt.Errorf("failed to update fish: %v", err)
	}

	return fish, nil
}
