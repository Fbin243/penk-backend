package mongorepo

import (
	"context"
	"log"

	mongodb "tenkhours/pkg/db/mongo"
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

	return &MetricRepo{mongodb.NewBaseRepo[entity.Metric, mongomodel.Metric](
		metricColl,
		true,
	)}
}
