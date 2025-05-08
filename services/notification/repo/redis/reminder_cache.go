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
	"github.com/samber/lo"
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

	var remindTime *time.Time
	if score != 0 {
		remindTime = lo.ToPtr(time.Unix(int64(score), 0))
	}

	return &core_entity.Reminder{
		BaseEntity: &base.BaseEntity{
			ID: id,
		},
		RemindTime: remindTime,
	}, nil
}

func (r *ReminderCache) SetReminder(ctx context.Context, reminder *core_entity.Reminder) error {
	var score float64
	if reminder.RemindTime != nil {
		score = float64(reminder.RemindTime.Unix())
	}

	if err := r.client.ZAdd(ctx, rediskeys.ReminderKey, &redis.Z{
		Score:  score,
		Member: reminder.ID,
	}).Err(); err != nil {
		return fmt.Errorf("failed to set reminder in redis: %v", err)
	}

	return nil
}

func (r *ReminderCache) SetReminders(ctx context.Context, reminders []core_entity.Reminder) error {
	pipe := r.client.Pipeline()

	for _, reminder := range reminders {
		var score float64
		if reminder.RemindTime != nil {
			score = float64(reminder.RemindTime.Unix())
		}

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
		// If score is 0, it means the reminder has no remind time
		var remindTime *time.Time
		if result.Score != 0 {
			remindTime = lo.ToPtr(time.Unix(int64(result.Score), 0))
		}

		reminder := core_entity.Reminder{
			BaseEntity: &base.BaseEntity{
				ID: result.Member.(string),
			},
			RemindTime: remindTime,
		}
		reminders = append(reminders, reminder)
	}

	return reminders, nil
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

func (r *ReminderCache) GetRemindersWithMinScore(ctx context.Context) ([]core_entity.Reminder, error) {
	// Get all elements with scores > 0
	results, err := r.client.ZRangeByScoreWithScores(ctx, rediskeys.ReminderKey, &redis.ZRangeBy{
		Min:    "1",    // Slightly above 0 to exclude 0 scores
		Max:    "+inf", // Any positive score
		Offset: 0,
		Count:  1, // Get just one to find the minimum score
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get min score from redis: %v", err)
	}
	if len(results) == 0 {
		return nil, nil
	}

	// Get all elements with the same minimum score
	minScore := results[0].Score
	results, err = r.client.ZRangeByScoreWithScores(ctx, rediskeys.ReminderKey, &redis.ZRangeBy{
		Min: fmt.Sprintf("%f", minScore),
		Max: fmt.Sprintf("%f", minScore),
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get reminders with min score from redis: %v", err)
	}

	reminders := make([]core_entity.Reminder, 0, len(results))
	for _, result := range results {
		remindTime := lo.ToPtr(time.Unix(int64(result.Score), 0))
		reminder := core_entity.Reminder{
			BaseEntity: &base.BaseEntity{
				ID: result.Member.(string),
			},
			RemindTime: remindTime,
		}
		reminders = append(reminders, reminder)
	}

	return reminders, nil
}
