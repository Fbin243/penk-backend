package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"
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

func buildTaskPipeline(p *entity.TaskPineline) []bson.M {
	pipeline := []bson.M{}

	// Add match stage
	if p.Filter != nil {
		matchStage := bson.M{}
		if p.Filter.CharacterID != nil {
			matchStage["character_id"] = mongodb.ToObjectID(*p.Filter.CharacterID)
		} else if p.Filter.CharacterIDs != nil {
			matchStage["character_id"] = bson.M{
				"$in": mongodb.ToObjectIDs(p.Filter.CharacterIDs),
			}
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

	// Add sorting stage
	if p.OrderBy != nil {
		sortStage := bson.M{}
		if p.OrderBy.Priority != nil {
			sortStage["priority"] = p.OrderBy.Priority.ToInt()
		}

		pipeline = append(pipeline, bson.M{"$sort": sortStage})
	}

	// Add pagination stage
	pipeline = append(pipeline, mongodb.ToPaginationPineline(p.Pagination)...)

	return pipeline
}

func (r *TaskRepo) Find(ctx context.Context, pineline entity.TaskPineline) ([]entity.Task, error) {
	return r.AggregateQuery(ctx, buildTaskPipeline(&pineline))
}
