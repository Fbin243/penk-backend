package mongorepo

import (
	"context"
	"log"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/timetracking/entity"
	mongomodel "tenkhours/services/timetracking/repo/mongo/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TimeTrackingRepo struct {
	*mongodb.BaseRepo[entity.TimeTracking, mongomodel.TimeTracking]
}

func NewTimeTrackingRepo(db *mongo.Database) *TimeTrackingRepo {
	timeTrackingColl := db.Collection(mongodb.TimeTrackingsCollecion)
	_, err := timeTrackingColl.Indexes().CreateMany(context.Background(),
		[]mongo.IndexModel{
			{
				Keys: bson.M{
					"character_id": 1,
				},
			},
			{
				Keys: bson.M{
					"category_id": 1,
				},
			},
			{
				Keys: bson.M{
					"end_time": 1,
				},
			},
		})
	if err != nil {
		log.Printf("failed to create indexes for time tracking collection: %v", err)
		return nil
	}

	return &TimeTrackingRepo{mongodb.NewBaseRepo(
		timeTrackingColl,
		&mongodb.Mapper[entity.TimeTracking, mongomodel.TimeTracking]{},
		false,
	)}
}
