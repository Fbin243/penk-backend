package mongorepo

import (
	"context"
	"fmt"
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/errors"
	"tenkhours/services/notification/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReminderRepo struct {
	*mongodb.BaseRepo[entity.Reminder, Reminder]
}

func NewReminderRepo(db *mongo.Database) *ReminderRepo {
	return &ReminderRepo{mongodb.NewBaseRepo[entity.Reminder, Reminder](
		db.Collection(mongodb.RemindersCollection),
		true,
	)}
}

func (r *ReminderRepo) CreateReminder(ctx context.Context, reminder *entity.Reminder) (*entity.Reminder, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	newReminder, err := r.InsertOne(ctx, reminder)
	if err != nil {
		return nil, fmt.Errorf("failed to create reminder: %w", err)
	}

	return newReminder, nil
}

func (r *ReminderRepo) GetRemindersByProfileID(ctx context.Context, profileID string) ([]*entity.Reminder, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var reminders []*entity.Reminder
	filter := bson.M{"profile_id": mongodb.ToObjectID(profileID)}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &reminders); err != nil {
		return nil, err
	}

	if len(reminders) == 0 {
		return nil, errors.ErrNotFound
	}

	return reminders, nil
}

func (r *ReminderRepo) GetReminderByID(ctx context.Context, reminderID string) (*entity.Reminder, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	reminder, err := r.FindByID(ctx, reminderID)
	if err != nil {
		return nil, err
	}

	return reminder, nil
}

func (r *ReminderRepo) UpdateReminder(ctx context.Context, reminder *entity.Reminder) (*entity.Reminder, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": mongodb.ToObjectID(reminder.ID)}
	update := bson.M{"$set": reminder}

	_, err := r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update reminder: %w", err)
	}

	var updatedReminder entity.Reminder
	err = r.Collection.FindOne(ctx, filter).Decode(&updatedReminder)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated reminder: %w", err)
	}

	return &updatedReminder, nil
}

func (r *ReminderRepo) DeleteReminder(ctx context.Context, reminderID string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": mongodb.ToObjectID(reminderID)}

	result, err := r.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to delete reminder: %w", err)
	}

	if result.DeletedCount == 0 {
		return false, errors.ErrNotFound
	}

	return true, nil
}

func (r *ReminderRepo) GetUpcomingReminders(ctx context.Context) ([]*entity.Reminder, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var reminders []*entity.Reminder
	filter := bson.M{"remind_time": bson.M{"$gt": time.Now()}}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &reminders); err != nil {
		return nil, err
	}

	if len(reminders) == 0 {
		return nil, errors.ErrNotFound
	}

	return reminders, nil
}
