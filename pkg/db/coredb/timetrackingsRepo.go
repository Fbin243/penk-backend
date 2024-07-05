package coredb

import (
	"context"
	"time"

	"tenkhours/pkg/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TimeTrackingsRepo struct {
	*mongo.Collection
}

func NewTimeTrackingsRepo(mongodb *mongo.Database) *TimeTrackingsRepo {
	return &TimeTrackingsRepo{mongodb.Collection(db.TimeTrackingsCollection)}
}

func (r *TimeTrackingsRepo) GetTimeTrackingByID(id primitive.ObjectID) (*TimeTracking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	timeTracking := &TimeTracking{}
	err := r.FindOne(ctx, bson.M{"_id": id}).Decode(timeTracking)

	return timeTracking, err
}

func (r *TimeTrackingsRepo) GetTimeTrackingsByCharacterID(characterID primitive.ObjectID) ([]TimeTracking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var timeTrackings []TimeTracking
	cursor, err := r.Find(ctx, bson.M{"character_id": characterID})
	if err != nil {
		return nil, err
	}

	cursor.All(ctx, &timeTrackings)

	return timeTrackings, err
}

func (r *TimeTrackingsRepo) CreateTimeTracking(timeTracking *TimeTracking) (*TimeTracking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.InsertOne(ctx, timeTracking)

	return timeTracking, err
}

func (r *TimeTrackingsRepo) UpdateTimeTracking(timeTracking *TimeTracking) (*TimeTracking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.FindOneAndUpdate(ctx, bson.M{"_id": timeTracking.ID}, bson.M{"$set": timeTracking}).Decode(timeTracking)

	return timeTracking, err
}

func (r *TimeTrackingsRepo) DeleteTimeTracking(id primitive.ObjectID) (*TimeTracking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	deletedTimeTracking := &TimeTracking{}
	err := r.FindOneAndDelete(ctx, bson.M{"_id": id}).Decode(deletedTimeTracking)

	return deletedTimeTracking, err
}
