package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/types"
	"tenkhours/services/core/entity"
)

func (b *MetricBusiness) Get(ctx context.Context, filter *entity.MetricFilter, orderBy *entity.MetricOrderBy, limit, offset *int) ([]entity.Metric, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
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
