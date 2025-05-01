package mongorepo

import (
	"context"
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/types"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *HabitLogRepo) Find(ctx context.Context, p entity.HabitLogPipeline) ([]entity.HabitLog, error) {
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

func (r *HabitLogRepo) CountByHabitID(ctx context.Context, habitID string) (int, error) {
	return r.Count(ctx, bson.M{"habit_id": mongodb.ToObjectID(habitID)})
}

func (r *HabitLogRepo) CountByFilter(ctx context.Context, filter *entity.HabitLogFilter) (int, error) {
	return r.AggregateCount(ctx,
		r.addCountStage(
			r.addMatchStage(
				[]bson.M{},
				filter,
			),
		),
	)
}

func (r *HabitLogRepo) addMatchStage(p []bson.M, filter *entity.HabitLogFilter) []bson.M {
	if filter == nil {
		return p
	}

	matchStage := bson.M{}
	if filter.HabitID != nil {
		matchStage["habit_id"] = mongodb.ToObjectID(*filter.HabitID)
	} else if filter.HabitIDs != nil {
		matchStage["habit_id"] = bson.M{
			"$in": mongodb.ToObjectIDs(filter.HabitIDs),
		}
	}

	timeRange := bson.M{}
	if filter.StartTime != nil {
		timeRange["$gte"] = filter.StartTime.Format(time.DateOnly)
	}
	if filter.EndTime != nil {
		timeRange["$lte"] = filter.EndTime.Format(time.DateOnly)
	}
	if len(timeRange) > 0 {
		matchStage["timestamp"] = timeRange
	}

	if filter.ResetDuration != nil {
		// Get all habit logs of a habit has this reset
		p = append(p, bson.M{"$lookup": bson.M{
			"from":         mongodb.HabitsCollection,
			"localField":   "habit_id",
			"foreignField": "_id",
			"as":           "habit",
		}})
		p = append(p, bson.M{"$unwind": "$habit"})
		matchStage["habit.reset_duration"] = *filter.ResetDuration
	}

	return append(p, bson.M{"$match": matchStage})
}

func (r *HabitLogRepo) addSortStage(p []bson.M, orderBy *entity.HabitLogOrderBy) []bson.M {
	if orderBy == nil {
		return p
	}

	sortStage := bson.M{}
	if orderBy.Timestamp != nil {
		sortStage["timestamp"] = orderBy.Timestamp.ToInt()
	}

	return append(p, bson.M{"$sort": sortStage})
}

func (r *HabitLogRepo) addPaginationStage(p []bson.M, pagination *types.Pagination) []bson.M {
	if pagination == nil {
		return p
	}

	return append(p, mongodb.ToPaginationPineline(pagination)...)
}

func (r *HabitLogRepo) addCountStage(p []bson.M) []bson.M {
	return append(p, bson.M{"$count": "count"})
}
