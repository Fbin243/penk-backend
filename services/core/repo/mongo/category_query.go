package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/types"
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
	return r.AggregateQuery(ctx,
		r.addPaginationStage(
			r.addSortStage(
				r.addMatchStage(
					[]bson.M{},
					pineline.Filter,
				),
				pineline.OrderBy,
			),
			pineline.Pagination,
		),
	)
}

func (r *CategoryRepo) CountByFilter(ctx context.Context, filter *entity.CategoryFilter) (int, error) {
	return r.AggregateCount(ctx,
		r.addCountStage(
			r.addMatchStage(
				[]bson.M{},
				filter,
			),
		),
	)
}

func (r *CategoryRepo) addMatchStage(p []bson.M, filter *entity.CategoryFilter) []bson.M {
	if filter == nil {
		return p
	}

	matchStage := bson.M{}
	if filter.CharacterID != nil {
		matchStage["character_id"] = mongodb.ToObjectID(*filter.CharacterID)
	}

	return append(p, bson.M{"$match": matchStage})
}

func (r *CategoryRepo) addCountStage(p []bson.M) []bson.M {
	return append(p, bson.M{"$count": "count"})
}

func (r *CategoryRepo) addSortStage(p []bson.M, _ *entity.CategoryOrderBy) []bson.M {
	return p
}

func (r *CategoryRepo) addPaginationStage(p []bson.M, pagination *types.Pagination) []bson.M {
	if pagination == nil {
		return p
	}

	return append(p, mongodb.ToPaginationPineline(pagination)...)
}
