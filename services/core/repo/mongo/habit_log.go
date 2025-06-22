package mongorepo

import (
	"context"
	"log"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"
	mongomodel "tenkhours/services/core/repo/mongo/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type HabitLogRepo struct {
	*mongodb.BaseRepo[entity.HabitLog, mongomodel.HabitLog]
}

func NewHabitLogRepo(db *mongo.Database) *HabitLogRepo {
	habitLogColl := db.Collection(mongodb.HabitLogsCollection)
	_, err := habitLogColl.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.M{
				"habit_id": 1,
			},
		},
		{
			Keys: bson.M{
				"timestamp": 1,
			},
		},
	})
	if err != nil {
		log.Printf("error creating index for %s: %v", mongodb.HabitLogsCollection, err)
	}
	// TODO: Uncomment this when you want to use time series collection
	// _ := db.CreateCollection(context.Background(), mongodb.HabitLogsCollection)
	// options.CreateCollection().
	// 	SetTimeSeriesOptions(
	// 		options.TimeSeries().
	// 			SetTimeField("timestamp").
	// 			SetMetaField("habit_id").
	// 			SetGranularity("hours"),
	// 	),

	return &HabitLogRepo{mongodb.NewBaseRepo[entity.HabitLog, mongomodel.HabitLog](
		habitLogColl,
		false,
	)}
}
