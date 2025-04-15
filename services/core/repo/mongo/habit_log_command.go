package mongorepo

import (
	"context"
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *HabitLogRepo) UpsertByTimestamp(ctx context.Context, timestamp time.Time, habitLog *entity.HabitLog) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := r.DeleteMany(ctx, bson.M{
		"habit_id":  mongodb.ToObjectID(habitLog.HabitID),
		"timestamp": bson.M{"$gte": utils.StartOfDay(timestamp)},
	})
	if err != nil {
		return err
	}

	_, err = r.InsertOne(ctx, habitLog)
	if err != nil {
		return err
	}

	return err
}

func (r *HabitLogRepo) DeleteByHabitID(ctx context.Context, habitID string) error {
	return r.DeleteMany(ctx, bson.M{"habit_id": mongodb.ToObjectID(habitID)})
}

func (r *HabitLogRepo) DeleteByHabitIDs(ctx context.Context, habitIDs []string) error {
	return r.DeleteMany(ctx, bson.M{"habit_id": bson.M{"$in": mongodb.ToObjectIDs(habitIDs)}})
}
