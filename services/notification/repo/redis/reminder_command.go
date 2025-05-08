package redisrepo

import (
	"context"
	"fmt"

	rediskeys "tenkhours/pkg/db/redis"
	core_entity "tenkhours/services/core/entity"
)

func (r *ReminderCache) SetReminder(ctx context.Context, reminder *core_entity.Reminder) error {
	err := r.client.ZAdd(ctx, rediskeys.ReminderKey, ToRedisZ(reminder)).Err()
	if err != nil {
		return fmt.Errorf("failed to set reminder in redis: %v", err)
	}

	return nil
}

func (r *ReminderCache) SetReminders(ctx context.Context, reminders []core_entity.Reminder) error {
	pipe := r.client.Pipeline()

	for _, reminder := range reminders {
		pipe.ZAdd(ctx, rediskeys.ReminderKey, ToRedisZ(&reminder))
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to cache reminders: %v", err)
	}

	return nil
}

func (r *ReminderCache) DeleteReminder(ctx context.Context, id string) error {
	if err := r.client.ZRem(ctx, rediskeys.ReminderKey, id).Err(); err != nil {
		return fmt.Errorf("failed to delete reminder from redis: %v", err)
	}

	return nil
}

func (r *ReminderCache) ClearReminders(ctx context.Context) error {
	if err := r.client.Del(ctx, rediskeys.ReminderKey).Err(); err != nil {
		return fmt.Errorf("failed to clear reminders from redis: %v", err)
	}
	return nil
}
