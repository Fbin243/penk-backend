package mongorepo

import (
	"context"
	"fmt"
	"log"
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	core_entity "tenkhours/services/core/entity"
	core_mongomodel "tenkhours/services/core/repo/mongo/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReminderRepo struct {
	*mongodb.BaseRepo[core_entity.Reminder, core_mongomodel.Reminder]
}

func NewReminderRepo(db *mongo.Database) *ReminderRepo {
	return &ReminderRepo{mongodb.NewBaseRepo[core_entity.Reminder, core_mongomodel.Reminder](
		db.Collection(mongodb.RemindersCollection),
		true,
	)}
}

func (r *ReminderRepo) GetTodayReminders(ctx context.Context) ([]core_entity.Reminder, error) {
	// Get today's date at midnight UTC
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	return r.FindMany(ctx, bson.M{
		"remind_time": bson.M{
			"$gte": today,
			"$lt":  today.Add(24 * time.Hour),
		},
	})
}

// BulkUpdateRemindTimes updates multiple reminders' remind times in a single bulk operation
func (r *ReminderRepo) BulkUpdateRemindTimes(ctx context.Context, reminders []core_entity.Reminder) error {
	if len(reminders) == 0 {
		return nil
	}

	// Create bulk write operations
	operations := make([]mongo.WriteModel, 0, len(reminders))
	for _, reminder := range reminders {
		update := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": mongodb.ToObjectID(reminder.ID)}).
			SetUpdate(bson.M{
				"$set": bson.M{
					"remind_time": reminder.RemindTime,
					"updated_at":  time.Now().UTC(),
				},
			})
		operations = append(operations, update)
	}

	// Execute bulk write
	opts := options.BulkWrite().SetOrdered(false)
	result, err := r.Collection.BulkWrite(ctx, operations, opts)
	if err != nil {
		return fmt.Errorf("failed to bulk update reminders: %v", err)
	}

	log.Printf("Bulk updated %d reminders (matched: %d, modified: %d)",
		len(reminders), result.MatchedCount, result.ModifiedCount)
	return nil
}

// GetOutdatedReminders gets all reminders that have passed their remind time
func (r *ReminderRepo) GetOutdatedReminders(ctx context.Context, now time.Time) ([]core_entity.Reminder, error) {
	return r.FindMany(ctx, bson.M{
		"remind_time": bson.M{
			"$lt": now,
		},
	})
}
