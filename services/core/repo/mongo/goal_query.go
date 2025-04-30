package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *GoalRepo) Exist(ctx context.Context, characterID, goalID string) error {
	return r.Exists(ctx, bson.M{
		"_id":          mongodb.ToObjectID(goalID),
		"character_id": mongodb.ToObjectID(characterID),
	})
}

func (r *GoalRepo) CountByCharacterID(ctx context.Context, characterID string) (int, error) {
	return r.Count(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
}

func (r *GoalRepo) Find(ctx context.Context, p entity.GoalPipeline) ([]entity.Goal, error) {
	pipeline := []bson.M{}

	// Add match stage
	if p.Filter != nil {
		pipeline = append(pipeline, bson.M{"$match": r.buildMatchStage(*p.Filter)})
	}

	// Add pagination stage
	pipeline = append(pipeline, mongodb.ToPaginationPineline(p.Pagination)...)

	return r.AggregateQuery(ctx, pipeline)
}

func (r *GoalRepo) CountByFilter(ctx context.Context, filter *entity.GoalFilter) (int, error) {
	pipeline := []bson.M{}

	// Add match stage
	if filter != nil {
		pipeline = append(pipeline, bson.M{"$match": r.buildMatchStage(*filter)})
	}

	// Add count stage
	pipeline = append(pipeline, bson.M{"$count": "count"})

	return r.AggregateCount(ctx, pipeline)
}

func (r *GoalRepo) buildMatchStage(filter entity.GoalFilter) bson.M {
	matchStage := bson.M{}
	if filter.CharacterID != nil {
		matchStage["character_id"] = mongodb.ToObjectID(*filter.CharacterID)
	}

	return matchStage
}
