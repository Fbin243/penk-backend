package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/types"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *TaskRepo) CountByCharacterID(ctx context.Context, characterID string) (int, error) {
	return r.Count(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
}

func (r *TaskRepo) CountByCategoryID(ctx context.Context, categoryID string) (int, error) {
	return r.Count(ctx, bson.M{"category_id": mongodb.ToObjectID(categoryID)})
}

func (r *TaskRepo) CountUnassigned(ctx context.Context, characterID string) (int, error) {
	return r.Count(ctx, bson.M{
		"character_id": mongodb.ToObjectID(characterID),
		"category_id":  nil,
	})
}

func (r *TaskRepo) Exist(ctx context.Context, characterID, taskID string) error {
	return r.Exists(ctx, bson.M{
		"_id":          mongodb.ToObjectID(taskID),
		"character_id": mongodb.ToObjectID(characterID),
	})
}

func (r *TaskRepo) Find(ctx context.Context, p entity.TaskPipeline) ([]entity.Task, error) {
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

func (r *TaskRepo) CountByFilter(ctx context.Context, filter *entity.TaskFilter) (int, error) {
	return r.AggregateCount(ctx,
		r.addCountStage(
			r.addMatchStage(
				[]bson.M{},
				filter,
			),
		),
	)
}

func (r *TaskRepo) addMatchStage(p []bson.M, filter *entity.TaskFilter) []bson.M {
	if filter == nil {
		return p
	}

	matchStage := bson.M{}
	if filter.CharacterID != nil {
		matchStage["character_id"] = mongodb.ToObjectID(*filter.CharacterID)
	} else if filter.CharacterIDs != nil {
		matchStage["character_id"] = bson.M{
			"$in": mongodb.ToObjectIDs(filter.CharacterIDs),
		}
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

func (r *TaskRepo) addSortStage(p []bson.M, orderBy *entity.TaskOrderBy) []bson.M {
	if orderBy == nil {
		return p
	}

	sortStage := bson.M{}
	if orderBy.Priority != nil {
		sortStage["priority"] = orderBy.Priority.ToInt()
	}

	return append(p, bson.M{"$sort": sortStage})
}

func (r *TaskRepo) addPaginationStage(p []bson.M, pagination *types.Pagination) []bson.M {
	if pagination == nil {
		return p
	}

	return append(p, mongodb.ToPaginationPineline(pagination)...)
}

func (r *TaskRepo) addCountStage(p []bson.M) []bson.M {
	return append(p, bson.M{"$count": "count"})
}
