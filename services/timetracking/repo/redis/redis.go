package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/timetracking/entity"

	mongodb "tenkhours/pkg/db/mongo"
	rdb "tenkhours/pkg/db/redis"

	"github.com/go-redis/redis/v8"
)

type RedisRepo struct {
	*redis.Client
}

func NewRedisRepo(rdb *redis.Client) *RedisRepo {
	return &RedisRepo{rdb}
}

func (r *RedisRepo) CreateTimeTracking(ctx context.Context, profileID string, timeTracking *entity.TimeTracking) error {
	// Save the time tracking to Redis
	timeTrackingJSON, err := json.Marshal(timeTracking)
	if err != nil {
		return fmt.Errorf("failed to serialize time tracking: %v", err)
	}

	err = r.Set(ctx, rdb.GetTimeTrackingKey(profileID), timeTrackingJSON, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to save time tracking to redis: %v", err)
	}

	return nil
}

func (r *RedisRepo) GetCurrentTimeTracking(ctx context.Context, profileID string) (*entity.TimeTracking, error) {
	// Get the current time tracking from Redis
	currentTimetrackingJSON, err := r.Get(ctx, rdb.GetTimeTrackingKey(profileID)).Result()
	if err == redis.Nil {
		return nil, errors.ErrRedisNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to get time tracking from redis: %v", err)
	}

	var currentTimetracking entity.TimeTracking
	err = json.Unmarshal([]byte(currentTimetrackingJSON), &currentTimetracking)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize time tracking: %v", err)
	}

	return &currentTimetracking, nil
}

func (r *RedisRepo) DeleteCurrentTimeTracking(ctx context.Context, profileID string) error {
	// Delete the current time tracking from Redis
	err := r.Del(ctx, rdb.GetTimeTrackingKey(profileID)).Err()
	if err != nil {
		return fmt.Errorf("failed to delete time tracking from redis: %v", err)
	}

	return nil
}

func (r *RedisRepo) GetCurrentCapturedRecord(ctx context.Context, profileID, characterID string) (*entity.CapturedRecord, error) {
	// Get the timetracking from the current captured record in Redis
	capturedRecordJSON, err := r.HGet(ctx, rdb.GetCapturedRecordKey(profileID), characterID).Result()
	if err == redis.Nil {
		return nil, errors.NewGQLError(errors.ErrCodeRedisNotFound, nil)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get captured record from redis: %v", err)
	}

	var capturedRecord entity.CapturedRecord
	err = json.Unmarshal([]byte(capturedRecordJSON), &capturedRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize captured record: %v", err)
	}

	return &capturedRecord, nil
}

func (r *RedisRepo) UpsertTimeTrackingInCapturedRecord(ctx context.Context, profileID string, timeTracking *entity.TimeTracking, duration int32) error {
	// Check if the captured record already exists in Redis
	capturedRecord := entity.CapturedRecord{}
	capturedRecordJSON, err := r.HGet(ctx, rdb.GetCapturedRecordKey(profileID), timeTracking.CharacterID).Result()
	if err == redis.Nil {
		// Make a new captured record if it doesn't exist
		capturedRecord = entity.CapturedRecord{
			ID:               mongodb.GenObjectID(),
			Timestamp:        utils.ResetTimeToBeginningOfDay(timeTracking.EndTime),
			TotalFocusedTime: 0,
			Metadata: entity.CapturedRecordMetadata{
				CharacterID: timeTracking.CharacterID,
				ProfileID:   profileID,
			},
			TimeTrackings: []entity.CapturedRecordTimeTracking{},
			Categories:    []entity.CapturedRecordCategory{},
		}
	} else if err != nil {
		return fmt.Errorf("failed to get captured record from redis: %v", err)
	} else {
		err = json.Unmarshal([]byte(capturedRecordJSON), &capturedRecord)
		if err != nil {
			return fmt.Errorf("failed to deserialize captured record: %v", err)
		}
	}

	// Add the time to the captured record for character
	capturedRecord.TotalFocusedTime += duration

	// Add the time tracking to the captured record
	capturedRecord.TimeTrackings = append(capturedRecord.TimeTrackings, entity.CapturedRecordTimeTracking{
		CategoryID: timeTracking.CategoryID,
		Time:       duration,
		StartTime:  timeTracking.StartTime,
		EndTime:    timeTracking.EndTime,
	})

	if timeTracking.CategoryID != nil {
		// Check if this custom metric already exists in the captured record
		found := false
		for j, capturedCustomMetric := range capturedRecord.Categories {
			// If it exists, add the time to it
			if capturedCustomMetric.ID == *timeTracking.CategoryID {
				capturedRecord.Categories[j].Time += int32(duration)
				found = true
			}
		}

		// If it doesn't exist, create a new one
		if !found {
			capturedRecord.Categories = append(capturedRecord.Categories, entity.CapturedRecordCategory{
				ID:   *timeTracking.CategoryID,
				Time: int32(duration),
			})
		}
	}

	// Upsert the captured record to Redis
	capturedRecordBytes, err := json.Marshal(capturedRecord)
	if err != nil {
		return fmt.Errorf("failed to serialize captured record: %v", err)
	}

	err = r.HSet(ctx, rdb.GetCapturedRecordKey(profileID), timeTracking.CharacterID, capturedRecordBytes).Err()
	if err != nil {
		return fmt.Errorf("failed to save captured record to redis: %v", err)
	}

	return nil
}
