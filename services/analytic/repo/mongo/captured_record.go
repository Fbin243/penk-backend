package mongorepo

import (
	"context"
	"time"

	"tenkhours/services/analytic/entity"

	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CapturedRecordRepo struct {
	*mongo.Collection
	*mongodb.Mapper[entity.CapturedRecord, CapturedRecord]
}

func NewCapturedRecordRepo(db *mongo.Database) *CapturedRecordRepo {
	return &CapturedRecordRepo{
		db.Collection(mongodb.CapturedRecordsCollection),
		&mongodb.Mapper[entity.CapturedRecord, CapturedRecord]{},
	}
}

func (r *CapturedRecordRepo) CreateCapturedRecord(ctx context.Context, capturedRecord *entity.CapturedRecord) (*entity.CapturedRecord, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.InsertOne(ctx, r.ToMongoEntity(capturedRecord))
	if err != nil {
		return nil, err
	}

	return capturedRecord, nil
}

func (r *CapturedRecordRepo) GetCapturedRecords(ctx context.Context, filter entity.GetCapturedRecordFilter) ([]entity.CapturedRecord, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	matchStage := bson.D{}
	if filter.CharacterID != nil {
		matchStage = append(matchStage, bson.E{Key: "metadata.character_id", Value: mongodb.ToObjectID(*filter.CharacterID)})
	} else {
		matchStage = append(matchStage, bson.E{Key: "metadata.profile_id", Value: mongodb.ToObjectID(filter.ProfileID)})
	}

	matchStage = append(matchStage, bson.E{Key: "timestamp", Value: bson.D{{Key: "$gte", Value: filter.StartTime}}})
	matchStage = append(matchStage, bson.E{Key: "timestamp", Value: bson.D{{Key: "$lte", Value: filter.EndTime}}})

	pipeline := mongo.Pipeline{}
	pipeline = append(pipeline, bson.D{{Key: "$match", Value: matchStage}})
	pipeline = append(pipeline, bson.D{{Key: "$sort", Value: bson.D{{Key: "timestamp", Value: 1}}}})

	cursor, err := r.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var capturedRecords []entity.CapturedRecord
	for cursor.Next(ctx) {
		var capturedRecord entity.CapturedRecord
		if err := cursor.Decode(&capturedRecord); err != nil {
			return nil, err
		}

		capturedRecords = append(capturedRecords, capturedRecord)
	}

	return capturedRecords, nil
}

func (r *CapturedRecordRepo) DeleteCapturedRecords(ctx context.Context, profileID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.DeleteMany(ctx, bson.M{"metadata.profile_id": mongodb.ToObjectID(profileID)})
	if err != nil {
		return err
	}

	return nil
}
