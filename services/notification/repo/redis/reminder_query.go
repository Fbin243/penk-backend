package redisrepo

import (
	"context"
	"fmt"

	rediskeys "tenkhours/pkg/db/redis"
	core_entity "tenkhours/services/core/entity"

	"github.com/go-redis/redis/v8"
)

func (r *ReminderCache) GetReminder(ctx context.Context, id string) (*core_entity.Reminder, error) {
	score, err := r.client.ZScore(ctx, rediskeys.ReminderKey, id).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get reminder from redis: %v", err)
	}

	return ToReminder(&redis.Z{
		Member: id,
		Score:  score,
	}), nil
}

func (r *ReminderCache) GetAllReminders(ctx context.Context) ([]core_entity.Reminder, error) {
	results, err := r.client.ZRangeWithScores(ctx, rediskeys.ReminderKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get reminders from redis: %v", err)
	}

	reminders := make([]core_entity.Reminder, 0, len(results))
	for _, result := range results {
		reminders = append(reminders, *ToReminder(&result))
	}

	return reminders, nil
}

// GetMinScore gets the minimum score value from the sorted set, ignoring score 0
func (r *ReminderCache) GetMinScore(ctx context.Context) (float64, error) {
	results, err := r.client.ZRangeByScoreWithScores(ctx, rediskeys.ReminderKey, &redis.ZRangeBy{
		Min:    "1",    // Ignore score 0
		Max:    "+inf", // Any positive score
		Offset: 0,
		Count:  1, // Get just one to find the minimum score
	}).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get min score from redis: %v", err)
	}
	if len(results) == 0 {
		return 0, nil
	}

	return results[0].Score, nil
}

// GetRemindersByScore gets all reminders with a specific score
func (r *ReminderCache) GetRemindersByScore(ctx context.Context, minScore, maxScore float64) ([]core_entity.Reminder, error) {
	results, err := r.client.ZRangeByScoreWithScores(ctx, rediskeys.ReminderKey, &redis.ZRangeBy{
		Min: fmt.Sprintf("%f", minScore),
		Max: fmt.Sprintf("%f", maxScore),
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get reminders by score from redis: %v", err)
	}

	reminders := make([]core_entity.Reminder, 0, len(results))
	for _, result := range results {
		reminders = append(reminders, *ToReminder(&result))
	}

	return reminders, nil
}
