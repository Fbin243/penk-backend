package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"
	mongomodel "tenkhours/services/core/repo/mongo/model"

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
