package business

import (
	"context"
	"fmt"
	"log"
	"time"

	core_entity "tenkhours/services/core/entity"

	"github.com/samber/lo"
	"github.com/teambition/rrule-go"
)

// SyncTodayReminders updates reminders in MongoDB with new remind times and syncs them to Redis
func (biz *NotificationBusiness) SyncTodayReminders(ctx context.Context) error {
	reminders, err := biz.ReminderCache.GetAllReminders(ctx)
	if err != nil {
		return fmt.Errorf("failed to get reminders from redis: %v", err)
	}

	if err := biz.ReminderRepo.BulkUpdateRemindTimes(ctx, reminders); err != nil {
		return fmt.Errorf("failed to update reminders in MongoDB: %v", err)
	}

	todayReminders, err := biz.ReminderRepo.GetTodayReminders(ctx)
	if err != nil {
		return fmt.Errorf("failed to get today's reminders from MongoDB: %v", err)
	}

	if err := biz.ReminderCache.ClearReminders(ctx); err != nil {
		return fmt.Errorf("failed to clear reminders from cache: %v", err)
	}

	if err := biz.ReminderCache.SetReminders(ctx, todayReminders); err != nil {
		return fmt.Errorf("failed to cache reminders: %v", err)
	}

	log.Printf("Successfully synced %d reminders to cache", len(todayReminders))
	return nil
}

// UpdateOutdatedReminders updates reminders that have passed their remind time
func (biz *NotificationBusiness) UpdateOutdatedReminders(ctx context.Context) error {
	// 1. Get all reminders from Redis and update MongoDB
	redisReminders, err := biz.ReminderCache.GetAllReminders(ctx)
	if err != nil {
		return fmt.Errorf("failed to get reminders from redis: %v", err)
	}

	if err := biz.ReminderRepo.BulkUpdateRemindTimes(ctx, redisReminders); err != nil {
		return fmt.Errorf("failed to update reminders in MongoDB: %v", err)
	}

	// 2. Get outdated reminders from MongoDB and calculate next occurrence
	now := time.Now()
	outdatedReminders, err := biz.ReminderRepo.GetOutdatedReminders(ctx, now)
	if err != nil {
		return fmt.Errorf("failed to get outdated reminders from MongoDB: %v", err)
	}

	// Calculate next occurrence for each outdated reminder
	updatedReminders := make([]core_entity.Reminder, 0, len(outdatedReminders))
	for _, reminder := range outdatedReminders {
		updatedReminders = append(updatedReminders, *updateRemindTime(&reminder))
	}

	// 3. Bulk update MongoDB with new remind times
	if err := biz.ReminderRepo.BulkUpdateRemindTimes(ctx, updatedReminders); err != nil {
		return fmt.Errorf("failed to update outdated reminders in MongoDB: %v", err)
	}

	// 4. Clear Redis cache, query all reminders for today and cache them again
	if err := biz.ReminderCache.ClearReminders(ctx); err != nil {
		return fmt.Errorf("failed to clear reminders from cache: %v", err)
	}

	todayReminders, err := biz.ReminderRepo.GetTodayReminders(ctx)
	if err != nil {
		return fmt.Errorf("failed to get today's reminders from MongoDB: %v", err)
	}

	if err := biz.ReminderCache.SetReminders(ctx, todayReminders); err != nil {
		return fmt.Errorf("failed to cache updated reminders: %v", err)
	}

	log.Printf("Successfully updated %d outdated reminders", len(updatedReminders))
	return nil
}

// ProcessRemindersWithMinScore processes all reminders with the minimum score
func (biz *NotificationBusiness) ProcessRemindersWithMinScore(ctx context.Context) (float64, error) {
	// 1. Get all reminders with min score
	reminders, err := biz.ReminderCache.GetRemindersWithMinScore(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get reminders with min score: %v", err)
	}

	// TODO: Query data from mongodb and compose message to Kafka
	// 2. Print all reminders to console
	fmt.Println("Processing reminders with minimum score:")
	for _, reminder := range reminders {
		fmt.Printf("Reminder ID: %s, RemindTime: %v, RRule: %s\n", reminder.ID, reminder.RemindTime, reminder.RRule)
	}

	// 3. Recalculate scores for all reminders
	updatedReminders := make([]core_entity.Reminder, 0, len(reminders))
	for _, reminder := range reminders {
		updatedReminders = append(updatedReminders, *updateRemindTime(&reminder))
	}

	// 4. Update reminders in Redis
	if err := biz.ReminderCache.SetReminders(ctx, updatedReminders); err != nil {
		return 0, fmt.Errorf("failed to update reminders in redis: %v", err)
	}

	log.Printf("Successfully processed %d reminders with minimum score", len(reminders))
	return biz.ReminderCache.GetMinScore(ctx)
}

func updateRemindTime(reminder *core_entity.Reminder) *core_entity.Reminder {
	rule, err := rrule.StrToRRule(reminder.RRule)
	if err != nil {
		log.Printf("Warning: failed to parse RRule for reminder %s: %v", reminder.ID, err)
		return reminder
	}

	nextOccurrence := rule.After(time.Now(), false)
	if nextOccurrence.IsZero() {
		reminder.RemindTime = nil
	} else {
		reminder.RemindTime = lo.ToPtr(time.Date(nextOccurrence.Year(), nextOccurrence.Month(), nextOccurrence.Day(), reminder.RemindTime.Hour(), reminder.RemindTime.Minute(), 0, 0, time.UTC))
	}

	return reminder
}
