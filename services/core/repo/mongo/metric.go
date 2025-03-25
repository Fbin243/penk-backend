package mongorepo

import (
	"context"
	"log"
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"
	mongomodel "tenkhours/services/core/repo/mongo/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MetricRepo struct {
	*mongodb.BaseRepo[entity.Metric, mongomodel.Metric]
}

func NewMetricRepo(db *mongo.Database) *MetricRepo {
	metricColl := db.Collection(mongodb.MetricsCollection)
	_, err := metricColl.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "character_id", Value: 1}}},
		{Keys: bson.D{{Key: "category_id", Value: 1}}},
	})
	if err != nil {
		log.Printf("failed to create indexes for %s collection\n", mongodb.MetricsCollection)
		return nil
	}

	return &MetricRepo{mongodb.NewBaseRepo(
		metricColl,
		&mongodb.Mapper[entity.Metric, mongomodel.Metric]{},
		true,
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

func (r *MetricRepo) UnassignCategory(ctx context.Context, categoryID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.UpdateMany(ctx, bson.M{"category_id": mongodb.ToObjectID(categoryID)}, bson.M{"$unset": bson.M{"category_id": ""}})
	return err
}

func (r *MetricRepo) DeleteByCharacterID(ctx context.Context, characterID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.DeleteMany(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
	return err
}

func (r *MetricRepo) DeleteByCharacterIDs(ctx context.Context, characterIDs []string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.DeleteMany(ctx, bson.M{"character_id": bson.M{"$in": mongodb.ToObjectIDs(characterIDs)}})
	return err
}
