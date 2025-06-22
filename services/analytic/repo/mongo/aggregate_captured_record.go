package mongorepo

import (
	"context"
	"log"

	"tenkhours/services/analytic/entity"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *TimeTrackingRepo) AggregateDailyCapturedRecord(ctx context.Context, filter entity.StatAnalyticFilter) ([]entity.CapturedRecord, error) {
	timeRange := bson.D{}
	if filter.StartTime != nil {
		timeRange = append(timeRange, bson.E{Key: "$gte", Value: filter.StartTime})
	}
	if filter.EndTime != nil {
		timeRange = append(timeRange, bson.E{Key: "$lte", Value: filter.EndTime})
	}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "character_id", Value: mongodb.ToObjectID(filter.CharacterID)},
			{Key: "timestamp", Value: timeRange},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "character_id", Value: "$character_id"},
			{Key: "category_id", Value: bson.D{{Key: "$ifNull", Value: bson.A{"$category_id", "unassigned"}}}},
			{Key: "date", Value: bson.D{{Key: "$dateToString", Value: bson.D{{Key: "format", Value: "%Y-%m-%d"}, {Key: "date", Value: "$timestamp"}}}}},
			{Key: "timestamp", Value: "$timestamp"},
			{Key: "year", Value: bson.D{{Key: "$year", Value: "$timestamp"}}},
			{Key: "month", Value: bson.D{{Key: "$month", Value: "$timestamp"}}},
			{Key: "day", Value: bson.D{{Key: "$dayOfMonth", Value: "$timestamp"}}},
			{Key: "week", Value: bson.D{{Key: "$isoWeek", Value: "$timestamp"}}},
			{Key: "time", Value: "$duration"},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "date", Value: 1}, {Key: "character_id", Value: 1}}}},
	}

	cursor, err := r.Collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []entity.CapturedRecord
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	for _, result := range results {
		log.Printf("Captured Record: %+v\n", utils.PrettyJSON(result))
	}

	return results, nil
}
