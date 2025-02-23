package redis

import (
	"context"
	"encoding/json"

	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/services/analytic/entity"

	"github.com/go-redis/redis/v8"
	"github.com/samber/lo"
)

type RedisRepo struct {
	*redis.Client
}

func NewRedisRepo(client *redis.Client) *RedisRepo {
	return &RedisRepo{
		Client: client,
	}
}

func (r *RedisRepo) GetCapturedRecords(ctx context.Context, filter entity.GetCapturedRecordFilter) ([]entity.CapturedRecord, error) {
	// Get the current captured record from redis
	currentCapturedRecords := []entity.CapturedRecord{}
	currentCapturedRecordsJSON, err := r.HGetAll(context.Background(), rdb.GetCapturedRecordKey(filter.ProfileID)).Result()
	if err == redis.Nil {
		return nil, errors.ErrRedisNotFound
	} else if err != nil {
		return nil, err
	}

	for characterID, capturedRecordJSON := range currentCapturedRecordsJSON {
		if filter.CharacterID != nil && lo.FromPtr(filter.CharacterID) != characterID {
			continue
		}

		var capturedRecord entity.CapturedRecord
		err = json.Unmarshal([]byte(capturedRecordJSON), &capturedRecord)
		if err != nil {
			return nil, err
		}

		if capturedRecord.Timestamp.After(filter.StartTime) && capturedRecord.Timestamp.Before(filter.EndTime) {
			currentCapturedRecords = append(currentCapturedRecords, capturedRecord)
		}
	}

	return currentCapturedRecords, nil
}
