package business

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"tenkhours/pkg/db"
	"tenkhours/services/analytics/graph/model"
	"tenkhours/services/analytics/repo"

	"github.com/go-redis/redis/v8"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FilterMethod int

const (
	FilterMethodServerSide FilterMethod = iota
	FilterMethodClientSide
)

type FilterType int

const (
	FilterTypeUser FilterType = iota
	FilterTypeCharacter
)

type CapturedRecordsFilter struct {
	FilterMethod         FilterMethod
	FilterType           FilterType
	CharacterID          primitive.ObjectID
	ProfileID            primitive.ObjectID
	StartTime            time.Time
	EndTime              time.Time
	RedisClient          *redis.Client
	CapturedRecordsRepo  *repo.CapturedRecordsRepo
	CapturedRecordLocals []model.CapturedRecord
}

func (c *CapturedRecordsFilter) Filter() ([]model.CapturedRecord, error) {
	switch c.FilterMethod {
	case FilterMethodServerSide:
		matchStage := bson.D{}
		if c.FilterType == FilterTypeCharacter {
			matchStage = append(matchStage, bson.E{Key: "metadata.character_id", Value: c.CharacterID})
		} else {
			matchStage = append(matchStage, bson.E{Key: "metadata.profile_id", Value: c.ProfileID})
		}

		matchStage = append(matchStage, bson.E{Key: "timestamp", Value: bson.D{{Key: "$gte", Value: c.StartTime}}})
		matchStage = append(matchStage, bson.E{Key: "timestamp", Value: bson.D{{Key: "$lte", Value: c.EndTime}}})

		pipeline := mongo.Pipeline{}
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: matchStage}})
		pipeline = append(pipeline, bson.D{{Key: "$sort", Value: bson.D{{Key: "timestamp", Value: 1}}}})

		// Get captured records from the database
		capturedRecords, err := c.CapturedRecordsRepo.GetCapturedRecords(pipeline)
		if err != nil {
			return nil, err
		}

		// Get the current captured record from redis
		currentCapturedRecords := []model.CapturedRecord{}
		currentCapturedRecordsJSON, err := c.RedisClient.HGetAll(context.Background(), db.CapturedRecordKey+c.ProfileID.Hex()).Result()
		if err == redis.Nil {
			log.Printf("no current captured record found in redis for profile: %s", c.ProfileID.Hex())
		} else if err != nil {
			log.Printf("failed to get current captured record from redis: %v", err)
		}

		for characterID, capturedRecordJSON := range currentCapturedRecordsJSON {
			if !c.CharacterID.IsZero() && characterID != c.CharacterID.Hex() {
				continue
			}

			var capturedRecord model.CapturedRecord
			err = json.Unmarshal([]byte(capturedRecordJSON), &capturedRecord)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal current captured record: %v", err)
			}

			if capturedRecord.Timestamp.After(c.StartTime) && capturedRecord.Timestamp.Before(c.EndTime) {
				currentCapturedRecords = append(currentCapturedRecords, capturedRecord)
			}
		}

		// Append the current captured record to the list if it is within the time range
		capturedRecords = append(capturedRecords, currentCapturedRecords...)

		return capturedRecords, nil

	case FilterMethodClientSide:
		capturedRecordLocals := c.CapturedRecordLocals
		if c.FilterType == FilterTypeCharacter {
			capturedRecordLocals = lo.Filter(capturedRecordLocals, func(record model.CapturedRecord, _ int) bool {
				return record.Metadata.CharacterID == c.CharacterID
			})
		}

		capturedRecordLocals = lo.Filter(capturedRecordLocals, func(record model.CapturedRecord, _ int) bool {
			return record.Timestamp.After(c.StartTime.Add(-24*time.Hour)) && record.Timestamp.Before(c.EndTime.Add(24*time.Hour))
		})

		sort.Slice(capturedRecordLocals, func(i, j int) bool {
			return capturedRecordLocals[i].Timestamp.Before(capturedRecordLocals[j].Timestamp)
		})

		return capturedRecordLocals, nil
	}

	return nil, nil
}
