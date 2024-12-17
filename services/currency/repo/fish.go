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

func (r *FishRepo) GetFishByProfileID(profileID primitive.ObjectID) (*Fish, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fish := Fish{}
	err := r.FindOne(ctx, bson.M{"profile_id": profileID}).Decode(&fish)

	return &fish, err
}

func (r *FishRepo) CreateFish(fish *Fish) (*Fish, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.InsertOne(ctx, fish)

	return fish, err
}

func (r *FishRepo) UpdateFishByProfileID(profileID primitive.ObjectID, fish *Fish) (*Fish, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": fish,
	}

	_, err := r.Collection.UpdateOne(ctx, bson.M{"profile_id": profileID}, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update fish: %v", err)
	}

	return fish, nil
}

func (r *FishRepo) DeleteFish(profileID primitive.ObjectID, fishType string) (*Fish, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use FindOneAndDelete to delete the fish document
	filter := bson.M{"profile_id": profileID, "type": fishType}
	result := r.FindOneAndDelete(ctx, filter)

	// Check for errors
	err := result.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("fish not found for deletion")
		}
		return nil, fmt.Errorf("failed to delete fish: %v", err)
	}

	var fish Fish
	err = result.Decode(&fish)
	if err != nil {
		return nil, fmt.Errorf("failed to decode deleted fish: %v", err)
	}

	return &fish, nil
}
