package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *TaskSessionRepo) Find(ctx context.Context, p entity.TaskSessionPipeline) ([]entity.TaskSession, error) {
	pipeline := []bson.M{}

	// Add match stage
	if p.Filter != nil {
		matchStage := bson.M{}
		if p.Filter.TaskID != nil {
			matchStage["task_id"] = mongodb.ToObjectID(*p.Filter.TaskID)
		} else if p.Filter.TaskIDs != nil {
			matchStage["task_id"] = bson.M{
				"$in": mongodb.ToObjectIDs(p.Filter.TaskIDs),
			}
		}

		timeRange := bson.M{}
		if p.Filter.StartTime != nil {
			timeRange["$gte"] = p.Filter.StartTime
		}
		if p.Filter.EndTime != nil {
			timeRange["$lte"] = p.Filter.EndTime
		}
		if len(timeRange) > 0 {
			matchStage["start_time"] = timeRange
		}

		if p.Filter.IsCompleted != nil {
			if *p.Filter.IsCompleted {
				matchStage["completed_time"] = bson.M{"$ne": nil}
			} else {
				matchStage["completed_time"] = nil
			}
		}

		pipeline = append(pipeline, bson.M{"$match": matchStage})
	}

	return r.AggregateQuery(ctx, pipeline)
}

func (r *TaskSessionRepo) CountByTaskID(ctx context.Context, taskID string) (int, error) {
	return r.Count(ctx, bson.M{"task_id": mongodb.ToObjectID(taskID)})
}

func (r *TaskSessionRepo) CountByCharacterID(ctx context.Context, characterID string) (int, error) {
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "tasks",
				"localField":   "task_id",
				"foreignField": "_id",
				"as":           "task",
			},
		},
		{
			"$match": bson.M{
				"task.character_id": mongodb.ToObjectID(characterID),
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
