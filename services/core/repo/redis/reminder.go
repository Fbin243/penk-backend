package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"tenkhours/services/core/entity"

	"github.com/go-redis/redis/v8"
)

type ReminderCache struct {
	*redis.Client
}

func NewReminderCache(client *redis.Client) *ReminderCache {
	return &ReminderCache{client}
}

func (r *ReminderCache) GetReminder(ctx context.Context, id string) (*entity.Reminder, error) {
	key := fmt.Sprintf("reminder:%s", id)
	data, err := r.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get reminder from redis: %v", err)
	}

	var reminder entity.Reminder
	if err := json.Unmarshal([]byte(data), &reminder); err != nil {
		return nil, fmt.Errorf("failed to unmarshal reminder: %v", err)
	}

	return &reminder, nil
}

func (r *ReminderCache) SetReminder(ctx context.Context, reminder *entity.Reminder) error {
	key := fmt.Sprintf("reminder:%s", reminder.ID)
	data, err := json.Marshal(reminder)
	if err != nil {
		return fmt.Errorf("failed to marshal reminder: %v", err)
	}

	if err := r.Set(ctx, key, data, 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("failed to set reminder in redis: %v", err)
	}

	return nil
}

func (r *ReminderCache) DeleteReminder(ctx context.Context, id string) error {
	key := fmt.Sprintf("reminder:%s", id)
	if err := r.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete reminder from redis: %v", err)
	}

	return nil
}
