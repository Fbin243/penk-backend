package mongorepo

import (
	"context"
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *HabitLogRepo) Find(ctx context.Context, p entity.HabitLogPipeline) ([]entity.HabitLog, error) {
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

		timeRange := bson.M{}
		if p.Filter.StartTime != nil {
			timeRange["$gte"] = p.Filter.StartTime.Format(time.DateOnly)
		}
		if p.Filter.EndTime != nil {
			timeRange["$lte"] = p.Filter.EndTime.Format(time.DateOnly)
		}
		if len(timeRange) > 0 {
			matchStage["timestamp"] = timeRange
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

	// Add pagination stage
	pipeline = append(pipeline, mongodb.ToPaginationPineline(p.Pagination)...)

	return r.AggregateQuery(ctx, pipeline)
}

func (r *HabitLogRepo) CountByHabitID(ctx context.Context, habitID string) (int, error) {
	return r.Count(ctx, bson.M{"habit_id": mongodb.ToObjectID(habitID)})
}

func (r *HabitLogRepo) CountByCharacterID(ctx context.Context, characterID string) (int, error) {
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         mongodb.HabitsCollection,
				"localField":   "habit_id",
				"foreignField": "_id",
				"as":           "habit",
			},
		},
		{
			"$match": bson.M{
				"habit.character_id": mongodb.ToObjectID(characterID),
			},
		},
		{
			"$count": "count",
		},
	}

	cursor, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)

	var result []struct {
		Count int `bson:"count"`
	}

	if err := cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}

	return result[0].Count, nil
}
