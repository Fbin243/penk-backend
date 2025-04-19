package business

import (
	"context"

	"tenkhours/pkg/auth"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/types"
	"tenkhours/services/core/entity"
)

func (b *MetricBusiness) Get(ctx context.Context, filter *entity.MetricFilter, orderBy *entity.MetricOrderBy, limit, offset *int) ([]entity.Metric, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	filter.CharacterID = &authSession.CurrentCharacterID
	return b.metricRepo.Find(ctx, entity.MetricPipeline{
		Filter:  filter,
		OrderBy: orderBy,
		Pagination: &types.Pagination{
			Limit:  limit,
			Offset: offset,
		},
	})
}
