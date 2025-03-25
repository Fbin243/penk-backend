package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"
)

type MetricBusiness struct {
	metricRepo    IMetricRepo
	characterRepo ICharacterRepo
	cateRepo      ICategoryRepo
}

func NewMetricBusiness(metricRepo IMetricRepo, characterRepo ICharacterRepo, cateRepo ICategoryRepo) *MetricBusiness {
	return &MetricBusiness{metricRepo, characterRepo, cateRepo}
}

func (b *MetricBusiness) GetMetrics(ctx context.Context, characterID string) ([]entity.Metric, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	err := b.characterRepo.Exist(ctx, authSession.ProfileID, characterID)
	if err != nil {
		return nil, err
	}

	return b.metricRepo.FindByCharacterID(ctx, characterID)
}

func (b *MetricBusiness) UpsertMetric(ctx context.Context, metricInput entity.MetricInput) (*entity.Metric, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	err := b.characterRepo.Exist(ctx, authSession.ProfileID, metricInput.CharacterID)
	if err != nil {
		return nil, err
	}

	if metricInput.CategoryID != nil {
		err = b.cateRepo.Exist(ctx, metricInput.CharacterID, *metricInput.CategoryID)
		if err != nil {
			return nil, err
		}
	}

	metric := &entity.Metric{
		BaseEntity: &base.BaseEntity{},
	}
	if metricInput.ID == nil {
		count, err := b.metricRepo.CountByCharacterID(ctx, metricInput.CharacterID)
		if err != nil {
			return nil, err
		}
		if count >= utils.LimitedMetricNumber {
			return nil, errors.ErrLimitMetric
		}
	} else {
		metric, err = b.metricRepo.FindByID(ctx, *metricInput.ID)
		if err != nil {
			return nil, err
		}
	}

	metric.Name = metricInput.Name
	metric.CharacterID = metricInput.CharacterID
	if metricInput.CategoryID != nil {
		metric.CategoryID = metricInput.CategoryID
	}
	metric.Value = metricInput.Value
	metric.Unit = metricInput.Unit

	if metricInput.ID != nil {
		return b.metricRepo.UpdateByID(ctx, *metricInput.ID, metric)
	}

	return b.metricRepo.InsertOne(ctx, metric)
}

func (b *MetricBusiness) DeleteMetric(ctx context.Context, metricID string) (*entity.Metric, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	metric, err := b.metricRepo.FindByID(ctx, metricID)
	if err != nil {
		return nil, err
	}

	err = b.characterRepo.Exist(ctx, authSession.ProfileID, metric.CharacterID)
	if err != nil {
		return nil, err
	}

	return b.metricRepo.DeleteByID(ctx, metricID)
}
