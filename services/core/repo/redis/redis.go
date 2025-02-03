package redisrepo

import (
	"context"
	"fmt"
	"time"

	"tenkhours/services/core/entity"

	rdb "tenkhours/pkg/db/redis"

	"github.com/go-redis/redis/v8"
)

type RedisRepo struct {
	*redis.Client
}

func NewRedisRepo(rdb *redis.Client) *RedisRepo {
	return &RedisRepo{rdb}
}

func (r *RedisRepo) DeleteProfileData(ctx context.Context, profile *entity.Profile) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Delete current captured record in redis
	err := r.Del(ctx, rdb.GetCapturedRecordKey(profile.ID)).Err()
	if err != nil {
		return fmt.Errorf("failed to delete current captured record in redis: %v", err)
	}

	// Delete current timetracking in redis
	err = r.Del(ctx, rdb.GetTimeTrackingKey(profile.ID)).Err()
	if err != nil {
		return fmt.Errorf("failed to delete current timetracking in redis: %v", err)
	}

	// Delete profile in redis
	err = r.Del(ctx, rdb.GetAuthSessionKey(profile.FirebaseUID)).Err()
	if err != nil {
		return fmt.Errorf("failed to delete profile in redis: %v", err)
	}

	return err
}
