package business

import (
	"context"

	"tenkhours/pkg/auth"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"

	"github.com/samber/lo"
)

func (b *HabitBusiness) GetHabitLogs(ctx context.Context, filter *entity.HabitLogFilter, orderBy *entity.HabitLogOrderBy, limit, offset *int) ([]entity.HabitLog, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	var habitID *string
	if filter != nil {
		habitID = filter.HabitID
	}

	var habitLogs []entity.HabitLog
	if habitID == nil {
		// Get habit logs of the current character
		habits, err := b.habitRepo.FindByCharacterID(ctx, authSession.CurrentCharacterID)
		if err != nil {
			return nil, err
		}

		habitIDs := lo.Map(habits, func(habit entity.Habit, _ int) string {
			return habit.ID
		})

		filter.HabitIDs = habitIDs
		habitLogs, err = b.habitLogRepo.FindByPineline(ctx, entity.HabitLogPineline{
			Filter:  filter,
			OrderBy: orderBy,
			Limit:   limit,
			Offset:  offset,
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

		habitLogs, err = b.habitLogRepo.FindByPineline(ctx, entity.HabitLogPineline{
			Filter:  filter,
			OrderBy: orderBy,
			Limit:   limit,
			Offset:  offset,
		})
		if err != nil {
			return nil, err
		}
	}

	return habitLogs, nil
}
