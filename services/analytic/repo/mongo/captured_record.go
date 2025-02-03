package mongorepo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"tenkhours/services/analytic/entity"

	mongodb "tenkhours/pkg/db/mongo"
	rdb "tenkhours/pkg/db/redis"

	"github.com/go-redis/redis/v8"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CapturedRecordRepo struct {
	*mongo.Collection
	*mongodb.Mapper[entity.CapturedRecord, CapturedRecord]
	*redis.Client
}

func NewCapturedRecordRepo(db *mongo.Database, rdb *redis.Client) *CapturedRecordRepo {
	return &CapturedRecordRepo{
		db.Collection(mongodb.CapturedRecordsCollection),
		&mongodb.Mapper[entity.CapturedRecord, CapturedRecord]{},
		rdb,
	}
}

func (r *CapturedRecordRepo) CreateCapturedRecord(ctx context.Context, capturedRecord *entity.CapturedRecord) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.InsertOne(ctx, r.ToMongoEntity(capturedRecord))
	if err != nil {
		return err
	}

	return nil
}

func (r *CapturedRecordRepo) GetCapturedRecords(ctx context.Context, filter entity.GetCapturedRecordFilter) ([]entity.CapturedRecord, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	matchStage := bson.D{}
	if filter.CharacterID != nil {
		matchStage = append(matchStage, bson.E{Key: "metadata.character_id", Value: filter.CharacterID})
	} else {
		matchStage = append(matchStage, bson.E{Key: "metadata.profile_id", Value: filter.ProfileID})
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

	// Get the current captured record from redis
	currentCapturedRecords := []entity.CapturedRecord{}
	currentCapturedRecordsJSON, err := r.HGetAll(context.Background(), rdb.GetCapturedRecordKey(filter.ProfileID)).Result()
	if err == redis.Nil {
		log.Printf("no current captured record found in redis for profile: %s", filter.ProfileID)
	} else if err != nil {
		log.Printf("failed to get current captured record from redis: %v", err)
	}

	for characterID, capturedRecordJSON := range currentCapturedRecordsJSON {
		if filter.CharacterID != nil && lo.FromPtr(filter.CharacterID) != characterID {
			continue
		}

		var capturedRecord entity.CapturedRecord
		err = json.Unmarshal([]byte(capturedRecordJSON), &capturedRecord)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal current captured record: %v", err)
		}

		if capturedRecord.Timestamp.After(filter.StartTime) && capturedRecord.Timestamp.Before(filter.EndTime) {
			currentCapturedRecords = append(currentCapturedRecords, capturedRecord)
		}
	}

	// Append the current captured record to the list if it is within the time range
	capturedRecords = append(capturedRecords, currentCapturedRecords...)

	return capturedRecords, nil
}

func (r *CapturedRecordRepo) DeleteCapturedRecords(ctx context.Context, profileID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.DeleteMany(ctx, bson.M{"metadata.profile_id": profileID})
	if err != nil {
		return err
	}

	return nil
}
