package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/types"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *TaskSessionRepo) Find(ctx context.Context, p entity.TaskSessionPipeline) ([]entity.TaskSession, error) {
	return r.AggregateQuery(ctx,
		r.addPaginationStage(
			r.addSortStage(
				r.addMatchStage(
					[]bson.M{},
					p.Filter,
				),
				p.OrderBy,
			),
			p.Pagination,
		),
	)
}

func (r *TaskSessionRepo) CountByTaskID(ctx context.Context, taskID string) (int, error) {
	return r.Count(ctx, bson.M{"task_id": mongodb.ToObjectID(taskID)})
}

func (r *TaskSessionRepo) CountByFilter(ctx context.Context, filter *entity.TaskSessionFilter) (int, error) {
	return r.AggregateCount(ctx,
		r.addCountStage(
			r.addMatchStage(
				[]bson.M{},
				filter,
			),
		),
	)
}

func (r *TaskSessionRepo) addMatchStage(p []bson.M, filter *entity.TaskSessionFilter) []bson.M {
	if filter == nil {
		return p
	}

	matchStage := bson.M{}
	if filter.TaskID != nil {
		matchStage["task_id"] = mongodb.ToObjectID(*filter.TaskID)
	} else if filter.TaskIDs != nil {
		matchStage["task_id"] = bson.M{
			"$in": mongodb.ToObjectIDs(filter.TaskIDs),
		}
	}

	timeRange := bson.M{}
	if filter.StartTime != nil {
		timeRange["$gte"] = filter.StartTime
	}
	if filter.EndTime != nil {
		timeRange["$lte"] = filter.EndTime
	}
	if len(timeRange) > 0 {
		matchStage["start_time"] = timeRange
	}

	if filter.IsCompleted != nil {
		if *filter.IsCompleted {
			matchStage["completed_time"] = bson.M{"$ne": nil}
		} else {
			matchStage["completed_time"] = nil
		}
	}

	return append(p, bson.M{"$match": matchStage})
}

func (r *TaskSessionRepo) addCountStage(p []bson.M) []bson.M {
	return append(p, bson.M{"$count": "count"})
}

func (r *TaskSessionRepo) addSortStage(p []bson.M, _ *entity.TaskSessionOrderBy) []bson.M {
	return p
}

func (r *TaskSessionRepo) addPaginationStage(p []bson.M, pagination *types.Pagination) []bson.M {
	if pagination == nil {
		return p
	}

	return append(p, mongodb.ToPaginationPineline(pagination)...)
}
