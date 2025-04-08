package mongorepo

import (
	"context"
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *HabitLogRepo) FindByHabitID(ctx context.Context, habitID string) ([]entity.HabitLog, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.Collection.Find(ctx, bson.M{"habit_id": mongodb.ToObjectID(habitID)})
	if err != nil {
		return nil, err
	}

	var habitLogs []entity.HabitLog
	if err := cursor.All(ctx, &habitLogs); err != nil {
		return nil, err
	}
	return habitLogs, nil
}
