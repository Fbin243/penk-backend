package mongorepo

import (
	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/mongo"
)

type TimeTrackingRepo struct {
	*mongo.Collection
}

func NewTimeTrackingRepo(db *mongo.Database) *TimeTrackingRepo {
	return &TimeTrackingRepo{
		Collection: db.Collection(mongodb.TimeTrackingsCollecion),
	}
}
