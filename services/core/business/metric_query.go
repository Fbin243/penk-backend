package business

import (
	"context"

	"tenkhours/pkg/auth"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"
)

func (b *MetricBusiness) GetMetrics(ctx context.Context) ([]entity.Metric, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	return b.metricRepo.FindByCharacterID(ctx, authSession.CurrentCharacterID)
}
