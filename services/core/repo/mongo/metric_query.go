package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/types"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *MetricRepo) Exist(ctx context.Context, characterID, metricID string) error {
	return r.Exists(ctx, bson.M{
		"_id":          mongodb.ToObjectID(metricID),
		"character_id": mongodb.ToObjectID(characterID),
	})
}

func (r *MetricRepo) CountByCharacterID(ctx context.Context, characterID string) (int, error) {
	return r.Count(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
}

func (r *MetricRepo) CountByCategoryID(ctx context.Context, categoryID string) (int, error) {
	return r.Count(ctx, bson.M{"category_id": mongodb.ToObjectID(categoryID)})
}

func (r *MetricRepo) CountUnassigned(ctx context.Context, characterID string) (int, error) {
	return r.Count(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID), "category_id": nil})
}

func (r *MetricRepo) Find(ctx context.Context, p entity.MetricPipeline) ([]entity.Metric, error) {
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

func (r *MetricRepo) CountByFilter(ctx context.Context, filter *entity.MetricFilter) (int, error) {
	return r.AggregateCount(ctx,
		r.addCountStage(
			r.addMatchStage(
				[]bson.M{},
				filter,
			),
		),
	)
}

func (r *MetricRepo) addMatchStage(p []bson.M, filter *entity.MetricFilter) []bson.M {
	if filter == nil {
		return p
	}

	matchStage := bson.M{}
	if filter.CharacterID != nil {
		matchStage["character_id"] = mongodb.ToObjectID(*filter.CharacterID)
	}

	return append(p, bson.M{"$match": matchStage})
}

func (r *MetricRepo) addSortStage(p []bson.M, orderBy *entity.MetricOrderBy) []bson.M {
	return p
}

func (r *MetricRepo) addPaginationStage(p []bson.M, pagination *types.Pagination) []bson.M {
	if pagination == nil {
		return p
	}

	return append(p, mongodb.ToPaginationPineline(pagination)...)
}

func (r *MetricRepo) addCountStage(p []bson.M) []bson.M {
	return append(p, bson.M{"$count": "count"})
}
