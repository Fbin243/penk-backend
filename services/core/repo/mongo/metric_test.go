package mongorepo_test

import (
	"context"
	"testing"

	"tenkhours/pkg/db/base"
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func NewMetric() *entity.Metric {
	return &entity.Metric{
		BaseEntity: &base.BaseEntity{
			ID:        mongodb.GenObjectID(),
			CreatedAt: utils.Now(),
			UpdatedAt: utils.Now(),
		},
		CharacterID: mongodb.GenObjectID(),
		CategoryID:  lo.ToPtr(mongodb.GenObjectID()),
		Name:        "Metric name",
		Value:       1.0,
		Unit:        "Metric unit",
	}
}

func assertMetric(t *testing.T, expected, actual entity.Metric) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.CreatedAt, actual.CreatedAt)
	assert.Equal(t, expected.UpdatedAt, actual.UpdatedAt)
	assert.Equal(t, expected.CharacterID, actual.CharacterID)
	assert.Equal(t, expected.CategoryID, actual.CategoryID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Value, actual.Value)
	assert.Equal(t, expected.Unit, actual.Unit)
}

func TestMetricRepo(t *testing.T) {
	characterID := mongodb.GenObjectID()
	for range 3 {
		metric := NewMetric()
		metric.CharacterID = characterID
		createdMetric, err := metricRepo.InsertOne(context.Background(), metric)
		assert.Nil(t, err)
		assertMetric(t, *metric, *createdMetric)
	}

	metrics, err := metricRepo.Find(context.Background(), entity.MetricPipeline{
		Filter: &entity.MetricFilter{
			CharacterID: &characterID,
		},
	})
	assert.Nil(t, err)
	assert.Len(t, metrics, 3)

	count, err := metricRepo.CountByCharacterID(context.Background(), characterID)
	assert.Nil(t, err)
	assert.Equal(t, int64(3), count)

	err = metricRepo.DeleteByCharacterID(context.Background(), characterID)
	assert.Nil(t, err)

	count, err = metricRepo.CountByCharacterID(context.Background(), characterID)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), count)
}

func TestExist(t *testing.T) {
	metric := NewMetric()
	err := metricRepo.Exist(context.Background(), metric.CharacterID, metric.CharacterID)
	assert.NotNil(t, err)
	assert.Equal(t, errors.ErrMongoNotFound, err)

	_, err = metricRepo.InsertOne(context.Background(), metric)
	assert.Nil(t, err)
	defer cleanUpMetric(t, metric.ID)

	err = metricRepo.Exist(context.Background(), metric.CharacterID, metric.ID)
	assert.Nil(t, err)
}

func TestUnassignCategory(t *testing.T) {
	categoryID := mongodb.GenObjectID()
	characterIDs := []string{
		mongodb.GenObjectID(),
		mongodb.GenObjectID(),
	}
	for i := range 2 {
		metric := NewMetric()
		metric.CategoryID = lo.ToPtr(categoryID)
		metric.CharacterID = characterIDs[i]
		_, err := metricRepo.InsertOne(context.Background(), metric)
		assert.Nil(t, err)
	}

	count, err := metricRepo.CountByCategoryID(context.Background(), categoryID)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), count)

	err = metricRepo.UnassignCategory(context.Background(), categoryID)
	assert.Nil(t, err)

	count, err = metricRepo.CountByCategoryID(context.Background(), categoryID)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), count)

	err = metricRepo.DeleteByCharacterIDs(context.Background(), characterIDs)
	assert.Nil(t, err)

	count, err = metricRepo.Count(context.Background(), bson.M{})
	assert.Nil(t, err)
	assert.Equal(t, int64(0), count)
}

func cleanUpMetric(t *testing.T, id string) {
	err := metricRepo.DeleteByID(context.Background(), id)
	assert.Nil(t, err)
}
