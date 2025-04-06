package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"

	"github.com/jinzhu/copier"
)

func (b *MetricBusiness) UpsertMetric(ctx context.Context, metricInput *entity.MetricInput) (*entity.Metric, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	permEntities := []PermissionEntity{}
	if metricInput.ID != nil {
		permEntities = append(permEntities, PermissionEntity{
			ID:   *metricInput.ID,
			Type: entity.EntityTypeMetric,
		})
	}

	if metricInput.CategoryID != nil {
		permEntities = append(permEntities, PermissionEntity{
			ID:   *metricInput.CategoryID,
			Type: entity.EntityTypeCategory,
		})
	}

	err := b.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, permEntities)
	if err != nil {
		return nil, err
	}

	metric := &entity.Metric{
		BaseEntity:  &base.BaseEntity{},
		CharacterID: authSession.CurrentCharacterID,
	}
	if metricInput.ID == nil {
		count, err := b.metricRepo.CountByCharacterID(ctx, authSession.CurrentCharacterID)
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

	err = copier.Copy(metric, metricInput)
	if err != nil {
		return nil, err
	}

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

	err := b.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, []PermissionEntity{
		{
			ID:   metricID,
			Type: entity.EntityTypeMetric,
		},
	})
	if err != nil {
		return nil, err
	}

	return b.metricRepo.FindOneAndDeleteByID(ctx, metricID)
}
