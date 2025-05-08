package redisrepo

import (
	"context"
	"fmt"
	"time"

	"tenkhours/pkg/db/base"
	rediskeys "tenkhours/pkg/db/redis"
	core_entity "tenkhours/services/core/entity"
	"tenkhours/services/notification/business"

	"github.com/go-redis/redis/v8"
)

type ReminderCache struct {
	client *redis.Client
}

func NewReminderCache(client *redis.Client) business.IReminderCache {
	return &ReminderCache{client: client}
}

func (r *ReminderCache) GetReminder(ctx context.Context, id string) (*core_entity.Reminder, error) {
	score, err := r.client.ZScore(ctx, rediskeys.ReminderKey, id).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get reminder from redis: %v", err)
	}

	return &core_entity.Reminder{
		BaseEntity: &base.BaseEntity{
			ID: id,
		},
		RemindTime: time.Unix(int64(score), 0),
	}, nil
}

func (r *ReminderCache) SetReminder(ctx context.Context, reminder *core_entity.Reminder) error {
	score := float64(reminder.RemindTime.Unix())
	if err := r.client.ZAdd(ctx, rediskeys.ReminderKey, &redis.Z{
		Score:  score,
		Member: reminder.ID,
	}).Err(); err != nil {
		return fmt.Errorf("failed to set reminder in redis: %v", err)
	}

	return nil
}

func (r *ReminderCache) DeleteReminder(ctx context.Context, id string) error {
	if err := r.client.ZRem(ctx, rediskeys.ReminderKey, id).Err(); err != nil {
		return fmt.Errorf("failed to delete reminder from redis: %v", err)
	}

	return nil
}

func (r *ReminderCache) SetReminders(ctx context.Context, reminders []core_entity.Reminder) error {
	pipe := r.client.Pipeline()

	for _, reminder := range reminders {
		score := float64(reminder.RemindTime.Unix())
		pipe.ZAdd(ctx, rediskeys.ReminderKey, &redis.Z{
			Score:  score,
			Member: reminder.ID,
		})
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to cache reminders: %v", err)
	}

	return nil
}

func (r *ReminderCache) GetAllReminders(ctx context.Context) ([]core_entity.Reminder, error) {
	results, err := r.client.ZRangeWithScores(ctx, rediskeys.ReminderKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get reminders from redis: %v", err)
	}

	reminders := make([]core_entity.Reminder, 0, len(results))
	for _, result := range results {
		reminder := core_entity.Reminder{
			BaseEntity: &base.BaseEntity{
				ID: result.Member.(string),
			},
			RemindTime: time.Unix(int64(result.Score), 0),
		}
		reminders = append(reminders, reminder)
	}

	return reminders, nil
}

func (r *ReminderCache) ClearReminders(ctx context.Context) error {
	if err := r.client.Del(ctx, rediskeys.ReminderKey).Err(); err != nil {
		return fmt.Errorf("failed to clear reminders from redis: %v", err)
	}
	return nil
}
