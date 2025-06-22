package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/types"
	"tenkhours/services/core/entity"
)

func (b *HabitBusiness) Get(ctx context.Context, filter *entity.HabitFilter, orderBy *entity.HabitOrderBy, limit, offset *int) ([]entity.Habit, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
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

func (b *HabitBusiness) Count(ctx context.Context, filter *entity.HabitFilter) (int, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return 0, err
	}

	if filter == nil {
		filter = &entity.HabitFilter{}
	}
	filter.CharacterID = &authSession.CurrentCharacterID

	return b.habitRepo.CountByFilter(ctx, filter)
}
