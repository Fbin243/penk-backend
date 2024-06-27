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

func NewTimeTrackingsRepo() *TimeTrackingsRepo {
	return &TimeTrackingsRepo{db.GetTimeTrackingsCollection()}
}

func (r *TimeTrackingsRepo) GetTimeTrackingByID(id primitive.ObjectID) (TimeTracking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var timeTracking TimeTracking
	err := r.FindOne(ctx, bson.M{"_id": id}).Decode(&timeTracking)

	return timeTracking, err
}

func (r *TimeTrackingsRepo) CreateTimeTracking(timeTracking TimeTracking) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	insertResult, err := r.InsertOne(ctx, timeTracking)

	return insertResult, err
}

func (r *TimeTrackingsRepo) UpdateTimeTracking(timeTracking TimeTracking) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updateResult, err := r.UpdateOne(ctx, bson.M{"_id": timeTracking.ID}, bson.M{"$set": timeTracking})

	return updateResult, err
}
