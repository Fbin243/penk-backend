package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/analytic/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *TimeTrackingRepo) AggregateDailyCapturedRecord(ctx context.Context, filter entity.GetCapturedRecordFilter) ([]entity.CapturedRecord, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "character_id", Value: mongodb.ToObjectID(filter.CharacterID)},
			{Key: "start_time", Value: bson.D{{Key: "$gte", Value: filter.StartTime}}},
			{Key: "end_time", Value: bson.D{{Key: "$lt", Value: filter.EndTime.AddDate(0, 0, 1)}}},
		}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "character", Value: "$character_id"},
				{Key: "category", Value: "$category_id"},
				{Key: "date", Value: bson.D{{Key: "$dateToString", Value: bson.D{{Key: "format", Value: "%Y-%m-%d"}, {Key: "date", Value: "$end_time"}}}}},
			}},
			{Key: "total_duration", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$subtract", Value: bson.A{"$end_time", "$start_time"}}}}}},
			{Key: "timestamp", Value: bson.D{{Key: "$first", Value: "$end_time"}}},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "character_id", Value: "$_id.character"},
			{Key: "category_id", Value: bson.D{{Key: "$ifNull", Value: bson.A{"$_id.category", "unassigned"}}}},
			{Key: "date", Value: "$_id.date"},
			{Key: "timestamp", Value: "$timestamp"},
			{Key: "year", Value: bson.D{{Key: "$year", Value: "$timestamp"}}},
			{Key: "month", Value: bson.D{{Key: "$month", Value: "$timestamp"}}},
			{Key: "day", Value: bson.D{{Key: "$dayOfMonth", Value: "$timestamp"}}},
			{Key: "week", Value: bson.D{{Key: "$isoWeek", Value: "$timestamp"}}},
			{Key: "time", Value: bson.D{{Key: "$divide", Value: bson.A{"$total_duration", 1000}}}},
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

	return results, nil
}
