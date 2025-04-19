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

	var habitID *string
	if filter != nil {
		habitID = filter.HabitID
	} else {
		filter = &entity.HabitLogFilter{}
	}

	var habitLogs []entity.HabitLog
	if habitID == nil {
		// Get habit logs of the current character
		habits, err := b.habitRepo.Find(ctx, entity.HabitPipeline{
			Filter: &entity.HabitFilter{
				CharacterID: &authSession.CurrentCharacterID,
			},
		})
		if err != nil {
			return nil, err
		}

		habitIDs := lo.Map(habits, func(habit entity.Habit, _ int) string {
			return habit.ID
		})

		filter.HabitIDs = habitIDs
		habitLogs, err = b.habitLogRepo.Find(ctx, entity.HabitLogPipeline{
			Filter:  filter,
			OrderBy: orderBy,
			Pagination: &types.Pagination{
				Limit:  limit,
				Offset: offset,
			},
		})
		if err != nil {
			return nil, err
		}
	} else {
		// Get habit logs of the habit
		err := b.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, []PermissionEntity{
			{
				ID:   *habitID,
				Type: entity.EntityTypeHabit,
			},
		})
		if err != nil {
			return nil, err
		}

		habitLogs, err = b.habitLogRepo.Find(ctx, entity.HabitLogPipeline{
			Filter:  filter,
			OrderBy: orderBy,
			Pagination: &types.Pagination{
				Limit:  limit,
				Offset: offset,
			},
		})
		if err != nil {
			return nil, err
		}
	}

	return habitLogs, nil
}
