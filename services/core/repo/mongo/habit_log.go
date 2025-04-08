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

func (r *HabitLogRepo) UpsertByTimestamp(ctx context.Context, timestamp time.Time, habitLog *entity.HabitLog) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.DeleteMany(ctx, bson.M{
		"habit_id":  mongodb.ToObjectID(habitLog.HabitID),
		"timestamp": bson.M{"$gte": utils.ResetTimeToBeginningOfDay(timestamp)},
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
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.DeleteMany(ctx, bson.M{"habit_id": mongodb.ToObjectID(habitID)})
	if err != nil {
		return err
	}

	return nil
}

func (r *HabitLogRepo) DeleteByHabitIDs(ctx context.Context, habitIDs []string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.DeleteMany(ctx, bson.M{"habit_id": bson.M{"$in": mongodb.ToObjectIDs(habitIDs)}})
	if err != nil {
		return err
	}

	return nil
}
