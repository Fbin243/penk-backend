package repo

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

	timeTracking := TimeTracking{}
	err := r.FindOne(ctx, bson.M{"_id": id}).Decode(&timeTracking)
	if err != nil {
		return nil, err
	}

	return &timeTracking, nil
}

func (r *TimeTrackingsRepo) GetTimeTrackingsByCharacterID(characterID primitive.ObjectID) ([]TimeTracking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var timeTrackings []TimeTracking
	cursor, err := r.Find(ctx, bson.M{"character_id": characterID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &timeTrackings)
	if err != nil {
		return nil, err
	}

	return timeTrackings, nil
}

func (r *TimeTrackingsRepo) GetCurrentTimeTrackingByCharacterID(characterID primitive.ObjectID) (*TimeTracking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"character_id": characterID,
		"end_time":     bson.M{"$exists": false},
	}

	var timeTracking TimeTracking
	err := r.FindOne(ctx, filter).Decode(&timeTracking)

	return &timeTracking, err
}

func (r *TimeTrackingsRepo) CreateTimeTracking(timeTracking *TimeTracking) (*TimeTracking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.InsertOne(ctx, timeTracking)
	if err != nil {
		return nil, err
	}

	return timeTracking, nil
}

func (r *TimeTrackingsRepo) UpdateTimeTracking(timeTracking *TimeTracking) (*TimeTracking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": timeTracking.ID}
	update := bson.M{"$set": bson.M{
		"character_id":      timeTracking.CharacterID,
		"custom_metric_id":  timeTracking.CustomMetricID,
		"start_time":        timeTracking.StartTime,
		"end_time":          timeTracking.EndTime,
		"min_duration_time": timeTracking.MinDurationTime,
		"max_duration_time": timeTracking.MaxDurationTime,
	}}

	err := r.FindOneAndUpdate(ctx, filter, update, db.FindOneAndUpdateOptions).Decode(timeTracking)
	if err != nil {
		return nil, err
	}

	return timeTracking, nil
}

func (r *TimeTrackingsRepo) DeleteTimeTracking(id primitive.ObjectID) (*TimeTracking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	deletedTimeTracking := &TimeTracking{}
	err := r.FindOneAndDelete(ctx, bson.M{"_id": id}).Decode(deletedTimeTracking)
	if err != nil {
		return nil, err
	}

	return deletedTimeTracking, nil
}
