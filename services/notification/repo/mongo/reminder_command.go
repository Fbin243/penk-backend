package mongorepo

import (
	"context"
	"fmt"
	"log"
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	core_entity "tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
