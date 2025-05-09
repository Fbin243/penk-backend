package mongorepo

import (
	"context"
	"fmt"
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	core_entity "tenkhours/services/core/entity"
	"tenkhours/services/notification/entity"

	"go.mongodb.org/mongo-driver/bson"
)

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

// GetOutdatedReminders gets all reminders that have passed their remind time
func (r *ReminderRepo) GetOutdatedReminders(ctx context.Context, now time.Time) ([]core_entity.Reminder, error) {
	return r.FindMany(ctx, bson.M{
		"remind_time": bson.M{
			"$lt": now,
		},
	})
}

func (r *ReminderRepo) GetRemindersAndMetadataByIDs(ctx context.Context, ids []string) ([]entity.ReminderWithMetadata, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id": bson.M{"$in": mongodb.ToObjectIDs(ids)},
			},
		},
		{
			"$lookup": bson.M{
				"from":         mongodb.HabitsCollection,
				"localField":   "reference_id",
				"foreignField": "_id",
				"as":           "habit",
			},
		},
		{
			"$lookup": bson.M{
				"from":         mongodb.TasksCollection,
				"localField":   "reference_id",
				"foreignField": "_id",
				"as":           "task",
			},
		},
		{
			"$addFields": bson.M{
				"habit": bson.M{
					"$arrayElemAt": []any{"$habit", 0},
				},
				"task": bson.M{
					"$arrayElemAt": []any{"$task", 0},
				},
			},
		},
		{
			"$project": bson.M{
				"reminder": "$$ROOT",
				"habit":    1,
				"task":     1,
			},
		},
	}

	cursor, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate reminders: %v", err)
	}
	defer cursor.Close(ctx)

	var results []entity.ReminderWithMetadata
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode reminder results: %v", err)
	}

	return results, nil
}
