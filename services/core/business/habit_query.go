package business

import (
	"context"

	"tenkhours/pkg/auth"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/types"
	"tenkhours/services/core/entity"
)

func (b *HabitBusiness) Get(ctx context.Context, filter *entity.HabitFilter, orderBy *entity.HabitOrderBy, limit, offset *int) ([]entity.Habit, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	if filter == nil {
		filter = &entity.HabitFilter{}
	}
	filter.CharacterID = &authSession.CurrentCharacterID

	return b.habitRepo.Find(ctx, entity.HabitPipeline{
		Filter:  filter,
		OrderBy: orderBy,
		Pagination: &types.Pagination{
			Limit:  limit,
			Offset: offset,
		},
	})
}
