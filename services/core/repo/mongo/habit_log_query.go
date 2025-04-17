package mongorepo

import (
	"context"
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func buildHabitLogPipeline(p *entity.HabitLogPineline) []bson.M {
	pipeline := []bson.M{}

	// Add match stage
	if p.Filter != nil {
		matchStage := bson.M{}
		if p.Filter.HabitID != nil {
			matchStage["habit_id"] = mongodb.ToObjectID(*p.Filter.HabitID)
		} else if p.Filter.HabitIDs != nil {
			matchStage["habit_id"] = bson.M{
				"$in": mongodb.ToObjectIDs(p.Filter.HabitIDs),
			}
		}

		matchStage["timestamp"] = bson.M{}
		if p.Filter.StartTime != nil {
			matchStage["timestamp"].(bson.M)["$gte"] = p.Filter.StartTime.Format(time.DateOnly)
		}
		if p.Filter.EndTime != nil {
			matchStage["timestamp"].(bson.M)["$lte"] = p.Filter.EndTime.Format(time.DateOnly)
		}

		if p.Filter.Reset != nil {
			// Get all habit logs of a habit has this reset
			pipeline = append(pipeline, bson.M{"$lookup": bson.M{
				"from":         mongodb.HabitsCollection,
				"localField":   "habit_id",
				"foreignField": "_id",
				"as":           "habit",
			}})
			pipeline = append(pipeline, bson.M{"$unwind": "$habit"})
			matchStage["habit.reset"] = *p.Filter.Reset
		}
		pipeline = append(pipeline, bson.M{"$match": matchStage})
	}

	// Add sorting stage
	if p.OrderBy != nil {
		sortStage := bson.M{}
		if p.OrderBy.Timestamp != nil {
			sortStage["timestamp"] = p.OrderBy.Timestamp.ToInt()
		}
		pipeline = append(pipeline, bson.M{"$sort": sortStage})
	}

	// Add limit stage
	if p.Limit != nil {
		pipeline = append(pipeline, bson.M{"$limit": *p.Limit})
	}

	// Add skip stage
	if p.Offset != nil {
		pipeline = append(pipeline, bson.M{"$skip": *p.Offset})
	}

	return pipeline
}

func (r *HabitLogRepo) FindByPineline(ctx context.Context, pineline entity.HabitLogPineline) ([]entity.HabitLog, error) {
	return r.AggregateQuery(ctx, buildHabitLogPipeline(&pineline))
}
