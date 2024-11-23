package repo

import (
	"context"
	"time"

	"tenkhours/pkg/db"
	"tenkhours/services/analytics/graph/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CapturedRecordsRepo struct {
	*mongo.Collection
}

func NewCapturedRecordsRepo(mongodb *mongo.Database) *CapturedRecordsRepo {
	return &CapturedRecordsRepo{mongodb.Collection(db.CapturedRecordsCollection)}
}

func (r *CapturedRecordsRepo) CreateCapturedRecord(capturedRecord *model.CapturedRecord) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.InsertOne(ctx, capturedRecord)
	if err != nil {
		return err
	}

	return nil
}

func (r *CapturedRecordsRepo) GetCapturedRecords(pineline mongo.Pipeline) ([]model.CapturedRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.Aggregate(ctx, pineline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var capturedRecords []model.CapturedRecord
	for cursor.Next(ctx) {
		var capturedRecord model.CapturedRecord
		if err := cursor.Decode(&capturedRecord); err != nil {
			return nil, err
		}

		capturedRecords = append(capturedRecords, capturedRecord)
	}

	return capturedRecords, nil
}

func (r *CapturedRecordsRepo) DeleteCapturedRecordsByProfileID(profileID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.DeleteMany(ctx, bson.M{"metadata.profile_id": profileID})
	if err != nil {
		return err
	}

	return nil
}
