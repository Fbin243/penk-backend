package mongorepo

import (
	"context"
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *MetricRepo) FindByCharacterID(ctx context.Context, characterID string) ([]entity.Metric, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.Find(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	metrics := []entity.Metric{}
	err = cursor.All(ctx, &metrics)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func (r *MetricRepo) Exist(ctx context.Context, characterID, metricID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	count, err := r.CountDocuments(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID), "_id": mongodb.ToObjectID(metricID)})
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.ErrMongoNotFound
	}

	return nil
}

func (r *MetricRepo) CountByCharacterID(ctx context.Context, characterID string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	count, err := r.CountDocuments(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
	return int(count), err
}

func (r *MetricRepo) CountByCategoryID(ctx context.Context, categoryID string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	count, err := r.CountDocuments(ctx, bson.M{"category_id": mongodb.ToObjectID(categoryID)})
	return int(count), err
}

func (r *MetricRepo) CountUnassigned(ctx context.Context, characterID string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	count, err := r.CountDocuments(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID), "category_id": nil})
	return int(count), err
}
