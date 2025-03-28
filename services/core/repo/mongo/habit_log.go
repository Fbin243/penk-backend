package mongorepo

import (
	"context"
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"
	mongomodel "tenkhours/services/core/repo/mongo/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HabitLogRepo struct {
	*mongodb.BaseRepo[entity.HabitLog, mongomodel.HabitLog]
}

func NewHabitLogRepo(db *mongo.Database) *HabitLogRepo {
	_ = db.CreateCollection(context.Background(), mongodb.HabitLogsCollection,
		options.CreateCollection().
			SetTimeSeriesOptions(
				options.TimeSeries().
					SetTimeField("timestamp").
					SetMetaField("habit_id").
					SetGranularity("hours"),
			),
	)

	return &HabitLogRepo{mongodb.NewBaseRepo(
		db.Collection(mongodb.HabitLogsCollection),
		&mongodb.Mapper[entity.HabitLog, mongomodel.HabitLog]{},
		false,
	)}
}

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

func (r *HabitLogRepo) DeleteByHabitIDAndTimestamp(ctx context.Context, habitID string, timestamp time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.DeleteOne(ctx, bson.M{
		"habit_id":  mongodb.ToObjectID(habitID),
		"timestamp": bson.M{"$gte": utils.ResetTimeToBeginningOfDay(timestamp)},
	})

	return err
}
