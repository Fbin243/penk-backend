package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"tenkhours/services/core/entity"

	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"

	"github.com/go-redis/redis/v8"
)

type RedisRepo struct {
	*redis.Client
}

func NewRedisRepo(rdb *redis.Client) *RedisRepo {
	return &RedisRepo{rdb}
}

func (r *RedisRepo) GetAuthSession(ctx context.Context, firebaseUID string) (*rdb.AuthSession, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	sessionJSON, err := r.Get(ctx, rdb.GetAuthSessionKey(firebaseUID)).Result()
	if err == redis.Nil {
		return nil, errors.ErrRedisNotFound
	}
	if err != nil {
		return nil, err
	}

	session := rdb.AuthSession{}
	err = json.Unmarshal([]byte(sessionJSON), &session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *RedisRepo) SetAuthSession(ctx context.Context, profile *entity.Profile, session *rdb.AuthSession) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	authSessionJSON, err := json.Marshal(session)
	if err != nil {
		return err
	}
	err = r.Set(ctx, rdb.GetAuthSessionKey(profile.FirebaseUID), string(authSessionJSON), time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to set auth session in redis: %v", err)
	}

	return nil
}

func (r *RedisRepo) DeleteAuthSession(ctx context.Context, firebaseUID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := r.Del(ctx, rdb.GetAuthSessionKey(firebaseUID)).Err()
	if err != nil {
		return fmt.Errorf("failed to delete auth session in redis: %v", err)
	}

	return nil
}

func (r *RedisRepo) DeleteProfileData(ctx context.Context, profile *entity.Profile) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Delete current timetracking in redis
	err := r.Del(ctx, rdb.GetTimeTrackingKey(profile.ID)).Err()
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
