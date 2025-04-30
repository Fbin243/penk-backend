package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"

	"github.com/go-redis/redis/v8"
)

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
