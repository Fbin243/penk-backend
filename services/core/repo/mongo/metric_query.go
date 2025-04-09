package mongorepo

import (
	"context"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *MetricRepo) FindByCharacterID(ctx context.Context, characterID string) ([]entity.Metric, error) {
	return r.FindMany(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
}

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
