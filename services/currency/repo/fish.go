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

func (r *FishRepo) UpdateFish(fish *Fish) (*Fish, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.FindOneAndUpdate(ctx, bson.M{"_id": fish.ID}, bson.M{"$set": fish}, db.FindOneAndUpdateOptions).Decode(fish)

	return fish, err
}
