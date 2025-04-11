package business

import (
	"context"
	"time"

	"tenkhours/pkg/auth"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/types"
	"tenkhours/services/core/entity"

	"github.com/samber/lo"
)

func (b *HabitBusiness) GetHabitLogs(ctx context.Context, habitID *string, startTime, endTime time.Time) ([]entity.HabitLog, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
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

		habitLogs, err = b.habitLogRepo.FindByHabitIDs(ctx, habitIDs, &types.TimeFilter{
			StartTime: lo.ToPtr(startTime),
			EndTime:   lo.ToPtr(endTime),
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

		habitLogs, err = b.habitLogRepo.FindByHabitID(ctx, *habitID, &types.TimeFilter{
			StartTime: lo.ToPtr(startTime),
			EndTime:   lo.ToPtr(endTime),
		})
		if err != nil {
			return nil, err
		}
	}

	return habitLogs, nil
}
