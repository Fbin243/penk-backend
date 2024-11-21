package repo

import (
	"context"
	"fmt"
	"log"
	"tenkhours/pkg/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *FishRepo) UpdateFish(fish *Fish, profileID primitive.ObjectID) (*Fish, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("gold f: ", fish.Gold, " Nor f: ", fish.Normal)

	update := bson.M{}
	if fish.Gold != 0 {
		update["gold"] = fish.Gold
	}
	if fish.Normal != 0 {
		update["normal"] = fish.Normal
	}

	if len(update) == 0 {
		return nil, fmt.Errorf("no valid values to update")
	}

	var updatedFish Fish
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := r.Collection.FindOneAndUpdate(ctx, bson.M{"profile_id": profileID}, bson.M{"$set": update}, opts).Decode(&updatedFish)
	log.Println("gold: ", updatedFish.Gold, " Nor: ", updatedFish.Normal)
	if err != nil {
		return nil, fmt.Errorf("failed to update fish: %v", err)
	}

	return &updatedFish, nil
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
