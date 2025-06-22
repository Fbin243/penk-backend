package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"

	"github.com/jinzhu/copier"
)

func (b *MetricBusiness) Upsert(ctx context.Context, metricInput *entity.MetricInput) (*entity.Metric, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
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

	err = b.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, permEntities)
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
		return b.metricRepo.FindAndUpdateByID(ctx, *metricInput.ID, metric)
	}

	return b.metricRepo.InsertOne(ctx, metric)
}

func (b *MetricBusiness) Delete(ctx context.Context, metricID string) (*entity.Metric, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	err = b.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, []PermissionEntity{
		{
			ID:   metricID,
			Type: entity.EntityTypeMetric,
		},
	})
	if err != nil {
		return nil, err
	}

	return b.metricRepo.FindAndDeleteByID(ctx, metricID)
}
