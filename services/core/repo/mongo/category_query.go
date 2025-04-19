package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *CategoryRepo) CountByCharacterID(ctx context.Context, characterID string) (int, error) {
	return r.Count(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
}

func (r *CategoryRepo) Exist(ctx context.Context, characterID, categoryID string) error {
	return r.Exists(ctx, bson.M{
		"_id":          mongodb.ToObjectID(categoryID),
		"character_id": mongodb.ToObjectID(characterID),
	})
}

func (r *CategoryRepo) Find(ctx context.Context, pineline entity.CategoryPipeline) ([]entity.Category, error) {
	pipeline := []bson.M{}

	// Add match stage
	if pineline.Filter != nil {
		matchStage := bson.M{}
		if pineline.Filter.CharacterID != nil {
			matchStage["character_id"] = mongodb.ToObjectID(*pineline.Filter.CharacterID)
		}
		pipeline = append(pipeline, bson.M{"$match": matchStage})
	}

	// Add pagination stage
	pipeline = append(pipeline, mongodb.ToPaginationPineline(pineline.Pagination)...)

	return r.AggregateQuery(ctx, pipeline)
}
