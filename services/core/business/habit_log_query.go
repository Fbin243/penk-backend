package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/types"
	"tenkhours/services/core/entity"

	"github.com/samber/lo"
)

func (b *HabitBusiness) GetHabitLogs(ctx context.Context, filter *entity.HabitLogFilter, orderBy *entity.HabitLogOrderBy, limit, offset *int) ([]entity.HabitLog, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	processedFilter, err := b.processFilter(ctx, authSession.CurrentCharacterID, filter)
	if err != nil {
		return nil, err
	}

	return b.habitLogRepo.Find(ctx, entity.HabitLogPipeline{
		Filter:  processedFilter,
		OrderBy: orderBy,
		Pagination: &types.Pagination{
			Limit:  limit,
			Offset: offset,
		},
	})
}

func (b *HabitBusiness) CountHabitLog(ctx context.Context, filter *entity.HabitLogFilter) (int, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return 0, err
	}

	processedFilter, err := b.processFilter(ctx, authSession.CurrentCharacterID, filter)
	if err != nil {
		return 0, err
	}

	return b.habitLogRepo.CountByFilter(ctx, processedFilter)
}

func (b *HabitBusiness) processFilter(
	ctx context.Context,
	characterID string,
	filter *entity.HabitLogFilter,
) (*entity.HabitLogFilter, error) {
	var habitID *string
	if filter != nil {
		habitID = filter.HabitID
	} else {
		filter = &entity.HabitLogFilter{}
	}

	if habitID == nil {
		// Get habit logs of the current character
		habits, err := b.habitRepo.Find(ctx, entity.HabitPipeline{
			Filter: &entity.HabitFilter{
				CharacterID: &characterID,
			},
		})
		if err != nil {
			return nil, err
		}
		habitIDs := lo.Map(habits, func(habit entity.Habit, _ int) string {
			return habit.ID
		})
		filter.HabitIDs = habitIDs
	} else {
		// Get habit logs of the habit
		err := b.permBiz.CheckOwnEntities(ctx, characterID, []PermissionEntity{
			{
				ID:   *habitID,
				Type: entity.EntityTypeHabit,
			},
		})
		if err != nil {
			return nil, err
		}
	}

	return filter, nil
}
