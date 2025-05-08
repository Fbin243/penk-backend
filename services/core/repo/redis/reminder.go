package redisrepo

import (
	"context"

	rdb "tenkhours/pkg/db/redis"
	"tenkhours/services/core/entity"

	errs "tenkhours/pkg/errors"

	"github.com/go-redis/redis/v8"
)

type ReminderCache struct {
	*redis.Client
}

func NewReminderCache(client *redis.Client) *ReminderCache {
	return &ReminderCache{client}
}

func (r *ReminderCache) Exist(ctx context.Context, reminder *entity.Reminder) error {
	exists, err := r.Exists(ctx, rdb.ReminderKey).Result()
	if err == redis.Nil || exists == 0 {
		return errs.ErrRedisNotFound
	}
	if err != nil {
		return err
	}

	return nil
}

// UpsertReminder adds or updates a reminder in the sorted set
// Score is the Unix timestamp of the reminder time
func (r *ReminderCache) UpsertReminder(ctx context.Context, reminder *entity.Reminder) error {
	return r.ZAdd(ctx, rdb.ReminderKey, &redis.Z{
		Score:  float64(reminder.RemindTime.Unix()),
		Member: reminder.ID,
	}).Err()
}

// DeleteReminder removes a reminder from the sorted set
func (r *ReminderCache) DeleteReminder(ctx context.Context, id string) error {
	return r.ZRem(ctx, rdb.ReminderKey, id).Err()
}
