package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *HabitRepo) CountByCharacterID(ctx context.Context, characterID string) (int, error) {
	return r.Count(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
}

func (r *HabitRepo) CountByCategoryID(ctx context.Context, categoryID string) (int, error) {
	return r.Count(ctx, bson.M{"category_id": mongodb.ToObjectID(categoryID)})
}

func (r *HabitRepo) CountUnassigned(ctx context.Context, characterID string) (int, error) {
	return r.Count(ctx, bson.M{
		"character_id": mongodb.ToObjectID(characterID),
		"category_id":  nil,
	})
}

func (r *HabitRepo) Exist(ctx context.Context, characterID, habitID string) error {
	return r.Exists(ctx, bson.M{
		"_id":          mongodb.ToObjectID(habitID),
		"character_id": mongodb.ToObjectID(characterID),
	})
}

func (r *HabitRepo) Find(ctx context.Context, p entity.HabitPipeline) ([]entity.Habit, error) {
	pipeline := []bson.M{}

	// Add match stage
	if p.Filter != nil {
		pipeline = append(pipeline, bson.M{"$match": r.buildMatchStage(*p.Filter)})
	}

	// Add pagination stage
	pipeline = append(pipeline, mongodb.ToPaginationPineline(p.Pagination)...)

	return r.AggregateQuery(ctx, pipeline)
}

func (r *HabitRepo) CountByFilter(ctx context.Context, filter *entity.HabitFilter) (int, error) {
	pipeline := []bson.M{}

	// Add match stage
	if filter != nil {
		pipeline = append(pipeline, bson.M{"$match": r.buildMatchStage(*filter)})
	}

	// Add count stage
	pipeline = append(pipeline, bson.M{"$count": "count"})

	return r.AggregateCount(ctx, pipeline)
}

func (r *HabitRepo) buildMatchStage(filter entity.HabitFilter) bson.M {
	matchStage := bson.M{}
	if filter.CharacterID != nil {
		matchStage["character_id"] = mongodb.ToObjectID(*filter.CharacterID)
	} else if filter.CharacterIDs != nil {
		matchStage["character_id"] = bson.M{
			"$in": mongodb.ToObjectIDs(filter.CharacterIDs),
		}
	}

	return matchStage
}
