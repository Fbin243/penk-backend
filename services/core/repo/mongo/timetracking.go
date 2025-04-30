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
					"reference_id": 1,
				},
			},
			{
				Keys: bson.M{
					"timestamp": 1,
				},
			},
		})
	if err != nil {
		log.Printf("failed to create indexes for time tracking collection: %v", err)
		return nil
	}

	return &TimeTrackingRepo{mongodb.NewBaseRepo[entity.TimeTracking, mongomodel.TimeTracking](
		timeTrackingColl,
		false,
	)}
}
