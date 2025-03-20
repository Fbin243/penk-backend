package mongorepo

import (
	"context"
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"
	mongorepo "tenkhours/services/core/repo/mongo/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MetricRepo struct {
	*mongodb.BaseRepo[entity.Metric, mongorepo.Metric]
}

func NewMetricRepo(db *mongo.Database) *MetricRepo {
	return &MetricRepo{mongodb.NewBaseRepo(
		db.Collection(mongodb.MetricsCollection),
		&mongodb.Mapper[entity.Metric, mongorepo.Metric]{},
	)}
}

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

func (r *MetricRepo) ValidateMetric(ctx context.Context, characterID, metricID string) error {
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
