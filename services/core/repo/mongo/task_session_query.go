package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func buildTaskSessionPipeline(p *entity.TaskSessionPineline) []bson.M {
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

	return pipeline
}

func (r *TaskSessionRepo) Find(ctx context.Context, pineline entity.TaskSessionPineline) ([]entity.TaskSession, error) {
	return r.AggregateQuery(ctx, buildTaskSessionPipeline(&pineline))
}
